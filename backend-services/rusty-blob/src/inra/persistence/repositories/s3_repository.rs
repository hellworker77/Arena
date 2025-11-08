use crate::core::application::contracts::data::blob_info::BlobInfo;
use crate::core::application::contracts::data::blob_metadata::BlobMetadata;
use crate::core::application::contracts::errors::blob_error::BlobError;
use crate::core::application::repositories::blob_repository::BlobRepository;
use crate::infrastructure::persistence::storages::s3::S3;
use async_trait::async_trait;
use std::collections::HashMap;
use std::path::Path;

pub struct S3Repository {
    pub(crate) client: Client,
    pub(crate) bucket: String,
}

impl S3Repository {
    pub async fn new(config: &S3Config) -> Result<Self, BlobError> {
        let shared_config = aws_config::defaults(BehaviorVersion::latest())
            .region(aws_sdk_s3::config::Region::new(config.region.clone()))
            .load()
            .await;

        let client = Client::new(&shared_config);
        Ok(Self {
            client,
            bucket: config.bucket.clone(),
        })
    }
}

#[async_trait]
impl BlobRepository for S3Repository {

}