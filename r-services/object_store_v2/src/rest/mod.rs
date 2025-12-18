mod middleware;
pub mod probes;
mod state;

use crate::error::error::StoreError;
use crate::rest::probes::{livez, readyz};
use crate::rest::state::AppState;
use crate::store::ObjectStore;
use axum::body::to_bytes;
use axum::response::Response;
use axum::{
    Router,
    body::Body,
    extract::{Path, State},
    http::{HeaderMap, StatusCode, header},
    response::IntoResponse,
    routing::{delete, get, put},
};
use std::{
    ops::RangeInclusive,
    sync::{
        Arc,
        atomic::{AtomicU64, Ordering},
    },
};
use axum::middleware::from_fn_with_state;
use tokio::io::{AsyncReadExt, AsyncSeekExt};
use tokio_util::io::ReaderStream;
use crate::rest::middleware::reject_writes_when_not_ready;

pub type SharedStore = Arc<tokio::sync::Mutex<ObjectStore>>;

#[derive(Default)]
pub struct Metrics {
    pub put_requests: AtomicU64,
    pub get_requests: AtomicU64,
    pub delete_requests: AtomicU64,
    pub bytes_in: AtomicU64,
    pub bytes_out: AtomicU64,
    pub range_gets: AtomicU64,
    pub not_modified: AtomicU64,
    pub precondition_failed: AtomicU64,
}

pub fn router_v1(store: SharedStore, ready_flag: Arc<std::sync::atomic::AtomicBool>) -> Router {
    let state = AppState {
        store,
        metrics: Arc::new(Metrics::default()),
        ready: ready_flag,
    };

    Router::new()
        .route("/api/v1/health", get(health))
        .route(
            "/api/v1/objects/{key}",
            get(get_object)
                .head(head_object)
                .put(put_object)
                .delete(delete_object),
        )
        .route("/api/v1/livez", get(livez))
        .route("/api/v1/readyz", get(readyz))
        .route("/metrics", get(metrics))
        .layer(from_fn_with_state(
            state.clone(),
            reject_writes_when_not_ready,
        ))
        .with_state(state)
}

async fn health() -> &'static str {
    "ok"
}

// ---------- PUT (streaming to temp file, then store.put) ----------
// Честный компромисс: не держим всё в RAM, но store.put всё равно читает весь файл.
// Чтобы сделать fully-streaming CAS ingest — надо менять storage формат на chunked.
async fn put_object(
    State(st): State<AppState>,
    Path(key): Path<String>,
    body: Body,
) -> impl IntoResponse {
    st.metrics.put_requests.fetch_add(1, Ordering::Relaxed);

    let bytes = match to_bytes(body, usize::MAX).await {
        Ok(b) => b,
        Err(_) => return StatusCode::BAD_REQUEST.into_response(),
    };

    st.metrics
        .bytes_in
        .fetch_add(bytes.len() as u64, Ordering::Relaxed);

    let mut store = st.store.lock().await;
    match store.put(key, &bytes) {
        Ok(_) => StatusCode::CREATED.into_response(),
        Err(e) => map_error(e).into_response(),
    }
}

// ---------- GET with Range + ETag/If-Match ----------

async fn get_object(
    State(st): State<AppState>,
    Path(key): Path<String>,
    headers: HeaderMap,
) -> impl IntoResponse {
    st.metrics.get_requests.fetch_add(1, Ordering::Relaxed);

    // locate
    let (seg_path, payload_off, payload_len, etag) = {
        let store = st.store.lock().await;
        match store.locate_for_read(&key) {
            Ok(x) => x,
            Err(e) => return map_error(e).into_response(),
        }
    };

    // ETag preconditions
    if let Some(if_match) = headers.get(header::IF_MATCH) {
        if if_match.to_str().ok().map(|s| s != etag).unwrap_or(true) {
            st.metrics
                .precondition_failed
                .fetch_add(1, Ordering::Relaxed);
            return StatusCode::PRECONDITION_FAILED.into_response();
        }
    }
    if let Some(if_none) = headers.get(header::IF_NONE_MATCH) {
        if if_none.to_str().ok().map(|s| s == etag).unwrap_or(false) {
            st.metrics.not_modified.fetch_add(1, Ordering::Relaxed);
            return (StatusCode::NOT_MODIFIED, [(header::ETAG, etag)]).into_response();
        }
    }

    // parse Range
    let range = match headers.get(header::RANGE).and_then(|v| v.to_str().ok()) {
        Some(r) => match parse_range(r, payload_len) {
            Ok(rr) => {
                st.metrics.range_gets.fetch_add(1, Ordering::Relaxed);
                Some(rr)
            }
            Err(code) => return code.into_response(),
        },
        None => None,
    };

    // open file and stream only needed bytes
    let file = match tokio::fs::File::open(&seg_path).await {
        Ok(f) => f,
        Err(_) => return StatusCode::INTERNAL_SERVER_ERROR.into_response(),
    };

    let (status, start, end) = if let Some(r) = range {
        let start = *r.start();
        let end = *r.end();
        (StatusCode::PARTIAL_CONTENT, start, end)
    } else {
        (StatusCode::OK, 0u64, payload_len.saturating_sub(1))
    };

    let length = if payload_len == 0 { 0 } else { end - start + 1 };

    // Seek to payload start + range start
    let mut file = file;
    if file
        .seek(std::io::SeekFrom::Start(payload_off + start))
        .await
        .is_err()
    {
        return StatusCode::INTERNAL_SERVER_ERROR.into_response();
    }

    // Limit reader to `length` bytes
    let reader = file.take(length);
    let stream = ReaderStream::new(reader);
    st.metrics.bytes_out.fetch_add(length, Ordering::Relaxed);

    let mut resp = Response::builder().status(status);

    resp = resp.header(header::ETAG, etag);
    resp = resp.header(header::ACCEPT_RANGES, "bytes");
    resp = resp.header(header::CONTENT_LENGTH, length.to_string());

    if status == StatusCode::PARTIAL_CONTENT {
        resp = resp.header(
            header::CONTENT_RANGE,
            format!("bytes {}-{}/{}", start, end, payload_len),
        );
    }

    resp.body(Body::from_stream(stream)).unwrap()
}

