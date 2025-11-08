use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::errors::blob_error::BlobError;
use crate::core::application::services::secure_service::SecureService;
use crate::inra::infrastructure::factories::build_key::build_key;
use aes::Aes256;
use block_modes::block_padding::Pkcs7;
use block_modes::{BlockMode, Cbc};
use flate2::Compression;
use flate2::read::GzDecoder;
use flate2::write::GzEncoder;
use rand::RngCore;
use rand::rngs::OsRng;
use std::io::{Read, Write};

type Aes256Cbc = Cbc<Aes256, Pkcs7>;
pub struct SecureServiceImpl {
    /// AES-256 encryption key (32 bytes)
    encryption_key: Option<Vec<u8>>,

    /// Enable compression before storage
    enable_compression: bool,

    /// Compression algorithm (e.g., "gzip")
    compression_algorithm: String,

    /// Enable encryption before storage
    enable_encryption: bool,

    /// Encryption algorithm (e.g., "aes-256-cbc")
    encryption_algorithm: String,

    compression_level: u8,
}

impl SecureServiceImpl {
    pub fn new(
        encryption_key: Option<Vec<u8>>,
        enable_compression: bool,
        compression_algorithm: String,
        enable_encryption: bool,
        encryption_algorithm: String,
        compression_level: u8,
    ) -> Self {
        Self {
            encryption_key,
            enable_compression,
            compression_algorithm,
            enable_encryption,
            encryption_algorithm,
            compression_level,
        }
    }

    fn compress(&self, data: &[u8]) -> Result<Vec<u8>, BlobError> {
        match self.compression_algorithm.to_lowercase().as_str() {
            "gzip" => {
                let level = match self.compression_level {
                    1..=9 => self.compression_level,
                    _ => 5,
                };
                let mut encoder = GzEncoder::new(Vec::new(), Compression::new(level.into()));
                encoder.write_all(data).map_err(BlobError::Io)?;
                encoder.finish().map_err(BlobError::Io)
            }
            other => Err(BlobError::InvalidConfig(format!(
                "Unsupported compression algorithm: {}",
                other
            ))),
        }
    }

    fn decompress(&self, data: &[u8], algorithm: &str) -> Result<Vec<u8>, BlobError> {
        match algorithm.to_lowercase().as_str() {
            "gzip" => {
                let mut decoder = GzDecoder::new(data);
                let mut decompressed = Vec::new();
                decoder
                    .read_to_end(&mut decompressed)
                    .map_err(BlobError::Io)?;
                Ok(decompressed)
            }
            other => Err(BlobError::InvalidConfig(format!(
                "Unsupported compression algorithm: {}",
                other
            ))),
        }
    }

    fn encrypt(&self, data: &[u8]) -> Result<Vec<u8>, BlobError> {
        if !self.enable_encryption {
            return Ok(data.to_vec());
        }

        match self.encryption_algorithm.to_lowercase().as_str() {
            "aes-256-cbc" => self.encrypt_aes256_cbc(data),
            other => Err(BlobError::InvalidConfig(format!(
                "Unsupported encryption algorithm: {}",
                other
            ))),
        }
    }

    fn decrypt(
        &self,
        data: &[u8],
        algorithm: &str,
        key: Option<Vec<u8>>,
    ) -> Result<Vec<u8>, BlobError> {
        match algorithm.to_lowercase().as_str() {
            "aes-256-cbc" => self.decrypt_aes256_cbc(data, key),
            other => Err(BlobError::InvalidConfig(format!(
                "Unsupported encryption algorithm: {}",
                other
            ))),
        }
    }

    fn encrypt_aes256_cbc(&self, data: &[u8]) -> Result<Vec<u8>, BlobError> {
        let key = self
            .encryption_key
            .as_ref()
            .ok_or_else(|| BlobError::InvalidConfig("Encryption key missing".into()))?;

        if key.len() != 32 {
            return Err(BlobError::InvalidConfig(
                "Encryption key must be 32 bytes for AES-256".into(),
            ));
        }

        let mut iv = [0u8; 16];
        OsRng.fill_bytes(&mut iv);

        let cipher = Aes256Cbc::new_from_slices(key, &iv)
            .map_err(|e| BlobError::InvalidConfig(format!("AES init error: {:?}", e)))?;

        let mut ciphertext = cipher.encrypt_vec(data);

        let mut result = iv.to_vec();
        result.append(&mut ciphertext);
        Ok(result)
    }

    fn decrypt_aes256_cbc(
        &self,
        data: &[u8],
        encryption_key: Option<Vec<u8>>,
    ) -> Result<Vec<u8>, BlobError> {
        let key = encryption_key
            .as_ref()
            .or(self.encryption_key.as_ref())
            .ok_or_else(|| BlobError::InvalidConfig("Encryption key missing".into()))?;

        if key.len() != 32 {
            return Err(BlobError::InvalidConfig(
                "Encryption key must be 32 bytes for AES-256".into(),
            ));
        }

        if data.len() < 16 {
            return Err(BlobError::InvalidConfig(
                "Invalid encrypted data (too short)".into(),
            ));
        }

        let (iv, ciphertext) = data.split_at(16);

        let cipher = Aes256Cbc::new_from_slices(key, iv)
            .map_err(|e| BlobError::InvalidConfig(format!("AES init error: {:?}", e)))?;

        cipher
            .decrypt_vec(ciphertext)
            .map_err(|e| BlobError::InvalidConfig(format!("AES decrypt error: {:?}", e)))
    }
}

impl SecureService for SecureServiceImpl {
    fn process_in(&self, data: &[u8]) -> Result<Vec<u8>, BlobError> {
        let mut result = data.to_vec();

        if self.enable_compression {
            result = self.compress(&result)?;
        }
        if self.enable_encryption {
            result = self.encrypt(&result)?;
        }

        Ok(result)
    }

    fn process_get(&self, data: &[u8], metadata: BlobMetadata) -> Result<Vec<u8>, BlobError> {
        let mut result = data.to_vec();

        let key = build_key(metadata.encryption_key.clone());

        if let Some(ref alg) = metadata.encryption_algorithm {
            if !alg.is_empty() {
                result = self.decrypt(&result, alg, key)?;
            }
        }

        if let Some(ref alg) = metadata.compression_algorithm {
            if !alg.is_empty() {
                result = self.decompress(&result, alg)?;
            }
        }

        Ok(result)
    }
}
