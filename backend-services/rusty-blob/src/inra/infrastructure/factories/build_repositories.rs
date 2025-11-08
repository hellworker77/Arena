use crate::core::application::contracts::data::config::{Config, StorageType};
use crate::core::application::repositories::blob_repository::BlobRepository;
use crate::inra::persistence::repositories::local_repository::LocalRepository;
use std::sync::Arc;

pub fn build_repositories(config: &Config) -> Arc<dyn BlobRepository> {
    let base_repository: Arc<dyn BlobRepository> = match config.storage_type {
        StorageType::Local => {
            let local = LocalRepository::new(&config.storage.base_path);
            Arc::new(local)
        }
        StorageType::S3 => {
            todo!()
        }
    };

    base_repository
}