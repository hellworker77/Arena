use domain::models::compression_kind::CompressionKind;
use domain::models::encryption_kind::EncryptionKind;
use uuid::Uuid;

#[derive(Debug, Clone)]
pub struct BlobMetadataBuilder {
    pub blob_id: Option<Uuid>,
    pub version: Option<u32>,
    pub key: Option<String>,
    pub blob_type: Option<String>,
    pub created_at_unix: Option<i64>,
    pub updated_at_unix: Option<i64>,
    pub size_original: Option<u64>,
    pub size_compressed: Option<u64>,
    pub size_encrypted: Option<u64>,
    pub pak_offset: Option<u64>,
    pub encryption_nonce: Option<[u8; 12]>,
    pub compression: Option<CompressionKind>,
    pub encryption: Option<EncryptionKind>,
    pub content_type: Option<String>,
    pub content_checksum_sha256: Option<String>,
    pub content_etag: Option<String>,
}

impl BlobMetadataBuilder {
    /// Creates a new instance of BlobMetadataBuilder with all fields set to None.
    /// #Returns
    /// * `Self` - Returns a new BlobMetadataBuilder instance
    pub fn new() -> Self {
        BlobMetadataBuilder {
            blob_id: None,
            version: None,
            key: None,
            blob_type: None,
            created_at_unix: None,
            updated_at_unix: None,
            size_original: None,
            size_compressed: None,
            size_encrypted: None,
            pak_offset: None,
            encryption_nonce: None,
            compression: None,
            encryption: None,
            content_type: None,
            content_checksum_sha256: None,
            content_etag: None,
        }
    }
    
    pub fn build(self) -> Result<domain::models::blob_metadata::BlobMetadata, String> {
        Ok(domain::models::blob_metadata::BlobMetadata {
            blob_id: self.blob_id.ok_or("blob_id is required")?,
            version: self.version.ok_or("version is required")?,
            key: self.key.ok_or("key is required")?,
            blob_type: self.blob_type.ok_or("blob_type is required")?,
            created_at_unix: self.created_at_unix.ok_or("created_at_unix is required")?,
            updated_at_unix: self.updated_at_unix.ok_or("updated_at_unix is required")?,
            size_original: self.size_original.ok_or("size_original is required")?,
            size_compressed: self.size_compressed.ok_or("size_compressed is required")?,
            size_encrypted: self.size_encrypted.ok_or("size_encrypted is required")?,
            pak_offset: self.pak_offset.ok_or("pak_offset is required")?,
            encryption_nonce: self.encryption_nonce.ok_or("encryption_nonce is required")?,
            compression: self.compression.ok_or("compression is required")?,
            encryption: self.encryption.ok_or("encryption is required")?,
            content_type: self.content_type.ok_or("content_type is required")?,
            content_checksum_sha256: self.content_checksum_sha256.ok_or("content_checksum_sha256 is required")?,
            content_etag: self.content_etag.ok_or("content_etag is required")?,
        })
    }
    
    /// Sets the blob_id field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `blob_id` - A Uuid representing the unique identifier for the blob entry
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the blob_id set
    pub fn with_blob_id(mut self, blob_id: Uuid) -> Self {
        self.blob_id = Some(blob_id);
        self
    }

    /// Sets the version field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `version` - A u32 representing the blob version
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the version set
    pub fn with_version(mut self, version: u32) -> Self {
        self.version = Some(version);
        self
    }

    /// Sets the key field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `key` - A String representing the blob key provided in service
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the key set
    pub fn with_key(mut self, key: String) -> Self {
        self.key = Some(key);
        self
    }

    /// Sets the blob_type field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `blob_type` - A String representing the blob type from YAML schema
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the blob_type set
    pub fn with_blob_type(mut self, blob_type: String) -> Self {
        self.blob_type = Some(blob_type);
        self
    }

    /// Sets the created_at_unix field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `created_at` - An i64 representing the UNIX timestamp of creation
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the created_at_un
    pub fn with_created_at_unix(mut self, created_at: i64) -> Self {
        self.created_at_unix = Some(created_at);
        self
    }

    /// Sets the updated_at_unix field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `updated_at` - An i64 representing the UNIX timestamp of last update
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the updated_at_un
    pub fn with_updated_at_unix(mut self, updated_at: i64) -> Self {
        self.updated_at_unix = Some(updated_at);
        self
    }

    /// Sets the size_original field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `size` - A u64 representing the original size in bytes
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the size_original set
    pub fn with_size_original(mut self, size: u64) -> Self {
        self.size_original = Some(size);
        self
    }

    /// Sets the size_compressed field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `size` - A u64 representing the compressed size in bytes
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the size_compressed
    pub fn with_size_compressed(mut self, size: u64) -> Self {
        self.size_compressed = Some(size);
        self
    }

    /// Sets the size_encrypted field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `size` - A u64 representing the encrypted size in bytes
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the size_encrypted
    pub fn with_size_encrypted(mut self, size: u64) -> Self {
        self.size_encrypted = Some(size);
        self
    }

    /// Sets the pak_offset field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `offset` - A u64 representing the offset in the PAK file
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the pak_offset set
    pub fn with_pak_offset(mut self, offset: u64) -> Self {
        self.pak_offset = Some(offset);
        self
    }

    /// Sets the encryption_nonce field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `nonce` - A [u8; 12] array representing the NONCE for AES_GCM
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the encryption_nonce
    pub fn with_encryption_nonce(mut self, nonce: [u8; 12]) -> Self {
        self.encryption_nonce = Some(nonce);
        self
    }

    /// Sets the encryption field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `encryption` - An EncryptionKind representing the encryption type
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the encryption set
    pub fn with_compression(mut self, compression: CompressionKind) -> Self {
        self.compression = Some(compression);
        self
    }
    
    /// Sets the encryption field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `encryption` - An EncryptionKind representing the encryption type
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the encryption set
    pub fn with_encryption(mut self, encryption: EncryptionKind) -> Self {
        self.encryption = Some(encryption);
        self
    }

    /// Sets the encryption field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `encryption` - An EncryptionKind representing the encryption type
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the encryption set
    pub fn with_content_type(mut self, content_type: String) -> Self {
        self.content_type = Some(content_type);
        self
    }

    /// Sets the content_checksum_sha256 field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `content_checksum_sha256` - A String representing the SHA256 checksum of the content
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the content_checksum_sha
    pub fn with_content_checksum_sha256(mut self, content_checksum_sha256: String) -> Self {
        self.content_checksum_sha256 = Some(content_checksum_sha256);
        self
    }

    /// Sets the content_etag field of the BlobMetadataBuilder.
    /// #Arguments
    /// * `content_etag` - A String representing the ETag of the content
    /// #Returns
    /// * `Self` - Returns the updated BlobMetadataBuilder instance with the content_etag
    pub fn with_content_etag(mut self, content_etag: String) -> Self {
        self.content_etag = Some(content_etag);
        self
    }
}