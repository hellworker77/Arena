use crate::core::application::contracts::data::cache_stats::CacheStats;
use crate::core::application::contracts::data::config::Config;
use crate::core::application::contracts::data::eviction_policy::EvictionPolicy;
use crate::core::application::repository::cache_repository::CacheRepository;
use crate::core::domain::cache_entry::CacheEntry;
use anyhow::Result;
use async_trait::async_trait;
use std::collections::{HashMap, HashSet, VecDeque};
use tokio::fs;
use std::path::Path;
use std::sync::Arc;
use tokio::sync::mpsc;
use tokio::sync::mpsc::{Receiver, Sender};
use tokio::sync::Mutex;
use tokio::time::{Duration, Instant};
use tokio::{task, time};
use crate::core::application::contracts::data::cache_dump::CacheDump;

#[derive(Clone)]
pub struct LruCacheRepository {
    store: Arc<Mutex<HashMap<String, CacheEntry>>>,
    watched_keys: Arc<Mutex<HashSet<String>>>,
    lru_queue: Arc<Mutex<VecDeque<String>>>,
    pubsub_channels: Arc<Mutex<HashMap<String, Vec<Sender<Vec<u8>>>>>>,
    eviction_policy: Arc<Mutex<EvictionPolicy>>,
    config: Config,
}

impl LruCacheRepository {
    pub fn new(config: Config) -> Self {
        let repo = Self {
            store: Arc::new(Mutex::new(HashMap::new())),
            watched_keys: Arc::new(Mutex::new(HashSet::new())),
            lru_queue: Arc::new(Mutex::new(VecDeque::new())),
            pubsub_channels: Arc::new(Mutex::new(HashMap::new())),
            eviction_policy: Arc::new(Mutex::new(EvictionPolicy::NoEviction)),
            config,
        };

        let store_clone = repo.store.clone();

        let lru_clone = repo.lru_queue.clone();

        let config_clone = repo.config.clone();

        task::spawn(async move {
            let mut interval = time::interval(Duration::from_secs(
                config_clone.cache_config.cleanup_interval_secs.clone(),
            ));

            loop {
                interval.tick().await;
                let mut store = store_clone.lock().await;
                let mut lru_queue = lru_clone.lock().await;
                let now = Instant::now();

                let expired_keys: Vec<String> = store
                    .iter()
                    .filter(|(_, v)| v.is_expired())
                    .map(|(k, _)| k.clone())
                    .collect();

                for key in &expired_keys {
                    store.remove(key);
                    lru_queue.retain(|k| k != key);
                }
            }
        });

        repo
    }

    async fn enforce_lru(&self) {
        if *self.eviction_policy.lock().await != EvictionPolicy::LRU {
            return;
        }

        let mut store = self.store.lock().await;
        let mut lru_queue = self.lru_queue.lock().await;

        while  store.values().map(|e| e.data.len()).sum::<usize>()
            > self.config.cache_config.max_total_cache_size
        {
            if let Some(oldest) = lru_queue.pop_front() {
                store.remove(&oldest);
            } else {
                break;
            }
        }
    }

    async fn touch_key(&self, key: &str) {
        let mut lru_queue = self.lru_queue.lock().await;

        lru_queue.retain(|k| k != key);
        lru_queue.push_back(key.to_string());
    }

    async fn cleanup_expired_entries(&self) -> usize {
        let mut store = self.store.lock().await;
        let mut lru_queue = self.lru_queue.lock().await;

        let expired_keys: Vec<String> = store
            .iter()
            .filter(|(_, v)| v.is_expired())
            .map(|(k, _)| k.clone())
            .collect();

        for key in &expired_keys {
            store.remove(key);
            lru_queue.retain(|k| k != key);
        }

        expired_keys.len()
    }
}

