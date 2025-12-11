use aes_gcm::{Aes256Gcm, Key, KeyInit, Nonce};
use anyhow::Result;
use application::codec::blob_codec::BlobCodec;
use flate2::{read::ZlibDecoder, write::ZlibEncoder, Compression};
use std::io::Write;
use aes_gcm::aead::Aead;
use rand::rngs::OsRng;
use rand_core::RngCore;

/// An example implementation of BlobCodec using zlib compression and AES encryption.
#[derive(Debug)]
pub struct ZlibAesCodec {
    pub key: [u8; 32],
}

impl ZlibAesCodec {
    pub fn new(key: [u8; 32]) -> Self {
        Self { key }
    }
}

impl BlobCodec for ZlibAesCodec {
    fn get_key(&self) -> &[u8; 32] {
        &self.key
    }

    fn encode(&self, data: &[u8]) -> Result<(Vec<u8>, [u8; 12], usize, usize)> {
        let mut encoder = ZlibEncoder::new(Vec::new(), Compression::default());
        encoder.write_all(data)?;
        let compressed = encoder.finish()?;
        let compressed_len = compressed.len();

        let key = Key::<Aes256Gcm>::from_slice(&self.key);
        let cipher = Aes256Gcm::new(key);

        let mut nonce = [0u8; 12];
        OsRng.fill_bytes(&mut nonce);
        let nonce_aead = Nonce::from_slice(&nonce);

        let ciphertext = cipher.encrypt(nonce_aead, compressed.as_ref())?;
        let encrypted_len = ciphertext.len();

        Ok((ciphertext, nonce, compressed_len, encrypted_len))
    }

    fn decode(&self, nonce: &[u8; 12], ciphertext: &[u8]) -> Result<Vec<u8>> {
        let key = Key::<Aes256Gcm>::from_slice(&self.key);
        let cipher = Aes256Gcm::new(key);
        let nonce_aead = Nonce::from_slice(nonce);

        let compressed = cipher.decrypt(nonce_aead, ciphertext)?;
        // decompress
        let mut decoder = ZlibDecoder::new(&compressed[..]);
        let mut out = Vec::new();
        std::io::copy(&mut decoder, &mut out)?;
        Ok(out)
    }
}