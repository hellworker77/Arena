use crate::core::application::contracts::data::cache_stats::CacheStats;
use crate::core::application::contracts::data::eviction_policy::EvictionPolicy;
use crate::core::application::repository::cache_repository::CacheRepository;
use crate::core::application::services::cache_service::CacheService;
use anyhow::Result;
use async_trait::async_trait;
use serde::Serialize;
use serde::de::DeserializeOwned;
use std::sync::Arc;
use std::sync::mpsc::Receiver;

pub struct CacheServiceImpl<R: CacheRepository> {
    repository: Arc<R>,
}

impl<R: CacheRepository> CacheServiceImpl<R> {
    pub fn new(repository: Arc<R>) -> Self { Self { repository, } }
}

#[async_trait]
impl<R: CacheRepository + Send + Sync> CacheService for CacheServiceImpl<R> {
    async fn set<T: Serialize + Send + Sync>(
        &self,
        key: &str,
        value: &T,
        ttl_secs: Option<u64>,
    ) -> Result<()> {
        let data = serde_json::to_vec(value)?;
        if let Some(ttl_secs) = ttl_secs {
            self.repository.set_with_ttl(key, data, ttl_secs).await
        } else {
            self.repository.set(key, data).await
        }
    }

    async fn get<T: DeserializeOwned + Send + Sync>(&self, key: &str) -> Result<Option<T>> {
        if let Some(bytes) = self.repository.get(key).await? {
            let value = serde_json::from_slice::<T>(&bytes)?;
            Ok(Some(value))
        } else {
            Ok(None)
        }
    }

    async fn delete(&self, key: &str) -> Result<()> {
        self.repository.delete(key).await
    }

    async fn exists(&self, key: &str) -> Result<bool> {
        self.repository.exists(key).await
    }

    async fn clear(&self) -> Result<()> {
        self.repository.clear().await
    }

    async fn retain(&self, key: &str, duration_secs: u64) -> Result<()> {
        self.repository.retain(key, duration_secs).await
    }

    async fn ttl(&self, key: &str) -> Result<Option<u64>> {
        self.repository.ttl(key).await
    }

    async fn mget<T: DeserializeOwned + Send + Sync>(
        &self,
        keys: &[String],
    ) -> Result<Vec<Option<T>>> {
        let bytes_list = self.repository.mget(keys).await?;
        let mut result = Vec::with_capacity(bytes_list.len());
        for bytes in bytes_list {
            if let Some(data) = bytes {
                result.push(Some(serde_json::from_slice::<T>(&data)?));
            } else {
                result.push(None);
            }
        }
        Ok(result)
    }

    async fn mset<T: Serialize + Send + Sync>(
        &self,
        items: Vec<(String, T)>,
        ttl_secs: Option<u64>,
    ) -> Result<()> {
        let data: Vec<(String, Vec<u8>)> = items
            .into_iter()
            .map(|(k, v)| (k, serde_json::to_vec(&v).unwrap()))
            .collect();

        if let Some(ttl_secs) = ttl_secs {
            for (k, v) in data {
                self.repository.set_with_ttl(&k, v, ttl_secs).await?;
            }
            Ok(())
        } else {
            self.repository.mset(data).await
        }
    }

    async fn publish<T: Serialize + Send + Sync>(&self, channel: &str, payload: T) -> Result<()> {
        let payload_bytes = serde_json::to_vec(&payload)?;

        self.repository.publish(channel, payload_bytes).await
    }

    async fn subscribe<T: DeserializeOwned + Send + Sync>(
        &self,
        channel: &str,
    ) -> Result<Receiver<T>> {
        unimplemented!("Subscribe method is not implemented yet")
    }

    async fn flushdb(&self) -> Result<()> {
        self.repository.flushdb().await
    }

    async fn flushall(&self) -> Result<()> {
        self.repository.flushall().await
    }

    async fn len(&self) -> Result<usize> {
        self.repository.len().await
    }

    async fn keys(&self) -> Result<Vec<String>> {
        self.repository.keys().await
    }

    async fn cleanup_expired(&self) -> Result<usize> {
        self.repository.cleanup_expired().await
    }

    async fn stats(&self) -> Result<CacheStats> {
        self.repository.stats().await
    }

    async fn set_eviction_policy(&self, policy: EvictionPolicy) -> Result<()> {
        self.repository.set_eviction_policy(policy).await
    }

    async fn get_eviction_policy(&self) -> Result<EvictionPolicy> {
        self.repository.get_eviction_policy().await
    }

    async fn save_to_disk(&self) -> Result<()> {
        self.repository.save_to_disk().await
    }

    async fn load_from_disk(&self) -> Result<()> {
        self.repository.load_from_disk().await
    }

    async fn backup(&self) -> Result<()> {
        self.repository.backup().await
    }

    async fn restore(&self) -> Result<()> {
        self.repository.restore().await
    }

    async fn snapshot(&self) -> Result<()> {
        self.repository.snapshot().await
    }

    async fn restore_snapshot(&self) -> Result<()> {
        self.repository.restore_snapshot().await
    }

    async fn rename(&self, old_key: &str, new_key: &str) -> Result<()> {
        self.repository.rename(old_key, new_key).await
    }

    async fn expire_at(&self, key: &str, timestamp: u64) -> Result<()> {
        self.repository.expire_at(key, timestamp).await
    }

    async fn dump(&self, key: &str) -> Result<Option<Vec<u8>>> {
        self.repository.dump(key).await
    }

    async fn restore_dump(&self, key: &str, data: Vec<u8>) -> Result<()> {
        self.repository.restore_dump(key, data).await
    }

    async fn scan(&self, cursor: u64, count: usize) -> Result<(u64, Vec<String>)> {
        self.repository.scan(cursor, count).await
    }

    async fn watch(&self, key: &str) -> Result<()> {
        self.repository.watch(key).await
    }

    async fn unwatch(&self) -> Result<()> {
        self.repository.unwatch().await
    }

    async fn check_watched(&self) -> Result<bool> {
        self.repository.check_watched().await
    }

    async fn eval(&self, script: &str, keys: Vec<String>, args: Vec<Vec<u8>>) -> Result<Vec<u8>> {
        self.repository.eval(script, keys, args).await
    }
}
