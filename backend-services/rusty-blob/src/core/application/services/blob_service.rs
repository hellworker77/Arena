use std::collections::HashMap;
use async_trait::async_trait;
use crate::core::application::contracts::data::blob_info::BlobInfo;
use crate::core::application::contracts::errors::blob_error::BlobError;

/// BlobService defines high-level operations for managing blobs.
///
/// Implementations of this trait should handle business logic, validation, and coordinate
/// with underlying BlobRepository implementations.
/// All methods are async to support network or file IO efficiently.
#[async_trait]
pub trait BlobService: Send + Sync {
    /// Uploads a blob with the given key and data.
    ///
    /// # Arguments
    /// * `key` - A string slice that holds the key for the blob.
    /// * `data` - A byte slice that holds the blob data.
    /// * `original_filename` - An optional string that holds the original filename of the blob
    ///
    /// # Returns
    /// * `Result<(), BlobError>` - Returns Ok(()) if the upload is successful, otherwise returns a BlobError.
    async fn upload_blob(&self, key: &str, data: &[u8], original_filename: Option<String>) -> Result<(), BlobError>;

    /// Downloads a blob with the given key and optional version.
    ///
    /// # Arguments
    /// `key` - A string slice that holds the key for the blob.
    /// * `version` - An optional u32 that holds the version of the blob.
    ///
    /// # Returns
    /// * `Result<(Vec<u8>, String), BlobError>` - Returns a tuple containing the blob data and original filename if successful, otherwise returns a BlobError.
    async fn download_blob(&self, key: &str, version: Option<u32>) -> Result<(Vec<u8>, String), BlobError>;

    /// Deletes a blob with the given key and optional version.
    ///
    /// # Arguments
    /// * `key` - A string slice that holds the key for the blob.
    /// * `version` - An optional u32 that holds the version of the blob.
    ///
    /// # Returns
    /// * `Result<(), BlobError>` - Returns Ok(()) if the deletion is successful, otherwise returns a BlobError.
    async fn delete_blob(&self, key: &str, version: Option<u32>) -> Result<(), BlobError>;

    /// Lists blobs with an optional prefix.
    ///
    /// # Arguments
    /// * `prefix` - An optional string slice that holds the prefix for filtering blobs.
    ///
    /// # Returns
    /// * `Result<HashMap<String, BlobInfo>, BlobError>` - Returns a hashmap of blob keys and their info if successful, otherwise returns a BlobError.
    async fn list_blobs(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobInfo>, BlobError>;

    /// Checks if a blob with the given key exists.
    ///
    /// # Arguments
    /// * `key` - A string slice that holds the key for the blob.
    /// * `version` - An optional u32 that holds the version of the blob.
    ///
    /// # Returns
    /// * `Result<bool, BlobError>` - Returns true if the blob exists, false otherwise. Returns a BlobError if an error occurs.
    async fn exist_blob(&self, key: &str) -> Result<bool, BlobError>;

    /// Copies a blob from one key to another, with an optional version.
    ///
    /// # Arguments
    /// * `from` - A string slice that holds the source key for the blob.
    /// * `to` - A string slice that holds the destination key for the blob.
    /// * `version` - An optional u32 that holds the version of the blob to copy.
    ///
    /// # Returns
    /// * `Result<(), BlobError>` - Returns Ok(()) if the copy is successful, otherwise returns a BlobError.
    async fn copy_blob(&self, from: &str, to: &str, version: Option<u32>) -> Result<(), BlobError>;

    /// Moves a blob from one key to another, with an optional version.
    ///
    /// # Arguments
    /// * `from` - A string slice that holds the source key for the blob.
    /// * `to` - A string slice that holds the destination key for the blob.
    /// * `version` - An optional u32 that holds the version of the blob to
    ///
    /// # Returns
    /// * `Result<(), BlobError>` - Returns Ok(()) if the move is successful, otherwise returns a BlobError.
    async fn move_blob(&self, from: &str, to: &str, version: Option<u32>) -> Result<(), BlobError>;
}