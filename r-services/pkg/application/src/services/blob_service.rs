use std::collections::HashMap;
use anyhow::Result;
use async_trait::async_trait;
use crate::models::blob_metadata::BlobMetadata;

/// English
/// Service for interacting with a blob storage system.
/// Defines core operations for blob storage.
/// Implementors can use local file systems.
/// All methods are asynchronous to efficiently support I/O.

/// Русский
/// Сервис для взаимодействия с блоб-хранилищем.
/// Определяет основные операции для блоб-хранилища.
/// Реализаторы могут использовать локальную файловую систему.
/// Все методы асинхронные для эффективной поддержки ввода-вывода.
#[async_trait]
pub trait BlobService {
    /// Upload a blob with the given key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// * `content` - The content to store.
    ///
    /// # Returns
    /// * `Ok(())` if successful.
    async fn upload(&self, key: &str, content: &[u8]) -> Result<()>;

    /// Download a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(Vec<u8>)` with the blob content if found.
    async fn download(&self, key: &str) -> Result<Vec<u8>>;

    /// Retrieve metadata for a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(BlobMetadata)` with the blob metadata if found.
    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata>;

    /// Delete a blob by its key.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(())` if deletion succeeds.
    async fn delete(&self, key: &str) -> Result<()>;

    /// Check if a blob exists.
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(bool)` indicating existence.
    async fn exists(&self, key: &str) -> Result<bool>;

    /// List blobs with an optional prefix.
    /// # Arguments
    /// * `prefix` - Optional prefix to filter blobs.
    ///
    /// # Returns
    /// * `Ok(HashMap<String, Vec<BlobMetadata>>)` with the list of blobs.
    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobMetadata>>;

    /// Copy a blob from source key to destination key.
    /// # Arguments
    /// * `from` - The key of the source blob.
    /// * `to` - The key of the destination blob.
    ///
    /// # Returns
    /// * `Ok(())` if copy succeeds.
    async fn copy(&self, from: &str, to: &str) -> Result<()>;

    /// Move a blob from source key to destination key.
    /// # Arguments
    /// * `from` - The key of the source blob.
    /// * `to` - The key of the destination blob.
    ///
    /// # Returns
    /// * `Ok(())` if move succeeds.
    async fn r#move(&self, from: &str, to: &str) -> Result<()>;
}