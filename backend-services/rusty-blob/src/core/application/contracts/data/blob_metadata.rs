use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};

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

impl BlobMetadata {
    pub fn new(
        key: &str,
        version: u32,
        size: u64,
        mime_type: Option<String>,
        original_filename: Option<String>,
        compression_algorithm: Option<String>,
        encryption_key: Option<String>,
        encryption_algorithm: Option<String>,
    ) -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();

        Self {
            version,
            key: key.to_string(),
            size,
            mime_type,
            created_at: now,
            original_filename,
            compression_algorithm,
            encryption_key,
            encryption_algorithm,
        }
    }
}
