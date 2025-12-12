use crate::models::compression_kind::CompressionKind;
use crate::models::encryption_kind::EncryptionKind;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlobMetadata {
    /// Unique identifier for the blob entry
    pub blob_id: Uuid,

    /// Blob version each upload increments version
    pub version: u32,

    /// Blob key provided in service
    pub key: String,

    /// Blob types (from YAML schema)
    pub blob_type: String,

    /// UNIX timestamps
    pub created_at_unix: i64,
    pub updated_at_unix: i64,

    /// Data sizes in bytes
    pub size_original: u64,
    pub size_compressed: u64,
    pub size_encrypted: u64,

    /// Offset in the PAK file
    pub pak_offset: u64,

    /// NONCE for AES_GCM (12 bytes)
    pub encryption_nonce: [u8; 12],

    /// Compression
    pub compression: CompressionKind,

    /// Encryption
    pub encryption: EncryptionKind,

    pub content_type: Option<String>,
    pub content_checksum_sha256: Option<String>,
}