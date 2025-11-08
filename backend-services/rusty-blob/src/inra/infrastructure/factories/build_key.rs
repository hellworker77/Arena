use base64::Engine;
use base64::engine::general_purpose;
use rand::RngCore;
use rand::rngs::OsRng;

pub fn build_key(src: Option<String>) -> Option<Vec<u8>> {
    match src.as_ref() {
        Some(k) => {
            let decoded = general_purpose::STANDARD.decode(k)
                .expect("Invalid base64 key");
            if decoded.len() != 32 {
                panic!(
                    "Encryption key must be 32 bytes for AES-256, got {} bytes (base64: {})",
                    decoded.len(),
                    k
                );
            }
            Some(decoded)
        },
        None => {
            let mut key = vec![0u8; 32];
            OsRng.fill_bytes(&mut key);
            Some(key)
        }
    }
}