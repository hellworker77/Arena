use axum::{
    extract::State,
    http::{header, HeaderValue},
    middleware::Next,
    response::Response,
};

use super::state::AppState;


/// Adds `Connection: close` during draining to discourage keep-alive.
///
/// Notes:
/// - Effective for HTTP/1.1 clients and load balancers.
/// - HTTP/2 does not use `Connection` header the same way, but keeping the behavior
///   consistent is still useful for mixed client stacks.
pub async fn add_connection_close_when_draining(
    State(state): State<AppState>,
    req: axum::http::Request<axum::body::Body>,
    next: Next
) -> Response {
    let mut resp = next.run(req).await;

    if !state.ready.load(std::sync::atomic::Ordering::Relaxed) {
        resp.headers_mut()
            .insert(header::CONNECTION, HeaderValue::from_static("close"));
    }

    resp
}