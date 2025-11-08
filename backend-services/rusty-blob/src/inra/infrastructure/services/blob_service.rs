use crate::core::application::contracts::data::blob_info::BlobInfo;
use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::data::config::Config;
use crate::core::application::contracts::errors::blob_error::BlobError;
use crate::core::application::repositories::blob_repository::BlobRepository;
use crate::core::application::services::blob_service::BlobService;
use crate::core::application::services::secure_service::SecureService;
use std::collections::HashMap;
use std::sync::Arc;
use async_trait::async_trait;

///Service layer for blob business logic
#[derive(Clone)]
pub struct BlobServiceImpl {
    repository: Arc<dyn BlobRepository + Send + Sync>,
    secure_service: Arc<dyn SecureService + Send + Sync>,
    config: Config,
}

impl BlobServiceImpl {
    pub fn new(
        repository: Arc<dyn BlobRepository + Send + Sync>,
        secure_service: Arc<dyn SecureService + Send + Sync>,
        config: &Config,
    ) -> Self {
        Self {
            repository,
            secure_service,
            config: config.clone(),
        }
    }

    fn validate_mime_type(&self, data: &[u8]) -> Result<(), BlobError> {
        if let Some(kind) = infer::get(data) {
            let mime_type = kind.mime_type();
            if !self.config.allowed_mime_types.contains(&mime_type.to_string()) {
                return Err(BlobError::InvalidConfig(format!(
                    "Mime type {} is not allowed",
                    mime_type
                )));
            }
        } else if !self.config.allowed_mime_types.contains(&"*".to_string()) {
            return Err(BlobError::InvalidConfig("Cannot detect mime type".to_string()));
        }
        Ok(())
    }

    async fn get_next_version(&self, key: &str) -> Result<u32, BlobError> {
        let all = self.repository.list(Some(key)).await?;
        let mut max_version = 0u32;

        if let Some(info) = all.get(key) {
            for version_name in info.versions.keys() {
                if version_name.ends_with(".meta.json") {
                    continue;
                }

                if let Some(v) = version_name.strip_prefix("v") {
                    if let Ok(num) = v.parse::<u32>() {
                        max_version = max_version.max(num);
                    }
                }
            }
        }

        Ok(max_version + 1)
    }

    fn meta_key(&self, key: &str, version: u32) -> String {
        format!("{}/v{}.meta.json", key, version)
    }

    fn latest_meta_key(&self, key: &str) -> String {
        format!("{}/latest.meta.json", key)
    }
}

#[async_trait]
impl BlobService for BlobServiceImpl {
    async fn upload_blob(
        &self,
        key: &str,
        data: &[u8],
        original_filename: Option<String>,
    ) -> Result<(), BlobError> {
        if data.len() > self.config.max_blob_size as usize {
            return Err(BlobError::InvalidConfig(format!(
                "Blob size exceeds maximum allowed size of {} bytes",
                self.config.max_blob_size
            )));
        }

        if self.config.auto_mime_type_detection {
            self.validate_mime_type(data)?;
        }

        let version = self.get_next_version(key).await?;
        let version_key = format!("{}/v{}", key, version);
        let latest_key = format!("{}/latest", key);

        let processed_data = self.secure_service.process_in(data)?;

        self.repository.put(&version_key, &processed_data).await?;

        let meta = BlobMetadata::new(
            key,
            version,
            processed_data.len() as u64,
            infer::get(data).map(|k| k.mime_type().to_string()),
            original_filename.clone(),
            self.config.compression_algorithm.clone(),
            self.config.encryption_key_base64.clone(),
            self.config.encryption_algorithm.clone(),
        );

        let meta_json = serde_json::to_vec_pretty(&meta)?;
        self.repository.put(&self.meta_key(key, version), &meta_json).await?;

        if self.repository.exists(&latest_key).await? {
            self.repository.delete(&latest_key).await?;
        }
        self.repository.copy(&version_key, &latest_key).await?;
        self.repository.put(&self.latest_meta_key(key), &meta_json).await?;

        Ok(())
    }

    async fn download_blob(&self, key: &str, version: Option<u32>) -> Result<(Vec<u8>, String), BlobError> {
        let blob_key = match version {
            Some(v) => format!("{}/v{}", key, v),
            None => format!("{}/latest", key),
        };

        let data = self.repository.get(&blob_key).await?;
        let meta = self.repository.get_metadata(&blob_key).await?;

        let decrypted = self.secure_service.process_get(&data, meta.clone())?;
        let filename = meta.original_filename.unwrap_or_else(|| key.to_string());

        Ok((decrypted, filename))
    }

    async fn delete_blob(&self, key: &str, version: Option<u32>) -> Result<(), BlobError> {
        let blob_key = match version {
            Some(v) => format!("{}/v{}", key, v),
            None => format!("{}/latest", key),
        };
        self.repository.delete(&blob_key).await
    }

    async fn list_blobs(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobInfo>, BlobError> {
        let all = self.repository.list(prefix).await?;
        let mut result = HashMap::new();

        for (key, info) in all {
            let mut filtered_versions = HashMap::new();
            for (vname, vinfo) in info.versions {
                if !vname.ends_with(".meta.json") {
                    filtered_versions.insert(vname, vinfo);
                }
            }
            result.insert(key, BlobInfo { versions: filtered_versions });
        }

        Ok(result)
    }

    async fn exist_blob(&self, key: &str) -> Result<bool, BlobError> {
        let latest_key = format!("{}/latest", key);
        self.repository.exists(&latest_key).await
    }

    async fn copy_blob(&self, from: &str, to: &str, version: Option<u32>) -> Result<(), BlobError> {
        let src_key = match version {
            Some(v) => format!("{}/v{}", from, v),
            None => format!("{}/latest", from),
        };
        self.repository.copy(&src_key, to).await
    }

    async fn move_blob(&self, from: &str, to: &str, version: Option<u32>) -> Result<(), BlobError> {
        let src_key = match version {
            Some(v) => format!("{}/v{}", from, v),
            None => format!("{}/latest", from),
        };
        self.repository.r#move(&src_key, to).await
    }
}
