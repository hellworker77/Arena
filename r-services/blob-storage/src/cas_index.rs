use std::collections::{HashMap, HashSet};

#[derive(Debug, Clone)]
pub struct CasEntry {
    pub segment_id: u64,
    pub offset: u64,
    pub size: u64,
    pub refcount: u64,
}

pub struct CasIndex {
    pub map: HashMap<[u8; 32], CasEntry>,
}

impl CasIndex {
    pub fn new() -> Self {
        CasIndex {
            map: HashMap::new(),
        }
    }
    
    pub fn insert(&mut self, hash: [u8; 32], entry: CasEntry) {
        self.map.insert(hash, entry);
    }

    pub fn inc_ref(&mut self, hash: &[u8; 32]) {
        if let Some(e) = self.map.get_mut(hash) {
            e.refcount += 1;
        }
    }


    pub fn dec_ref(&mut self, hash: &[u8; 32]) {
        if let Some(e) = self.map.get_mut(hash) {
            e.refcount -= 1;
        }
    }
    
    pub fn live_hashes(&self) -> HashSet<[u8; 32]> {
       self.map
           .iter()
           .filter(|(_, e)| e.refcount > 0)
           .map(|(h, _)| *h)
           .collect()
    }
}