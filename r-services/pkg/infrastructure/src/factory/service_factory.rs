use crate::services::blob_service_impl::BlobServiceImpl;
use persistence::repository::pak_repository::PakRepository;
use std::sync::Arc;
use tokio::sync::Mutex;

#[derive(Clone)]
pub struct ServiceFactory {
    // for simplicity: single shared service instance
    // (manages many .pak files under base_dir)
    inner: Arc<Mutex<Option<Arc<BlobServiceImpl<Arc<PakRepository>>>>>>,
}

impl ServiceFactory {
    pub fn new() -> Self {
        Self {
            inner: Arc::new(Mutex::new(None)),
        }
    }

    pub async fn get_or_init(
        &self,
        base_dir: impl Into<std::path::PathBuf>,
        codec_key: [u8; 32],
    ) -> Arc<BlobServiceImpl<Arc<PakRepository>>> {
        let mut lock = self.inner.lock().await;
        if let Some(svc) = &*lock {
            return svc.clone();
        }
        
        //Create a PakRepository that will store many .pak files in base_dir
        let repo = PakRepository::new(base_dir.into(), codec_key).expect("pak repo create");
        let repo = Arc::new(repo);
        
        //Create service and wrap in Arc
        let service = Arc::new(BlobServiceImpl::new(repo.clone()));
        *lock = Some(service.clone());
        service
    }
}
