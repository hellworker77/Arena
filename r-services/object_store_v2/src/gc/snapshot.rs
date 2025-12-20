use std::collections::{HashMap, HashSet};
use crate::error::error::StoreResult;
use crate::gc::types::{GcSnapshot, Hash32};
use crate::segment::{scan_all_segments, SegmentObjectInfo};

pub fn build_snapshot(
    // These are your current in-memory views.
    cas_entries: HashMap<Hash32, crate::store::CasEntry>,
    key_store: &crate::index::key::KeyIndexStore,
    segments: HashMap<u64, std::path::PathBuf>,
    active_id: u64,
    sealed_ids: HashSet<u64>,
) -> StoreResult<GcSnapshot> {
    // Mark: compute live hashes from latest KeyIndex view.
    let live_hashes = mark_live_hashes(key_store)?;

    // Scan sealed segments only (validation source of truth).
    let sealed_segments: HashMap<u64, std::path::PathBuf> = sealed_ids
        .iter()
        .filter_map(|id| segments.get(id).map(|p| (*id, p.clone())))
        .collect();

    let sealed_scan: HashMap<u64, HashMap<Hash32, SegmentObjectInfo>> =
    scan_all_segments(&sealed_segments)?;

    Ok(GcSnapshot{
        live_hashes,
        cas_entries,
        sealed_scan,
        segments,
        active_id,
        sealed_ids,
    })
}

fn mark_live_hashes(
    key_store: &crate::index::key::KeyIndexStore,
) -> StoreResult<HashSet<Hash32>> {
    let mut live = HashSet::new();

    // The KeyIndexStore API may differ in your codebase.
    // This assumes you can iterate a latest-view (mem + sst merged).
    // If you only have mem right now, this is still correct for correctness in dev mode.
    for rec in key_store.iter_latest()? {
        if let crate::index::key::KeyRecord::Put { hash, .. } = rec {
            live.insert(hash);
        }
    }

    Ok(live)
}