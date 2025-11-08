use std::time::{SystemTime, UNIX_EPOCH};
use serde::{Deserialize, Serialize};

#[derive(Clone, Serialize, Deserialize)]
pub struct CacheEntry {
    pub data: Vec<u8>,
    pub expires_at: Option<u64>,
    pub version: u64,
}

impl CacheEntry {
    pub fn new(data: Vec<u8>, ttl_secs: u64, version: u64) -> Self {
        let expires_at = Some(ttl_secs).map(|ttl| {
            let now = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs();
            now + ttl
        });
        Self {
            data,
            expires_at,
            version,
        }
    }

    pub fn is_expired(&self) -> bool {
        if let Some(ts) = self.expires_at {
            let now = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs();
            now >= ts
        } else {
            false
        }
    }

    pub fn ttl(&self) -> Option<u64> {
        self.expires_at.map(|exp| {
            let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs();
            exp.saturating_sub(now)
        })
    }
}