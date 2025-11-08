use thiserror::Error;

#[derive(Error, Debug)]
pub enum BlobError {

    /// An error local file system operations.
    #[error("IO Error: {0}")]
    Io(#[from] std::io::Error),

    /// An error local file system operations with custom message.
    #[error("IO Error: {0}")]
    IoError(String),

    /// An parsing error.
    #[error("Config Parse Error: {0}")]
    ConfigParse(#[from] toml::de::Error),

    /// Wrong configuration provided.
    #[error("Invalid Configuration: {0}")]
    InvalidConfig(String),

    /// An error while working with blob
    #[error("Blob not found: {0}")]
    NotFound(String),

    /// An error while working with network
    #[error("Network Error: {0}")]
    Network(String),
    
    #[error("Storage Error: {0}")]
    Storage(String),

    /// An unknown error
    #[error("Unknown Error: {0}")]
    Unknown(String),
}

impl From<serde_json::Error> for BlobError {
    fn from(err: serde_json::Error) -> Self {
        BlobError::IoError(err.to_string())
    }

}