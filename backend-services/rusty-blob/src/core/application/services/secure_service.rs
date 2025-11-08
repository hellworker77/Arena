use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::errors::blob_error::BlobError;

pub trait SecureService {
    /// Processes data for secure storage (e.g., encryption, compression).
    ///
    /// # Arguments
    /// * `data` - A byte slice that holds the data to be processed.
    ///
    /// # Returns
    /// * `Result<Vec<u8>, BlobError>` - Returns the processed data if successful, otherwise returns a BlobError.
    fn process_in(&self, data: &[u8]) -> Result<Vec<u8>, BlobError>;

    /// Processes data retrieved from secure storage (e.g., decryption, decompression).
    ///
    /// # Arguments
    /// * `data` - A byte slice that holds the data to be processed.
    /// * `metadata` - Metadata associated with the blob.
    ///
    /// # Returns
    /// * `Result<Vec<u8>, BlobError>` - Returns the processed data if successful, otherwise returns a BlobError.
    fn process_get(&self, data: &[u8], metadata: BlobMetadata) -> Result<Vec<u8>, BlobError>;
}