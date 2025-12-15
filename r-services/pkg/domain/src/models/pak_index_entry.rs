use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Index entry persisted in the INDEX section (binary)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PakIndexEntry {
    /// Blob key
    pub key: String,

    /// Version number
    pub version: u32,

    /// Offset in PAK file
    pub offset: u64,

    /// Original size
    pub size_original: u64,

    /// Compressed size
    pub size_compressed: u64,
    
    /// Nonce for encryption
    /// 12 bytes for AES-GCM
    pub nonce: [u8; 12],
    
    /// Unique identifier for the blob
    pub blob_id: Uuid
}