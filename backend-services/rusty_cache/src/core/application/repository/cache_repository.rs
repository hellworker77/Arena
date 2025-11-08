use crate::core::application::contracts::data::cache_stats::CacheStats;
use crate::core::application::contracts::data::eviction_policy::EvictionPolicy;
use anyhow::Result;
use async_trait::async_trait;
use tokio::sync::mpsc::Receiver;

/// Trait defining the contract for cache repository.
///
/// Provides asynchronous methods for cache operations including:
/// - Basic CRUD & TTL management
/// - Batch operations (MGET / MSET)
/// - Pub/Sub messaging
/// - Eviction policies
/// - Persistence (save/load/backup/snapshot)
/// - Advanced Redis-like features (rename, expire_at, dump/restore, scan, watch/unwatch, eval)
#[async_trait]
pub trait CacheRepository: Send + Sync {
    // --------------- Basic Cache Operations ---------------
    /// Sets a value in the cache for the given key.
    ///
    /// # Arguments
    /// * `key` - The key under which the data will be stored.
    /// * `data` - The data to be stored in the cache.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn set(&self, key: &str, data: Vec<u8>) -> Result<()>;

    /// Sets a value in the cache for the given key with a specified TTL (time-to-live).
    ///
    /// # Arguments
    /// * `key` - The key under which the data will be stored.
    /// * `data` - The data to be stored in the cache.
    /// * `ttl_secs` - The time-to-live in seconds for the cache entry.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn set_with_ttl(&self, key: &str, data: Vec<u8>, ttl_secs: u64) -> Result<()>;

    /// Retrieves a value from the cache for the given key.
    ///
    /// # Arguments
    /// * `key` - The key whose associated value is to be returned.
    ///
    /// # Returns
    /// * `Result<Option<Vec<u8>>>` - The value associated with the key, or None if the key does not exist.
    async fn get(&self, key: &str) -> Result<Option<Vec<u8>>>;

    /// Deletes a value from the cache for the given key.
    ///
    /// # Arguments
    /// * `key` - The key whose associated value is to be deleted.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn delete(&self, key: &str) -> Result<()>;

    /// Checks if a key exists in the cache.
    ///
    /// # Arguments
    /// * `key` - The key to check for existence.
    ///
    /// # Returns
    /// * `Result<bool>` - True if the key exists, false otherwise.
    async fn exists(&self, key: &str) -> Result<bool>;

    /// Clears all entries from the cache.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn clear(&self) -> Result<()>;

    /// Retains a key in the cache for a specified duration.
    ///
    /// # Arguments
    /// * `key` - The key to retain.
    /// * `duration_secs` - The duration in seconds to retain the key.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn retain(&self, key: &str, duration_secs: u64) -> Result<()>;

    /// Retrieves the TTL (time-to-live) for a given key in the cache.
    ///
    /// # Arguments
    /// * `key` - The key whose TTL is to be retrieved.
    ///
    /// # Returns
    /// * `Result<Option<u64>>` - The TTL in seconds, or None if the key does not exist or has no TTL.
    async fn ttl(&self, key: &str) -> Result<Option<u64>>;

    // --------------- Batch Operations ---------------
    /// Retrieves multiple values from the cache for the given keys.
    ///
    /// # Arguments
    /// * `keys` - A slice of keys whose associated values are to be returned.
    ///
    /// # Returns
    /// * `Result<Vec<Option<Vec<u8>>>>` - A vector of values associated with the keys, where each value is an Option.
    async fn mget(&self, keys: &[String]) -> Result<Vec<Option<Vec<u8>>>>;

    /// Sets multiple key-value pairs in the cache.
    ///
    /// # Arguments
    /// * `items` - A vector of tuples containing key-value pairs to be stored in the cache.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn mset(&self, items: Vec<(String, Vec<u8>)>) -> Result<()>;

    // --------------- Pub/Sub ---------------
    /// Publishes a message to a specified channel.
    ///
    /// # Arguments
    /// * `channel` - The channel to which the message will be published.
    /// * `payload` - The message payload to be published.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn publish(&self, channel: &str, payload: Vec<u8>) -> Result<()>;

    /// Subscribes to a specified channel to receive messages.
    ///
    /// # Arguments
    /// * `channel` - The channel to subscribe to.
    ///
    /// # Returns
    /// * `Result<tokio::sync::mpsc::Receiver<Vec<u8>>>` - A receiver to receive messages from the subscribed channel.
    async fn subscribe(&self, channel: &str) -> Result<Receiver<Vec<u8>>>;

    // --------------- Cleaning ---------------
    /// Flushes the current database in the cache.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn flushdb(&self) -> Result<()>;

