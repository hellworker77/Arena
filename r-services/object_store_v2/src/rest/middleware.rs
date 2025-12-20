use crate::rest::state::AppState;
use axum::extract::State;
use axum::middleware::Next;
use axum::response::{IntoResponse, Response};
use http::{Method, StatusCode};

pub async fn reject_writes_when_not_ready(
    State(state): State<AppState>,
    req: axum::http::Request<axum::body::Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let method = req.method();

    let is_write = matches!(
        *method,
        Method::PUT | Method::POST | Method::DELETE | Method::PATCH
    );
    
    if is_write && !state.ready.load(std::sync::atomic::Ordering::Relaxed) {
        return Err(StatusCode::SERVICE_UNAVAILABLE)
    }
    
    Ok(next.run(req).await)
}
