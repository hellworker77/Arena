use crate::core::application::contracts::errors::blob_error::BlobError;
use crate::core::application::services::blob_service::BlobService;
use axum::Json;
use axum::body::{Body, Bytes};
use axum::extract::{Multipart, Path, Query, State};
use axum::http::StatusCode;
use axum::http::header::CONTENT_DISPOSITION;
use axum::response::{IntoResponse, Response};
use serde::Deserialize;
use std::sync::Arc;

fn handle_error(err: BlobError) -> (StatusCode, String) {
    let (status, msg) = match err {
        BlobError::Io(e) => (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()),
        BlobError::IoError(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
        BlobError::ConfigParse(e) => (StatusCode::BAD_REQUEST, e.to_string()),
        BlobError::InvalidConfig(msg) => (StatusCode::BAD_REQUEST, msg),
        BlobError::NotFound(msg) => (StatusCode::NOT_FOUND, msg),
        BlobError::Network(msg) => (StatusCode::BAD_GATEWAY, msg),
        BlobError::Storage(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
        BlobError::Unknown(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
    };
    (status, msg)
}

fn error_into_response(err: BlobError) -> Response {
    let (status, msg) = handle_error(err);

    (status, msg).into_response()
}

#[derive(Deserialize)]
pub struct VersionParam {
    version: Option<u32>,
}

#[derive(Deserialize)]
pub struct PrefixParam {
    prefix: Option<String>,
}

pub async fn upload_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path(key): Path<String>,
    mut multipart: Multipart,
) -> impl IntoResponse {
    while let Some(field) = multipart.next_field().await.unwrap() {
        let filename = field.file_name().map(|s| s.to_string());
        let data = field.bytes().await.unwrap();

        match service.upload_blob(&key, &data, filename.clone()).await {
            Ok(_) => return StatusCode::CREATED.into_response(),
            Err(err) => return error_into_response(err),
        }
    }
    StatusCode::BAD_REQUEST.into_response()
}

pub async fn download_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path(key): Path<String>,
    Query(query): Query<VersionParam>,
) -> impl IntoResponse {
    match service.download_blob(&key, query.version).await {
        Ok((data, filename)) => {
            let body = Bytes::from(data);
            let mut res = Response::new(Body::from(body));
            res.headers_mut().insert(
                CONTENT_DISPOSITION,
                format!("attachment; filename=\"{}\"", filename)
                    .parse()
                    .unwrap(),
            );
            res
        }
        Err(err) => error_into_response(err),
    }
}

pub async fn delete_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path(key): Path<String>,
    Query(query): Query<VersionParam>,
) -> impl IntoResponse {
    match service.delete_blob(&key, query.version).await {
        Ok(_) => StatusCode::NO_CONTENT.into_response(),
        Err(err) => error_into_response(err),
    }
}

pub async fn list_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Query(query): Query<PrefixParam>,
) -> impl IntoResponse {
    match service.list_blobs(query.prefix.as_deref()).await {
        Ok(list) => (StatusCode::OK, Json(list)).into_response(),
        Err(err) => error_into_response(err),
    }
}

pub async fn exists_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    match service.exist_blob(&key).await {
        Ok(exists) => (StatusCode::OK, Json(exists)).into_response(),
        Err(err) => error_into_response(err),
    }
}

pub async fn move_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path((from, to)): Path<(String, String)>,
    Query(query): Query<VersionParam>,
) -> impl IntoResponse {
    match service.move_blob(&from, &to, query.version).await {
        Ok(_) => StatusCode::CREATED.into_response(),
        Err(err) => error_into_response(err),
    }
}

pub async fn copy_handler(
    State(service): State<Arc<dyn BlobService + Send + Sync>>,
    Path((from, to)): Path<(String, String)>,
    Query(query): Query<VersionParam>,
) -> impl IntoResponse {
    match service.copy_blob(&from, &to, query.version).await {
        Ok(_) => StatusCode::CREATED.into_response(),
        Err(err) => error_into_response(err),
    }
}
