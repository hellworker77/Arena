use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CompressionKind {
    None,
    Zlib,
    Zstd,
}