#[async_trait]
impl CacheRepository for LruCacheRepository {
    async fn set(&self, key: &str, data: Vec<u8>) -> Result<()> {
        let mut store = self.store
            .lock()
            .await;


        let version = match store.get(key) {
            Some(entry) => entry.version + 1,
            None => 1,
        };

        store.insert(key.to_string(), CacheEntry::new(data, self.config.cache_config.cache_item_expiration_secs, version));

        self.touch_key(key).await;
        self.enforce_lru().await;
        Ok(())
    }

    async fn set_with_ttl(&self, key: &str, data: Vec<u8>, ttl_secs: u64) -> Result<()> {
        let mut store = self.store
            .lock()
            .await;

        let version = match store.get(key) {
            Some(entry) => entry.version + 1,
            None => 1,
        };

        store.insert(key.to_string(), CacheEntry::new(data, ttl_secs, version));

        self.touch_key(key).await;
        self.enforce_lru().await;
        Ok(())
    }

    async fn get(&self, key: &str) -> Result<Option<Vec<u8>>> {
        self.cleanup_expired_entries().await;
        if let Some(entry) = self.store.lock().await.get(key) {
            self.touch_key(key).await;
            return Ok(Some(entry.data.clone()));
        }
        Ok(None)
    }

    async fn delete(&self, key: &str) -> Result<()> {
        self.store.lock().await.remove(key);
        self.lru_queue.lock().await.retain(|k| k != key);
        Ok(())
    }

    async fn exists(&self, key: &str) -> Result<bool> {
        self.cleanup_expired_entries().await;
        Ok(self.store.lock().await.contains_key(key))
    }

    async fn clear(&self) -> Result<()> {
        self.store.lock().await.clear();
        self.lru_queue.lock().await.clear();
        Ok(())
    }

    async fn retain(&self, key: &str, duration_secs: u64) -> Result<()> {
        if let Some(entry) = self.store.lock().await.get_mut(key) {
            entry.expires_at = Some(duration_secs);
        }
        Ok(())
    }

    async fn ttl(&self, key: &str) -> Result<Option<u64>> {
        self.cleanup_expired_entries().await;
        let store = self.store.lock().await;
        if let Some(entry) = store.get(key) {
            return Ok(entry.ttl());
        }
        Ok(None)
    }

    async fn mget(&self, keys: &[String]) -> Result<Vec<Option<Vec<u8>>>> {
        self.cleanup_expired_entries().await;
        let store = self.store.lock().await;
        let result = keys
            .iter()
            .map(|k| store.get(k).map(|e| e.data.clone()))
            .collect();
        Ok(result)
    }

    async fn mset(&self, items: Vec<(String, Vec<u8>)>) -> Result<()> {
        let mut store = self.store.lock().await;
        for (k, v) in items {
            let version = match store.get(&k) {
                Some(entry) => entry.version + 1,
                None => 1,
            };
            store.insert(
                k.clone(),
                CacheEntry::new(v, self.config.cache_config.cache_item_expiration_secs, version),
            );
            self.touch_key(&k).await;
        }
        self.enforce_lru().await;
        Ok(())
    }

    async fn publish(&self, channel: &str, payload: Vec<u8>) -> Result<()> {
        let channels = self.pubsub_channels.lock().await;
        if let Some(subs) = channels.get(channel) {
            for tx in subs {
                let _ = tx.send(payload.clone()).await;
            }
        }
        Ok(())
    }

    async fn subscribe(&self, channel: &str) -> Result<Receiver<Vec<u8>>> {
        let (tx, rx) = mpsc::channel(100);
        self.pubsub_channels.lock().await.entry(channel.to_string()).or_default().push(tx);
        Ok(rx)
    }

    async fn flushdb(&self) -> Result<()> {
        self.clear().await
    }

    async fn flushall(&self) -> Result<()> {
        self.clear().await
    }

    async fn len(&self) -> Result<usize> {
        self.cleanup_expired_entries().await;
        Ok(self.store.lock().await.len())
    }

    async fn keys(&self) -> Result<Vec<String>> {
        self.cleanup_expired_entries().await;
        Ok(self.store.lock().await.keys().cloned().collect())
    }

