use crate::models::blob_metadata::BlobMetadata;
use anyhow::Result;
use async_trait::async_trait;
use std::collections::HashMap;
use std::sync::Arc;

/// Ru
/// Репозиторий для взаимодействия с блоб-хранилищем.
/// Определяет основные операции для блоб-хранилища.
/// Реализаторы могут использовать локальную файловую систему,
/// S3-совместимое хранилище или любую другую систему хранения.
/// Все методы асинхронные для эффективной поддержки сетевого
/// или файлового ввода-вывода.
/// Eng
/// Repository for interacting with a blob storage system.
/// Defines core operations for blob storage.
/// Implementors can use local file systems,
/// S3-compatible storage, or any other storage system.
/// All methods are asynchronous to efficiently support network
/// or file I/O.
#[async_trait]
pub trait BlobRepository: Send + Sync {
    /// Store a blob with the given key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// * `content` - The content to store.
    /// # Returns
    /// * `Ok(())` if successful.
    async fn put(&self, key: &str, content: &[u8]) -> Result<()>;

    /// Retrieve a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// # Returns
    /// * `Ok(Vec<u8>)` with the blob content if found.
    async fn get(&self, key: &str) -> Result<Vec<u8>>;

    /// Retrieve metadata for a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// # Returns
    /// * `Ok(BlobMetadata)` with the blob metadata if found.
    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata>;

    /// Delete a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// # Returns
    /// * `Ok(())` if deletion succeeds.
    async fn delete(&self, key: &str) -> Result<()>;

    /// Check if a blob exists.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// # Returns
    /// * `Ok(bool)` indicating existence.
    async fn exists(&self, key: &str) -> Result<bool>;

    /// List blobs with an optional prefix.
    /// # Arguments
    /// * `prefix` - Optional prefix to filter blobs.
    /// # Returns
    /// * `Ok(HashMap<String, BlobMetadata>)` mapping keys to their BlobMetadata
    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobMetadata>>;

    /// Copy a blob from source_key to destination_key.
    /// # Arguments
    /// * `source_key` - The key of the source blob.
    /// * `destination_key` - The key of the destination blob.
    /// # Returns
    /// * `Ok(())` if copy succeeds.
    async fn copy(&self, source_key: &str, destination_key: &str) -> Result<()>;

    /// Move a blob from source_key to destination_key.
    /// # Arguments
    /// * `source_key` - The key of the source blob.
    /// * `destination_key` - The key of the destination blob.
    /// # Returns
    /// * `Ok(())` if move succeeds.
    async fn r#move(&self, source_key: &str, destination_key: &str) -> Result<()>;
}

/// Blanket implementation of BlobRepository for Arc-wrapped types.
#[async_trait]
impl<T: BlobRepository + Send + Sync> BlobRepository for Arc<T> {
    async fn put(&self, key: &str, content: &[u8]) -> anyhow::Result<()> {
        (**self).put(key, content).await
    }

    async fn get(&self, key: &str) -> anyhow::Result<Vec<u8>> {
        (**self).get(key).await
    }

    async fn get_metadata(&self, key: &str) -> anyhow::Result<BlobMetadata> {
        (**self).get_metadata(key).await
    }

    async fn delete(&self, key: &str) -> anyhow::Result<()> {
        (**self).delete(key).await
    }

    async fn exists(&self, key: &str) -> anyhow::Result<bool> {
        (**self).exists(key).await
    }

    async fn list(&self, prefix: Option<&str>) -> anyhow::Result<HashMap<String, BlobMetadata>> {
        (**self).list(prefix).await
    }

    async fn copy(&self, source_key: &str, destination_key: &str) -> anyhow::Result<()> {
        (**self).copy(source_key, destination_key).await
    }

    async fn r#move(&self, source_key: &str, destination_key: &str) -> anyhow::Result<()> {
        (**self).r#move(source_key, destination_key).await
    }
}