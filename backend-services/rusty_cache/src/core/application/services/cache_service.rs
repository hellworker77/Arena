use std::sync::mpsc::Receiver;
use anyhow::Result;
use async_trait::async_trait;
use serde::de::DeserializeOwned;
use serde::Serialize;
use crate::core::application::contracts::data::cache_stats::CacheStats;
use crate::core::application::contracts::data::eviction_policy::EvictionPolicy;

#[async_trait]
pub trait CacheService: Send + Sync {
    async fn set<T: Serialize + Send + Sync>(&self, key: &str, value: &T, ttl_secs: Option<u64>) -> Result<()>;
    
    async fn get<T: DeserializeOwned + Send + Sync>(&self, key: &str) -> Result<Option<T>>;
    
    async fn delete(&self, key: &str) -> Result<()>;
    
    async fn exists(&self, key: &str) -> Result<bool>;
    
    async fn clear(&self) -> Result<()>;
    
    async fn retain(&self, key: &str, duration_secs: u64) -> Result<()>;
    
    async fn ttl(&self, key: &str) -> Result<Option<u64>>;
    
    async fn mget<T: DeserializeOwned + Send + Sync>(&self, keys: &[String]) -> Result<Vec<Option<T>>>;
    
    async fn mset<T: Serialize + Send + Sync>(&self, items: Vec<(String, T)>, ttl_secs: Option<u64>) -> Result<()>;
    
    async fn publish<T: Serialize + Send + Sync>(&self, channel: &str, payload: T) -> Result<()>;
    
    async fn subscribe<T: DeserializeOwned + Send + Sync>(&self, channel: &str) -> Result<Receiver<T>>;
    
    async fn flushdb(&self) -> Result<()>;
    
    async fn flushall(&self) -> Result<()>;
    
    async fn len(&self) -> Result<usize>;
    
    async fn keys(&self) -> Result<Vec<String>>;
    
    async fn cleanup_expired(&self) -> Result<usize>;
    
    async fn stats(&self) -> Result<CacheStats>;
    
    async fn set_eviction_policy(&self, policy: EvictionPolicy) -> Result<()>;
    
    async fn get_eviction_policy(&self) -> Result<EvictionPolicy>;
    
    async fn save_to_disk(&self) -> Result<()>;
    
    async fn load_from_disk(&self) -> Result<()>;
    
    async fn backup(&self) -> Result<()>;
    
    async fn restore(&self) -> Result<()>;
    
    async fn snapshot(&self) -> Result<()>;
    
    async fn restore_snapshot(&self) -> Result<()>;
    
    async fn rename(&self, old_key: &str, new_key: &str) -> Result<()>;
    
    async fn expire_at(&self, key: &str, timestamp: u64) -> Result<()>;
    
    async fn dump(&self, key: &str) -> Result<Option<Vec<u8>>>;
    
    async fn restore_dump(&self, key: &str, data: Vec<u8>) -> Result<()>;
    
    async fn scan(&self, cursor: u64, count: usize) -> Result<(u64, Vec<String>)>;
    
    async fn watch(&self, key: &str) -> Result<()>;
    
    async fn unwatch(&self) -> Result<()>;

    async fn check_watched(&self) -> Result<bool>;
    
    async fn eval(&self, script: &str, keys: Vec<String>, args: Vec<Vec<u8>>) -> Result<Vec<u8>>;
}