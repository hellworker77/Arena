#[derive(thiserror::Error, Debug)]
pub enum BrokerError {
    /// An IO error from sled database operations.
    #[error("IO error: {0}")]
    Io(#[from] sled::Error),

    /// An error during serialization or deserialization.
    #[error("Serde error: {0}")]
    Serde(#[from] serde_json::Error),

    /// An error indicating that a requested item was not found.
    #[error("Not found: {0}")]
    NotFound(String),

    /// An error indicating that a message could not be delivered.
    #[error("Other: {0}")]
    Other(String),
}