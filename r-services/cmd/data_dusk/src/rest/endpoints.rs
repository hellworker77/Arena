use axum::body::{Body, Bytes};
use axum::extract::{Path, Query, State};
use axum::http::{HeaderValue, StatusCode};
use axum::Json;
use axum::response::{IntoResponse, Response};
use application::services::blob_service::BlobService;
use crate::get_config::get_config;
use crate::types::app_state::AppState;
use crate::types::list_query::ListQuery;

pub async fn upload_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
    body: Bytes
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.upload(&key, &body).await {
        Ok(_) => StatusCode::CREATED.into_response(),
        Err(e) => {
            eprintln!("upload err: {:?}", e);
            (StatusCode::INTERNAL_SERVER_ERROR, format!("error: {}", e)).into_response()
        }
    }
}

pub async fn download_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.download(&key).await {
        Ok(bytes) => {
            let mime = mime_guess::from_path(&key)
                .first()
                .map(|m| m.essence_str().to_string())
                .or_else(|| infer::get(&bytes).map(|t| t.mime_type().to_string()))
                .unwrap_or_else(|| "application/octet-stream".to_string());

            println!("{}", mime);

            let body = Body::from(bytes);

            let mut resp = axum::response::Response::new(body);
            *resp.status_mut() = StatusCode::OK;

            resp.headers_mut().insert(
                axum::http::header::CONTENT_TYPE,
                axum::http::HeaderValue::from_static("application/octet-stream"),
            );

            resp.into_response()
        }
        Err(e) => {
            eprintln!("download err: {:?}", e);
            (
                StatusCode::NOT_FOUND,
                format!("error: {}", e)
            ).into_response()
        }
    }
}

pub async fn metadata_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.get_metadata(&key).await {
        Ok(meta) => (StatusCode::OK, Json(meta)).into_response(),
        Err(e) => {
            eprintln!("metadata err: {:?}", e);
            (StatusCode::NOT_FOUND, format!("error: {}", e)).into_response()
        }
    }
}

pub async fn exists_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.exists(&key).await {
        Ok(exists) => (StatusCode::OK, Json(exists)).into_response(),
        Err(e) => {
            eprintln!("exists err: {:?}", e);
            (StatusCode::INTERNAL_SERVER_ERROR, format!("error: {}", e)).into_response()
        }
    }
}

pub async fn list_handler(
    State(state): State<AppState>,
    Query(query): Query<ListQuery>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.list(query.prefix.as_deref()).await {
        Ok(listing) => (StatusCode::OK, Json(listing)).into_response(),
        Err(e) => {
            eprintln!("list err: {:?}", e);
            (StatusCode::INTERNAL_SERVER_ERROR, format!("error: {}", e)).into_response()
        }
    }
}

pub async fn delete_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.delete(&key).await {
        Ok(_) => StatusCode::NO_CONTENT.into_response(),
        Err(e) => {
            eprintln!("delete err: {:?}", e);
            (StatusCode::INTERNAL_SERVER_ERROR, format!("error: {}", e)).into_response()
        }
    }
}

pub async fn head_info_handler(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> impl IntoResponse {
    let (base_dir, codec_key) = get_config();
    let service = state.factory.get_or_init(base_dir, codec_key).await;

    match service.get_metadata(&key).await {
        Ok(meta) => {
            let mut resp = Response::new(Body::empty());

            resp.headers_mut().insert(
                axum::http::header::CONTENT_TYPE,
                HeaderValue::from_str(&meta.content_type).unwrap(),
            );

            resp.headers_mut().insert(
                "X-Blob-Version",
                HeaderValue::from_str(&meta.version.to_string()).unwrap(),
            );

            resp.headers_mut().insert(
                "X-Blob-Size",
                HeaderValue::from_str(&meta.size_original.to_string()).unwrap(),
            );

            resp.headers_mut().insert(
                "X-Blob-SHA256",
                HeaderValue::from_str(&meta.content_checksum_sha256).unwrap(),
            );

            resp.headers_mut().insert(
                "X-Blob-ETag",
                HeaderValue::from_str(&meta.content_etag).unwrap(),
            );

            *resp.status_mut() = StatusCode::OK;
            resp
        }
        Err(_) => StatusCode::NOT_FOUND.into_response(),
    }
}