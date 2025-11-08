use crate::core::application::contracts::data::blob_info::{BlobInfo, BlobVersionInfo};
use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::errors::blob_error::BlobError;
use crate::core::application::repositories::blob_repository::BlobRepository;
use async_trait::async_trait;
use std::collections::HashMap;
use std::path::PathBuf;
use tokio::fs as async_fs;
use tokio::io::AsyncReadExt;

#[derive(Clone)]
pub struct LocalRepository {
    pub base_path: PathBuf
}

impl LocalRepository {
    pub fn new(base_path: impl Into<PathBuf>) -> Self {
        Self {
            base_path: base_path.into(),
        }
    }

    fn key_to_path(&self, key: &str) -> PathBuf {
        self.base_path.join(key)
    }
}

#[async_trait]
impl BlobRepository for LocalRepository {
    async fn put(&self, key: &str, data: &[u8]) -> Result<(), BlobError> {
        let path = self.key_to_path(key);
        if let Some(parent) = path.parent() {
            async_fs::create_dir_all(parent).await.map_err(BlobError::Io)?;
        }
        async_fs::write(&path, data).await.map_err(BlobError::Io)?;
        Ok(())
    }

    async fn get(&self, key: &str) -> Result<Vec<u8>, BlobError> {
        let path = self.key_to_path(key);
        let mut file = async_fs::File::open(&path)
            .await
            .map_err(|_| BlobError::NotFound(key.to_string()))?;

        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer).await.map_err(BlobError::Io)?;
        Ok(buffer)
    }

    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata, BlobError> {
        let meta_key = format!("{}.meta.json", key);
        let data = self.get(&meta_key).await?;
        let meta: BlobMetadata = serde_json::from_slice(&data).map_err(|e| BlobError::InvalidConfig(format!("Invalid metadata JSON: {:?}", e)))?;
        Ok(meta)
    }

    async fn delete(&self, key: &str) -> Result<(), BlobError> {
        let path = self.key_to_path(key);
        if path.exists() {
            async_fs::remove_file(path).await.map_err(BlobError::Io)?;
        }
        Ok(())
    }

    async fn exists(&self, key: &str) -> Result<bool, BlobError> {
        Ok(self.key_to_path(key).exists())
    }

    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobInfo>, BlobError> {
        let mut result = HashMap::new();

        let mut entries = async_fs::read_dir(&self.base_path).await.map_err(BlobError::Io)?;
        while let Some(entry) = entries.next_entry().await.map_err(BlobError::Io)? {
            let path = entry.path();
            if path.is_file() {
                let key = path.strip_prefix(&self.base_path)
                    .unwrap()
                    .to_str()
                    .unwrap()
                    .to_string();

                if let Some(prefix) = prefix {
                    if !key.starts_with(prefix) {
                        continue;
                    }
                }

                let blob_info = result.entry(key.clone())
                    .or_insert_with(|| BlobInfo { versions: HashMap::new() });

                let file_name = path.file_name().unwrap().to_str().unwrap().to_string();
                blob_info.versions.insert(file_name, BlobVersionInfo { size: async_fs::metadata(&path).await.map_err(BlobError::Io)?.len() });
            }
        }

        Ok(result)
    }

    async fn copy(&self, from: &str, to: &str) -> Result<(), BlobError> {
        let from_path = self.key_to_path(from);
        let to_path = self.key_to_path(to);
        if let Some(parent) = to_path.parent() {
            async_fs::create_dir_all(parent).await.map_err(BlobError::Io)?;
        }
        async_fs::copy(&from_path, &to_path).await.map_err(BlobError::Io)?;
        Ok(())
    }

    async fn r#move(&self, from: &str, to: &str) -> Result<(), BlobError> {
        self.copy(from, to).await?;
        self.delete(from).await
    }
}