use std::collections::HashSet;
use crate::index::key::{KeyIndexStore, KeyRecord};

pub struct GcWorker;

impl GcWorker {
    pub fn mark(key_store: &KeyIndexStore) -> HashSet<[u8; 32]> {
        let mut live = HashSet::new();
        for r in key_store.mem.values() {
            if let KeyRecord::Put { hash, .. } = r {
                live.insert(*hash);
            }
        }
        live
    }
}