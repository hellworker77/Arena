use anyhow::Result;
use application::repository::blob_repository::BlobRepository;
use application::services::blob_service::BlobService;
use async_trait::async_trait;
use std::collections::HashMap;
use domain::models::blob_metadata::BlobMetadata;

/// Implementation of the BlobService trait.
pub struct BlobServiceImpl<Repo: BlobRepository> {
    repository: Repo,
}

impl<Repo: BlobRepository> BlobServiceImpl<Repo> {
    /// Create a new BlobServiceImpl with the given repository.
    pub fn new(repository: Repo ) -> Self {
        Self { repository }
    }

    /// Validate the blob key.
    /// # Arguments
    /// * `key` - The blob key to validate.
    ///
    /// # Returns
    /// * `Ok(())` if the key is valid.
    /// * `Err` if the key is invalid.
    fn validate_key(&self, key: &str) -> Result<()> {
        if key.trim().is_empty() {
            anyhow::bail!("empty key");
        }

        if key.contains("..") {
            anyhow::bail!("invalid key: contains ..");
        }

        Ok(())
    }

    /// Auto generate metadata with incremented version.
    /// # Arguments
    /// * `key` - The blob key.
    /// * `size` - The size of the blob.
    /// * `content_type` - The content type of the blob.
    ///
    /// # Returns
    /// * `Ok(())` if metadata is updated successfully.
    async fn update_metadata(&self, key: &str, size: u64, content_type: Option<String>) -> Result<()> {
        todo!()
    }
}

#[async_trait]
impl<Repo: BlobRepository + Send + Sync> BlobService for BlobServiceImpl<Repo> {
    async fn upload(&self, key: &str, payload_bytes: &[u8]) -> Result<()> {
        self.validate_key(key)?;
        self.repository.put(key, payload_bytes).await?;
        Ok(())
    }

    async fn download(&self, key: &str) -> Result<Vec<u8>> {
        let data = self.repository.get(key).await?;
        Ok(data)
    }

    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata> {
        let metadata = self.repository.get_metadata(key).await?;
        Ok(metadata)
    }

    async fn delete(&self, key: &str) -> Result<()> {
        self.repository.delete(key).await?;
        Ok(())
    }

    async fn exists(&self, key: &str) -> Result<bool> {
        let exists = self.repository.exists(key).await?;
        Ok(exists)
    }

    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobMetadata>> {
        let blobs = self.repository.list(prefix).await?;
        Ok(blobs)
    }

    async fn copy(&self, from: &str, to: &str) -> Result<()> {
        self.repository.copy(from, to).await?;
        Ok(())
    }

    async fn r#move(&self, from: &str, to: &str) -> Result<()> {
        self.repository.r#move(from, to).await?;
        Ok(())
    }
}