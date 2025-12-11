use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EncryptionKind {
    None,
    Aes256Gcm,
}