use serde::{Deserialize, Serialize};
use crate::core::domain::cache_entry::CacheEntry;

#[derive(Serialize, Deserialize)]
pub struct CacheDump {
    pub store: std::collections::HashMap<String, CacheEntry>,
    pub lru_queue: Vec<String>,
}
