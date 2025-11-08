use serde::Deserialize;
use std::fs;
use crate::core::application::contracts::errors::blob_error::BlobError;

/// Enum representing different storage types
#[derive(Debug, Deserialize, Clone)]
#[serde(rename_all = "lowercase")]
pub enum StorageType {
    Local,
    S3,
}

/// Main configuration struct for blob storage
#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    /// Server mode: "rest" or "grpc"
    pub server_mode: String,
    
    /// Server address (e.g., "
    pub server_address: String,
    
    /// Type of storage backend: "local" or "s3"
    pub storage_type: StorageType,

    /// Allowed MIME types for blobs
    #[serde(default)]
    pub allowed_mime_types: Vec<String>,

    /// Maximum allowed blob size in bytes
    #[serde(default)]
    pub max_blob_size: u64,
    
    /// Allowed CORS origins
    #[serde(default)]
    pub cors_allowed_origins: Vec<String>,

    /// Allowed CORS methods
    #[serde(default)]
    pub cors_allowed_methods: Vec<String>,

    /// Allowed CORS headers
    #[serde(default)]
    pub cors_allowed_headers: Vec<String>,

    /// Automatic MIME type detection
    #[serde(default)]
    pub auto_mime_type_detection: bool,

    /// Allow anonymous access
    #[serde(default)]
    pub allow_anonymous_access: bool,

    /// Enable logging
    #[serde(default)]
    pub enable_logging: bool,

    /// Enable encryption
    #[serde(default)]
    pub enable_encryption: bool,

    /// Enable compression
    #[serde(default)]
    pub enable_compression: bool,
    
    /// Maximum number of concurrent connections
    #[serde(default = "default_max_connections")]
    pub max_concurrent_connections: usize,

    /// 32-byte base64-encoded encryption key (optional)
    pub encryption_key_base64: Option<String>,

    /// Compression algorithm (optional)
    pub compression_algorithm: Option<String>,
    
    /// Compression level (optional)
    #[serde(default = "default_compression_level")]
    pub compression_level: u8,
    
    /// Encryption algorithm (optional)
    pub encryption_algorithm: Option<String>,

    /// Logging level
    #[serde(default = "default_logging_level")]
    pub logging_level: String,

    /// Local storage config
    pub storage: StorageConfig,

    /// S3 storage configuration (optional)
    pub s3: Option<S3Config>,

    /// JWT configuration (optional)
    pub jwt: Option<JwtConfig>,
}

/// Default values
fn default_max_connections() -> usize { 10 }
fn default_logging_level() -> String { "info".into() }
fn default_compression_level() -> u8 { 5 }

/// Local filesystem storage configuration
#[derive(Debug, Deserialize, Clone)]
pub struct StorageConfig {
    /// Base path for local storage
    pub base_path: String,
}

/// S3 storage configuration
#[derive(Debug, Deserialize, Clone)]
pub struct S3Config {
    pub bucket: String,
    pub region: String,
    pub access_key: String,
    pub secret_key: String,
}

/// JWT configuration
#[derive(Debug, Deserialize, Clone)]
pub struct JwtConfig {
    pub authority: String,
    pub audience: String,
    
    /// JWKS rotation interval in seconds
    #[serde(default = "default_rotation_interval_secs")]
    pub rotation_interval_secs: u64,
    
    /// JWKS auto-refresh interval in seconds
    #[serde(default = "default_auto_refresh_interval_secs")]
    pub auto_refresh_interval_secs: u64,
}

fn default_rotation_interval_secs() -> u64 { 3600 }
fn default_auto_refresh_interval_secs() -> u64 { 300 }

impl Config {
    /// Load configuration from TOML file
    pub fn load(path: &str) -> Result<Self, BlobError> {
        let content = fs::read_to_string(path)?;
        let config: Config = toml::from_str(&content)?;
        config.validate()?;
        Ok(config)
    }

    /// Validate configuration fields
    fn validate(&self) -> Result<(), BlobError> {
        match self.storage_type {
            StorageType::Local => {
                if self.storage.base_path.is_empty() {
                    return Err(BlobError::InvalidConfig("Local storage path cannot be empty".into()));
                }
            },
            StorageType::S3 => {
                if let Some(s3) = &self.s3 {
                    if s3.bucket.is_empty() || s3.region.is_empty() || s3.access_key.is_empty() || s3.secret_key.is_empty() {
                        return Err(BlobError::InvalidConfig("S3 configuration is incomplete".into()));
                    }
                } else {
                    return Err(BlobError::InvalidConfig("S3 configuration is missing".into()));
                }
            },
        }
        Ok(())
    }
}