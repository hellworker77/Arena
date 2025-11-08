use std::collections::HashMap;
use std::path::{Path, PathBuf};
use async_trait::async_trait;
use crate::core::application::contracts::data::blob_info::BlobInfo;
use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::errors::blob_error::BlobError;

/// BlobRepository defines the core operations for a blob storage backend.
///
/// Implementations can be local file system, S3-compatible storage, or any other storage system.
/// All methods are async to support network or file IO efficiently.
#[async_trait]
pub trait BlobRepository: Send + Sync {
    /// Store a blob with the given key.
    ///
    /// # Arguments
    /// * `prefix_path` - The path prefix where the blob will be stored.
    /// * `data` - The content to store.
    ///
    /// # Returns
    /// * `Ok(())` if successful.
    /// * `Err(BlobError)` if storing fails.
    async fn put(&self, prefix_path: &str, data: &[u8]) -> Result<(), BlobError>;

    /// Retrieve a blob by its key.
    ///
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(Vec<u8>)` with the blob content if found.
    /// * `Err(BlobError::NotFound)` if the blob does not exist.
    /// * `Err(BlobError)` for other failures.
    async fn get(&self, key: &str) -> Result<Vec<u8>, BlobError>;
    
    /// Retrieve metadata for a blob by its key.
    /// 
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    /// 
    /// # Returns
    /// * `Ok(BlobMetadata)` with the blob metadata if found.
    /// * `Err(BlobError::NotFound)` if the blob does not exist.
    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata, BlobError>;

    /// Delete a blob by its key.
    ///
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(())` if deletion succeeds.
    /// * `Err(BlobError)` if deletion fails or blob does not exist.
    async fn delete(&self, key: &str) -> Result<(), BlobError>;

    /// Check if a blob exists.
    ///
    /// # Arguments
    /// * `key` - The unique identifier for the blob.
    ///
    /// # Returns
    /// * `Ok(true)` if the blob exists.
    /// * `Ok(false)` if the blob does not exist.
    /// * `Err(BlobError)` if the check fails.
    async fn exists(&self, key: &str) -> Result<bool, BlobError>;

    /// List all blobs optionally filtered by a prefix.
    ///
    /// # Arguments
    /// * `prefix` - Optional prefix to filter the keys.
    ///
    /// # Returns
    /// * `Ok(Vec<String>)` with all matching blob keys.
    /// * `Err(BlobError)` if listing fails.
    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobInfo>, BlobError>;

    /// Copy a blob from one key to another.
    ///
    /// # Arguments
    /// * `from` - The source blob key.
    /// * `to` - The destination blob key.
    ///
    /// # Returns
    /// * `Ok(())` if copy succeeds.
    /// * `Err(BlobError)` if copy fails.
    async fn copy(&self, from: &str, to: &str) -> Result<(), BlobError>;

    /// Move a blob from one key to another (copy + delete).
    ///
    /// # Arguments
    /// * `from` - The source blob key.
    /// * `to` - The destination blob key.
    ///
    /// # Returns
    /// * `Ok(())` if move succeeds.
    /// * `Err(BlobError)` if move fails.
    async fn r#move(&self, from: &str, to: &str) -> Result<(), BlobError>;
}