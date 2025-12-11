use crate::models::eviction_policy::EvictionPolicy;
use anyhow::Result;
use async_trait::async_trait;

/// Ru
/// Репозиторий для взаимодействия с кэшем.
/// Определяет методы для установки, получения, удаления
/// и управления записями кэша, а также обработки политик
/// удаления. Реализаторы этого трейта должны быть потокобезопасными (Send + Sync).
/// Eng
/// Trait defining the interface for a cache repository.
/// This trait includes methods for setting, getting, deleting,
/// and managing cache entries, as well as handling eviction policies.
/// Implementors of this trait must be thread-safe (Send + Sync).
#[async_trait]
pub trait CacheRepository: Send + Sync {
    /// Sets a value in the cache with a specified TTL (time to live).
    /// # Arguments
    /// * `key` - The key under which the value is stored.
    /// * `value` - The value to be stored.
    /// * `ttl` - The time to live for the cache entry in seconds.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn set(&self, key: &str, value: &[u8], ttl: u64) -> Result<()>;

    /// Sets multiple values in the cache with specified TTLs.
    /// # Arguments
    /// * `items` - A vector of tuples containing key, value, and ttl.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn mset(&self, items: Vec<(&str, &[u8], u64)>) -> Result<()>;

    /// Gets a value from the cache by key.
    /// # Arguments
    /// * `key` - The key of the value to retrieve.
    /// # Returns
    /// * `Result<Option<Vec<u8>>>` - Ok(Some(value)) if found,
    async fn get(&self, key: &str) -> Result<Option<Vec<u8>>>;

    /// Gets multiple values from the cache by keys.
    /// # Arguments
    /// * `keys` - A vector of keys to retrieve.
    /// # Returns
    /// * `Result<Vec<Option<Vec<u8>>>>` - Ok(vector of values) corresponding to the keys.
    async fn mget(&self, keys: Vec<&str>) -> Result<Vec<Option<Vec<u8>>>>;

    /// Deletes a value from the cache by key.
    /// # Arguments
    /// * `key` - The key of the value to delete.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn delete(&self, key: &str) -> Result<()>;

    /// Checks if a key exists in the cache.
    /// # Arguments
    /// * `key` - The key to check for existence.
    /// # Returns
    /// * `Result<bool>` - Ok(true) if the key exists, Ok(false)
    async fn exists(&self, key: &str) -> Result<bool>;

    /// Clears all entries in the cache.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn clear(&self) -> Result<()>;

    /// Gets the time to live (TTL) for a cache entry by key.
    /// # Arguments
    /// * `key` - The key of the cache entry.
    /// # Returns
    /// * `Result<Option<u64>>` - Ok(Some(ttl)) if the key
    async fn ttl(&self, key: &str) -> Result<Option<u64>>;

    /// Gets the length of the value stored at a key.
    /// # Arguments
    /// * `key` - The key of the cache entry.
    /// # Returns
    /// * `Result<usize>` - Ok(length) if the key exists, Err otherwise
    async fn len(&self, key: &str) -> Result<usize>;

    /// Retrieves all keys matching a given pattern.
    /// # Arguments
    /// * `pattern` - The pattern to match keys against.
    /// # Returns
    /// * `Result<Vec<String>>` - Ok(vector of matching keys).
    async fn keys(&self, pattern: &str) -> Result<Vec<String>>;

    /// Cleans up expired entries in the cache.
    /// # Returns
    /// * `Result<usize>` - Ok(number of cleaned entries).
    async fn cleanup_expired(&self) -> Result<usize>;

    /// Sets the eviction policy for the cache.
    /// # Arguments
    /// * `policy` - The eviction policy to set.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn set_eviction_policy(&self, policy: EvictionPolicy) -> Result<()>;

    /// Gets the current eviction policy of the cache.
    /// # Returns
    /// * `Result<EvictionPolicy>` - Ok(current eviction policy).   
    async fn get_eviction_policy(&self) -> Result<EvictionPolicy>;

    /// Renames a key in the cache.
    /// # Arguments
    /// * `old_key` - The current key name.
    /// * `new_key` - The new key name.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn rename(&self, old_key: &str, new_key: &str) -> Result<()>;

    /// Renews the TTL of a cache entry.
    /// # Arguments
    /// * `key` - The key of the cache entry.
    /// * `ttl` - The new time to live for the cache entry in seconds.
    /// # Returns
    /// * `Result<()>` - Ok if the operation was successful, Err otherwise.
    async fn renew_ttl(&self, key: &str, ttl: u64) -> Result<()>;
}
