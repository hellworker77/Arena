use serde::Deserialize;

#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    pub server_config: ServerConfig,
    pub logging_config: LoggingConfig,
    pub jwt_config: JwtConfig,
    pub cache_config: CacheConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ServerConfig {
    /// Address and port the server listens on
    pub server_address: String,

    /// Server mode "rest" or "grpc"
    pub server_mode: String,

    /// Cors allowed origins
    pub cors_allowed_origins: Vec<String>,

    /// Cors allowed methods
    pub cors_allowed_methods: Vec<String>,

    /// Cors allowed headers
    pub  cors_allowed_headers: Vec<String>,

    /// Maximum number of concurrent connections
    pub max_concurrent_connections: usize,
}

#[derive(Debug, Deserialize, Clone)]
pub struct LoggingConfig {
    /// Enable or disable logging
    pub enable_logging: bool,

    /// Logging level (e.g., "info", "debug", "error")
    pub logging_level: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct JwtConfig {
    /// JWT authority (issuer)
    pub authority: String,

    /// JWT audience
    pub audience: String,

    /// JWT key rotation interval in seconds
    pub rotation_interval_secs: u64,

    /// JWT auto-refresh interval in seconds
    pub auto_refresh_interval_secs: u64,
}

#[derive(Debug, Deserialize, Clone)]
pub struct CacheConfig {
    /// Maximum size of individual cache items in bytes
    pub max_cache_size: usize,

    /// Maximum total size of the cache in bytes
    pub max_total_cache_size: usize,

    /// Expiration time for cache items in seconds
    pub cache_item_expiration_secs: u64,

    /// Interval for cleaning up expired cache items in seconds
    pub cleanup_interval_secs: u64,

    /// Path to the file where cache dumps are stored
    pub dump_file_path: String,
}
