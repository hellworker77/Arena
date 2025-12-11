use async_trait::async_trait;

/// English
/// Service for encoding and decoding binary large objects (blobs).
/// Defines methods for transforming blob data.
/// Implementors can use various encoding schemes.
/// All methods return results for error handling.
///
/// Русский
/// Сервис для кодирования и декодирования бинарных больших объектов (блобов).
/// Определяет методы для преобразования данных блобов.
/// Реализаторы могут использовать различные схемы кодирования.
/// Все методы возвращают результаты для обработки ошибок.

#[async_trait]
pub trait BlobCodec {
    /// Get the encryption key.
    /// # Returns
    /// * `&[u8; 32]` - The encryption key used for encoding and decoding. 
    fn get_key(&self) -> &[u8; 32];
    
    /// Encode the given data.
    /// # Arguments
    /// * `data` - The binary data to encode.
    ///
    /// # Returns
    /// * `Ok(Vec<u8>)` with the encoded data.
    fn encode(&self, data: &[u8]) -> anyhow::Result<(Vec<u8>, [u8; 12], usize, usize)>;

    /// Decode the given data.
    /// # Arguments
    /// * `data` - The binary data to decode.
    ///
    /// # Returns
    /// * `Ok(Vec<u8>)` with the decoded data.
    fn decode(&self, nonce: &[u8; 12], ciphertext: &[u8]) -> anyhow::Result<Vec<u8>>;
}