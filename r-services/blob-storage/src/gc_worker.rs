use std::collections::HashSet;
use crate::cas_index::CasIndex;
use crate::key_index::KeyIndex;

pub struct GcWorker;

impl GcWorker {
    pub fn run(index: &KeyIndex, cas: &mut CasIndex) {
        let mut live = HashSet::new();

        for e in index.latest.values() {
            live.insert(e.hash);
        }

        cas.map.retain(|h, _| live.contains(h));
    }

    ///Phase 1: MARK - compute live hashes from KeyIndex snapshot
    pub fn mark(index: &KeyIndex)-> HashSet<[u8; 32]> {
        index.latest.values().map(|e| e.hash).collect()
    }

    ///Phase 2: SWEEP - remove uncreachable CAS entries
    pub fn sweep(cas: &mut CasIndex, live: &HashSet<[u8; 32]>) {
        cas.map.retain(|h, _| live.contains(h));
    }
}