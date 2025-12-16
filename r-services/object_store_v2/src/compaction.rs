use crate::segment::{scan_segment, SegmentWriter};
use crate::store::{CasEntry, CasIndex};
use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::{Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};
use crate::error::error::{StoreError, StoreResult};

/// Rewrite new segment with only live objects.
/// Production notes:
/// - We validate offsets using a scanned segment index, not blind CAS offsets.
/// - We preserve stored bytes (nonce+cipher) and do not re-encrypt here.
pub struct Compactor;

impl Compactor {
    pub fn compact(
        live: &HashSet<[u8; 32]>,
        cas: &mut CasIndex,
        segments: &HashMap<u64, PathBuf>,
        out_path: &Path,
        out_segment_id: u64,
    ) -> StoreResult<PathBuf> {
        // Pre-scan segments to validate what actually exists
        let mut seg_scan: HashMap<u64, HashMap<[u8; 32], crate::segment::SegmentObjectInfo>> = HashMap::new();
        for (id, p) in segments {
            seg_scan.insert(*id, scan_segment(p)?);
        }

        let mut writer = SegmentWriter::create(out_path, out_segment_id)?;
        let mut new_map: HashMap<[u8; 32], CasEntry> = HashMap::new();

        for (hash, entry) in cas.map.iter() {
            if !live.contains(hash) || entry.refcount <= 0 {
                continue;
            }
            let seg_path = segments.get(&entry.segment_id).ok_or(StoreError::SegmentMissing)?;
            let scan = seg_scan.get(&entry.segment_id).ok_or(StoreError::SegmentMissing)?;
            let info = scan.get(hash).ok_or(StoreError::CasDanglingObject)?;

            // Read exactly the object at scanned offset (trusted more than CAS offset)
            let mut f = File::open(seg_path)?;
            f.seek(SeekFrom::Start(info.offset))?;

            // read & copy: hash + nonce + size_plain + size_cipher + cipher
            let mut stored_hash = [0u8; 32];
            f.read_exact(&mut stored_hash)?;
            if stored_hash != *hash {
                return Err(StoreError::HashMismatch);
            }

            let mut nonce = [0u8; 12];
            f.read_exact(&mut nonce)?;

            let mut sp = [0u8; 8];
            f.read_exact(&mut sp)?;
            let size_plain = u64::from_le_bytes(sp);

            let mut sc = [0u8; 8];
            f.read_exact(&mut sc)?;
            let size_cipher = u64::from_le_bytes(sc);

            let mut cipher = vec![0u8; size_cipher as usize];
            f.read_exact(&mut cipher)?;

            let new_off = writer.write_object(*hash, nonce, &cipher, size_plain)?;

            new_map.insert(
                *hash,
                CasEntry {
                    segment_id: out_segment_id,
                    offset: new_off,
                    size: size_cipher,
                    refcount: entry.refcount,
                },
            );
        }

        let new_path = writer.seal()?;
        cas.map = new_map;
        Ok(new_path)
    }
}