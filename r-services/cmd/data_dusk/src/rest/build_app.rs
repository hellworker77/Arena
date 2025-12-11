use crate::rest::endpoints::{
    delete_handler, download_handler, exists_handler, list_handler, metadata_handler,
    upload_handler,
};
use crate::types::app_state::AppState;
use axum::http::StatusCode;
use axum::routing::{delete, get, post};
use axum::Router;
use infrastructure::factory::service_factory::ServiceFactory;

pub async fn build_app() -> Router {
    let factory = ServiceFactory::new();
    let state = AppState { factory };

    let api_v1: Router = Router::new()
        .route("/upload/{key}", post(upload_handler))
        .route("/download/{key}", get(download_handler))
        .route("/metadata/{key}", get(metadata_handler))
        .route("/exists/{key}", get(exists_handler))
        .route("/list", get(list_handler))
        .route("/delete/{key}", delete(delete_handler))
        .with_state(state);

    let app = Router::new()
        .nest("/api/v1", api_v1)
        .fallback(|| async { StatusCode::NOT_FOUND });

    app
}
