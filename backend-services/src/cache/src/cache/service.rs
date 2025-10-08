use crate::cache::model::CacheStore;
use crate::proto::cache_server::{Cache, CacheServer};
use crate::proto::{
    ClearRequest, ClearResponse, DeleteRequest, DeleteResponse, GetRequest, GetResponse,
    SetRequest, SetResponse, UpdateRequest, UpdateResponse,
};
use tonic::{Request, Response, Status};

#[derive(Clone)]
pub struct CacheService {
    pub store: CacheStore,
}

impl CacheService {
    pub fn new() -> Self {
        Self {
            store: CacheStore::default(),
        }
    }

    pub fn into_server(self) -> CacheServer<Self> {
        CacheServer::new(self)
    }
}

#[tonic::async_trait]
impl Cache for CacheService {
    async fn get(&self, req: Request<GetRequest>) -> Result<Response<GetResponse>, Status> {
        let key = req.into_inner().key;
        let value = self.store.get(&key);

        Ok(Response::new(GetResponse {
            found: value.is_some(),
            value: value.unwrap_or_default(),
        }))
    }

    async fn set(&self, request: Request<SetRequest>) -> Result<Response<SetResponse>, Status> {
        let req = request.into_inner();

        let ttl = if req.ttl_milliseconds == 0 {
            60000
        } else {
            req.ttl_milliseconds
        };
        self.store.insert(req.key, req.value, ttl);

        Ok(Response::new(SetResponse { saved: true }))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        let key = request.into_inner().key;
        let removed = self.store.delete(&key);

        Ok(Response::new(DeleteResponse { removed }))
    }

    async fn update(
        &self,
        request: Request<UpdateRequest>,
    ) -> Result<Response<UpdateResponse>, Status> {
        let req = request.into_inner();

        let ttl = if req.ttl_milliseconds == 0 {
            60000
        } else {
            req.ttl_milliseconds
        };
        let mut updated = false;
        if self.store.get(&req.key).is_some() {
            self.store.insert(req.key, req.value, ttl);
            updated = true;
        }
        Ok(Response::new(UpdateResponse { updated }))
    }

    async fn clear(
        &self,
        _request: Request<ClearRequest>,
    ) -> Result<Response<ClearResponse>, Status> {
        self.store.clear();
        Ok(Response::new(ClearResponse { cleared: true }))
    }
}