// ---------- DELETE ----------

async fn delete_object(State(st): State<AppState>, Path(key): Path<String>) -> impl IntoResponse {
    st.metrics.delete_requests.fetch_add(1, Ordering::Relaxed);

    let mut store = st.store.lock().await;
    match store.delete(key) {
        Ok(_) => StatusCode::NO_CONTENT.into_response(),
        Err(e) => map_error(e).into_response(),
    }
}

// ---------- /metrics (Prometheus text) ----------

async fn metrics(State(st): State<AppState>) -> impl IntoResponse {
    let m = &st.metrics;
    let text = format!(
        "\
# TYPE store_put_requests_total counter
store_put_requests_total {}
# TYPE store_get_requests_total counter
store_get_requests_total {}
# TYPE store_delete_requests_total counter
store_delete_requests_total {}
# TYPE store_bytes_in_total counter
store_bytes_in_total {}
# TYPE store_bytes_out_total counter
store_bytes_out_total {}
# TYPE store_range_gets_total counter
store_range_gets_total {}
# TYPE store_not_modified_total counter
store_not_modified_total {}
# TYPE store_precondition_failed_total counter
store_precondition_failed_total {}
",
        m.put_requests.load(Ordering::Relaxed),
        m.get_requests.load(Ordering::Relaxed),
        m.delete_requests.load(Ordering::Relaxed),
        m.bytes_in.load(Ordering::Relaxed),
        m.bytes_out.load(Ordering::Relaxed),
        m.range_gets.load(Ordering::Relaxed),
        m.not_modified.load(Ordering::Relaxed),
        m.precondition_failed.load(Ordering::Relaxed),
    );

    (
        StatusCode::OK,
        [(header::CONTENT_TYPE, "text/plain; version=0.0.4")],
        text,
    )
}

//  -------- HEAD (like GET but no body) ----------

async fn head_object(
    State(st): State<AppState>,
    Path(key): Path<String>,
    headers: HeaderMap,
) -> impl IntoResponse {
    // locate object (no IO on segment body)
    let (_seg_path, _off, len, etag) = {
        let store = st.store.lock().await;
        match store.locate_for_read(&key) {
            Ok(x) => x,
            Err(e) => return map_error(e).into_response(),
        }
    };

    // If-Match
    if let Some(v) = headers.get(header::IF_MATCH) {
        if v.to_str().ok().map(|s| s != etag).unwrap_or(true) {
            return StatusCode::PRECONDITION_FAILED.into_response();
        }
    }

    // If-None-Match
    if let Some(v) = headers.get(header::IF_NONE_MATCH) {
        if v.to_str().ok().map(|s| s == etag).unwrap_or(false) {
            return StatusCode::NOT_MODIFIED.into_response();
        }
    }

    // Build HEAD response (nobody!)
    let resp = Response::builder()
        .status(StatusCode::OK)
        .header(header::ETAG, etag)
        .header(header::CONTENT_LENGTH, len.to_string())
        .header(header::ACCEPT_RANGES, "bytes");

    resp.body(Body::empty()).unwrap()
}

// ---------- helpers ----------

fn map_error(err: StoreError) -> StatusCode {
    match err {
        StoreError::NotFound => StatusCode::NOT_FOUND,
        StoreError::Deleted => StatusCode::GONE,
        StoreError::HashMismatch => StatusCode::CONFLICT,
        StoreError::CasMiss | StoreError::SegmentMissing => StatusCode::INTERNAL_SERVER_ERROR,
        _ => StatusCode::INTERNAL_SERVER_ERROR,
    }
}

// Range: "bytes=start-end" only (поддержка suffix можно добавить позже)
fn parse_range(s: &str, total: u64) -> Result<RangeInclusive<u64>, StatusCode> {
    if !s.starts_with("bytes=") {
        return Err(StatusCode::RANGE_NOT_SATISFIABLE);
    }
    let v = &s["bytes=".len()..];
    let (a, b) = v.split_once('-').ok_or(StatusCode::RANGE_NOT_SATISFIABLE)?;
    let start: u64 = a.parse().map_err(|_| StatusCode::RANGE_NOT_SATISFIABLE)?;
    let end: u64 = if b.is_empty() {
        total.saturating_sub(1)
    } else {
        b.parse().map_err(|_| StatusCode::RANGE_NOT_SATISFIABLE)?
    };

    if total == 0 || start >= total || end >= total || start > end {
        return Err(StatusCode::RANGE_NOT_SATISFIABLE);
    }
    Ok(start..=end)
}