    /// Flushes all databases in the cache.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn flushall(&self) -> Result<()>;

    // --------------- Metadata and management ---------------
    /// Retrieves the number of entries in the cache.
    ///
    /// # Returns
    /// * `Result<usize>` - The number of entries in the cache.
    async fn len(&self) -> Result<usize>;

    /// Retrieves all keys currently stored in the cache.
    ///
    /// # Returns
    /// * `Result<Vec<String>>` - A vector of all keys in the cache.
    async fn keys(&self) -> Result<Vec<String>>;

    /// Cleans up expired entries in the cache.
    ///
    /// # Returns
    /// * `Result<usize>` - The number of expired entries that were removed.
    async fn cleanup_expired(&self) -> Result<usize>;

    /// Retrieves statistics about the cache.
    ///
    /// # Returns
    /// * `Result<CacheStats>` - Statistics about the cache usage and performance.
    async fn stats(&self) -> Result<CacheStats>;

    // -------------- Eviction policy ---------------
    /// Sets the eviction policy for the cache.
    ///
    /// # Arguments
    /// * `policy` - The eviction policy to be set.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn set_eviction_policy(&self, policy: EvictionPolicy) -> Result<()>;

    /// Gets the current eviction policy of the cache.
    ///
    /// # Returns
    /// * `Result<EvictionPolicy>` - The current eviction policy of the cache.
    async fn get_eviction_policy(&self) -> Result<EvictionPolicy>;

    //  ---------- Persistence ----------
    /// Saves the current state of the cache to disk.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn save_to_disk(&self) -> Result<()>;

    /// Loads the cache state from disk.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn load_from_disk(&self) -> Result<()>;

    /// Creates a backup of the current cache state.
    /// 
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn backup(&self) -> Result<()>;

    /// Restores the cache state from a backup.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn restore(&self) -> Result<()>;

    /// Creates a snapshot of the current cache state.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn snapshot(&self) -> Result<()>;

    /// Restores the cache state from a snapshot.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn restore_snapshot(&self) -> Result<()>;

    //  ---------- Advanced Redis-like Features ----------
    /// Renames an existing key to a new key.
    ///
    /// # Arguments
    /// * `old_key` - The current key name.
    /// * `new_key` - The new key name.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn rename(&self, old_key: &str, new_key: &str) -> Result<()>;

    /// Sets the expiration time for a key as an absolute UNIX timestamp.
    ///
    /// # Arguments
    /// * `key` - The key to expire.
    /// * `timestamp` - UNIX timestamp at which the key will expire.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn expire_at(&self, key: &str, timestamp: u64) -> Result<()>;

    /// Dumps the serialized value of a key for backup or transfer.
    ///
    /// # Arguments
    /// * `key` - The key whose value is to be dumped.
    ///
    /// # Returns
    /// * `Result<Option<Vec<u8>>>` - The serialized data of the key, or None if not found.
    async fn dump(&self, key: &str) -> Result<Option<Vec<u8>>>;

    /// Restores a dumped key into the cache.
    ///
    /// # Arguments
    /// * `key` - The key to restore.
    /// * `data` - The serialized data to restore.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn restore_dump(&self, key: &str, data: Vec<u8>) -> Result<()>;

    /// Iterates over keys in the cache, supporting cursor-based scanning.
    ///
    /// # Arguments
    /// * `cursor` - The current cursor position (0 for new scan).
    /// * `count` - The number of items to return per iteration.
    ///
    /// # Returns
    /// * `Result<(u64, Vec<String>)>` - The next cursor and a list of keys.
    async fn scan(&self, cursor: u64, count: usize) -> Result<(u64, Vec<String>)>;

    /// Watches a key for changes (used for optimistic transactions).
    ///
    /// # Arguments
    /// * `key` - The key to watch.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn watch(&self, key: &str) -> Result<()>;

    /// Unwatches all keys currently being watched.
    ///
    /// # Returns
    /// * `Result<()>` - Indicates success or failure of the operation.
    async fn unwatch(&self) -> Result<()>;
    
    /// Checks if any watched keys have been modified.
    /// 
    /// # Returns
    /// * `Result<bool>` - True if any watched keys were modified, false otherwise.
    async fn check_watched(&self) -> Result<bool>;

    /// Evaluates a script (e.g., Lua) in the cache context.
    ///
    /// # Arguments
    /// * `script` - The script code to execute.
    /// * `keys` - The list of keys referenced by the script.
    /// * `args` - The arguments to pass to the script.
    ///
    /// # Returns
    /// * `Result<Vec<u8>>` - The result of script execution.
    async fn eval(&self, script: &str, keys: Vec<String>, args: Vec<Vec<u8>>) -> Result<Vec<u8>>;
}