    async fn cleanup_expired(&self) -> Result<usize> {
        Ok(self.cleanup_expired_entries().await)
    }

    async fn stats(&self) -> Result<CacheStats> {
        Ok(CacheStats {
            items: self.store.lock().await.len(),
            total_size_bytes: self.store.lock().await.values().map(|e| e.data.len()).sum(),
            hits: 0,
            misses: 0,
        })
    }

    async fn set_eviction_policy(&self, policy: EvictionPolicy) -> Result<()> {
        *self.eviction_policy.lock().await = policy;
        Ok(())
    }

    async fn get_eviction_policy(&self) -> Result<EvictionPolicy> {
        Ok(*self.eviction_policy.lock().await)
    }

    async fn save_to_disk(&self) -> Result<()> {
        let store = self.store.lock().await;
        let lru_queue = self.lru_queue.lock().await;

        let dump = CacheDump {
            store: store.clone(),
            lru_queue: lru_queue.iter().cloned().collect(),
        };

        let bytes = serde_json::to_vec(&dump)?;
        fs::write(&self.config.cache_config.dump_file_path, bytes).await?;
        Ok(())
    }

    async fn load_from_disk(&self) -> Result<()> {
        let path = Path::new(&self.config.cache_config.dump_file_path);

        if !path.exists() {
            return Ok(());
        }

        let bytes = fs::read(path).await?;
        let dump: CacheDump = serde_json::from_slice(&bytes)?;

        let mut store = self.store.lock().await;
        let mut lru_queue = self.lru_queue.lock().await;

        *store = dump.store;
        *lru_queue = std::collections::VecDeque::from(dump.lru_queue);
        Ok(())
    }

    async fn backup(&self) -> Result<()> {
        self.save_to_disk().await
    }

    async fn restore(&self) -> Result<()> {
        self.load_from_disk().await
    }

    async fn snapshot(&self) -> Result<()> {
        self.save_to_disk().await
    }

    async fn restore_snapshot(&self) -> Result<()> {
        self.load_from_disk().await
    }

    async fn rename(&self, old_key: &str, new_key: &str) -> Result<()> {
        let mut store = self.store.lock().await;
        if let Some(entry) = store.remove(old_key) {
            store.insert(new_key.to_string(), entry);
            self.touch_key(new_key).await;
        }
        Ok(())
    }

    async fn expire_at(&self, key: &str, timestamp: u64) -> Result<()> {
        let mut store = self.store.lock().await;
        if let Some(entry) = store.get_mut(key) {
            entry.expires_at = Some(timestamp);
        }
        Ok(())
    }

    async fn dump(&self, key: &str) -> Result<Option<Vec<u8>>> {
        Ok(self.get(key).await?)
    }

    async fn restore_dump(&self, key: &str, data: Vec<u8>) -> Result<()> {
        self.set(key, data).await
    }

    async fn scan(&self, cursor: u64, count: usize) -> Result<(u64, Vec<String>)> {
        let keys: Vec<String> = self.keys().await?;
        let start = cursor as usize;
        let end = (start + count).min(keys.len());
        let next_cursor = if end >= keys.len() { 0 } else { end as u64 };
        Ok((next_cursor, keys[start..end].to_vec()))
    }

    async fn watch(&self, key: &str) -> Result<()> {
        let mut watched = self.watched_keys.lock().await;
        watched.insert(key.to_string());
        Ok(())
    }

    async fn unwatch(&self) -> Result<()> {
        let mut watched = self.watched_keys.lock().await;
        watched.clear();
        Ok(())
    }

    async fn check_watched(&self) -> Result<bool> {
        let store = self.store.lock().await;
        let watched = self.watched_keys.lock().await;
        
        for key in watched.iter() {
            if !store.contains_key(key) {
                return Ok(true);
            }
        }
        
        Ok(false)
    }

    async fn eval(
        &self,
        _script: &str,
        _keys: Vec<String>,
        _args: Vec<Vec<u8>>,
    ) -> Result<Vec<u8>> {
        Ok(Vec::new())
    }
}