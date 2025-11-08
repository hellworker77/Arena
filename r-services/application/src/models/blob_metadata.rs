use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BlobMetadata {
    /// Version of the blob
    pub version: u32,

    /// Key identifying the blob
    pub key: String,

    /// UNIX timestamp
    pub created_at: u64,

    /// Size of the blob in bytes
    pub size: u64,

    /// MIME type of the blob
    pub mime_type: Option<String>,

    /// Original filename of the blob
    pub original_filename: Option<String>,

    /// Compression algorithm used when storing the blob
    pub compression_algorithm: Option<String>,

    /// Encryption key used for the blob (if applicable)
    pub encryption_key: Option<String>,

    /// Encryption algorithm used when storing the blob
    pub encryption_algorithm: Option<String>,
}