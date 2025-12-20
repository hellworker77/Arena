use axum::{
    extract::State,
    http::StatusCode,
    response::IntoResponse,
};
use std::sync::{
    atomic::{AtomicBool, Ordering},
    Arc,
};
use super::state::AppState;

/// Liveness probe.
///
/// Contract:
/// - Returns 200 as long as the process is running.
/// - Must never block.
/// - Must never perform I/O.
pub async fn livez() -> impl IntoResponse {
    StatusCode::OK
}

/// Readiness probe.
///
/// Contract:
/// - Returns 200 only if the service can accept traffic.
/// - Returns 503 during shutdown or failed initialization.
/// - Must be fast and non-blocking.
pub async fn readyz(
    State(state): State<AppState>,
) -> impl IntoResponse {
    if state.ready.load(Ordering::Relaxed) {
        StatusCode::OK
    } else {
        StatusCode::SERVICE_UNAVAILABLE
    }
}