use crate::core::application::contracts::data::config::Config;
use crate::core::application::repositories::blob_repository::BlobRepository;
use crate::core::application::services::blob_service::BlobService;
use crate::inra::infrastructure::factories::build_key::build_key;
use crate::inra::infrastructure::services::blob_service::BlobServiceImpl;
use crate::inra::infrastructure::services::secure_service::SecureServiceImpl;
use std::sync::Arc;

pub fn build_services(
    repo: Arc<dyn BlobRepository + Send + Sync>,
    config: &Config,
) -> Arc<dyn BlobService> {
    let key = build_key(config.encryption_key_base64.clone());

    let secured_service = Arc::new(
        SecureServiceImpl::new(
            key,
            config.enable_compression,
            config.compression_algorithm.clone().unwrap_or_default(),
            config.enable_encryption,
            config.encryption_algorithm.clone().unwrap_or_default(),
            config.compression_level.clone(),
        ),
    );

    Arc::new(
        BlobServiceImpl::new(
            repo,
            secured_service,
            config,
        ),
    )
}