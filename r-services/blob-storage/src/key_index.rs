use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct KeyIndexEntry {
    pub key: String,
    pub version: u32,
    pub hash: [u8; 32],
}

pub struct KeyIndex {
    pub latest: HashMap<String, KeyIndexEntry>
}

impl KeyIndex {
    pub fn new() -> Self {
        Self {
            latest: HashMap::new()
        }
    }

    pub fn put(&mut self, key: String, hash: [u8; 32]) -> u32 {
        let version = self
            .latest
            .get(&key)
            .map(|e| e.version + 1)
            .unwrap_or(1);
        
        self.latest.insert(
            key.clone(),
            KeyIndexEntry {
                key,
                version,
                hash,
            },
        );
        version
    }
    
    pub fn get(&self, key: &str) -> Option<&KeyIndexEntry> {
        self.latest.get(key)
    }
}