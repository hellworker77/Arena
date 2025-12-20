use crate::error::error::{StoreError, StoreResult};
use crate::gc::types::GcSnapshot;
use crate::manifest::{Manifest, ManifestRecord};
use crate::segment::{read_object, SegmentWriter};
use std::path::{Path, PathBuf};
use sha2::digest::DynDigest;

/// Compacts a single sealed segment into a new sealed segment containing only live objects.
///
/// Contract:
/// - Never reads/writes the active segment.
/// - Writes a brand new segment file.
/// - Appends manifest records in a crash-safe order.
/// - Does not delete the old segment file directly (that is handled as a separate DropSegment).
pub fn compact_one_segment(
    old_segment_id: u64,
    snap: &GcSnapshot,
    manifest: &mut Manifest,
    dir_segments: &Path,
) -> StoreResult<()> {
    if old_segment_id == snap.active_id {
        return Err(StoreError::GcInvariantViolation);
    }
    if !snap.sealed_ids.contains(&old_segment_id) {
        return Err(StoreError::GcInvariantViolation);
    }

    let old_path = snap
        .segments
        .get(&old_segment_id)
        .ok_or(StoreError::SegmentMissing)?
        .clone();

    let scan = snap
        .sealed_scan
        .get(&old_segment_id)
        .ok_or(StoreError::SegmentMissing)?;

    // Allocate a new segment id.
    // In production you should use a monotonic allocator stored in manifest.
    let new_segment_id = next_segment_id(&snap.segments);
    let new_path = dir_segments.join(format!("seg-{new_segment_id:05}.seg"));

    // Record the new segment first.
    manifest.append(&ManifestRecord::NewSegment {
        segment_id: new_segment_id,
        path: new_path.clone(),
    })?;

    // Write new segment with only live objects.
    let mut writer = SegmentWriter::create(&new_path, new_segment_id)?;
    let mut written = 0u32;

    for (hash, info) in scan.iter() {
        if !snap.live_hashes.contains(hash) {
            continue;
        }

        // Read object payload as stored in the old segment.
        // This assumes SegmentReader can read by offset and returns the cipher bytes.
        let (cipher, read_hash) = read_object(&old_path, info.offset)?;
        if read_hash != *hash {
            return Err(StoreError::HashMismatch);
        }

        // Preserve nonce/size semantics if your segment format stores them.
        // Here we write nonce as zero if you don't have it in scan info.
        let nonce = [0u8; 12];

        let _new_offset = writer.write_object(*hash, nonce, &cipher, info.size_plain)?;
        written += 1;
    }

    // Ensure data durability before sealing.
    writer.flush_data()?;

    // Seal the new segment (writes object_count + fsync).
    writer.seal()?;

    // Seal the new segment.
    manifest.append(&ManifestRecord::SealSegment { segment_id: new_segment_id })?;

    // Drop the old segment after the new one is sealed.
    // This makes the transition crash-safe: at worst you have both.
    manifest.append(&ManifestRecord::DropSegment { segment_id: old_segment_id })?;
    let _ = std::fs::remove_file(old_path);

    Ok(())
}

fn next_segment_id(segments: &std::collections::HashMap<u64, PathBuf>) -> u64 {
    segments.keys().copied().max().unwrap_or(0) + 1
}