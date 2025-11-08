use crate::core::application::contracts::data::config::Config;
use crate::core::application::services::blob_service::BlobService;
use crate::presentation::endpoints::{
    copy_handler, delete_handler, download_handler, exists_handler, list_handler,
    move_handler, upload_handler,
};
use axum::http::{HeaderName, HeaderValue, Method, StatusCode};
use axum::routing::{delete, get, post, put};
use axum::{middleware, Router};
use jwt_auth::jwks_cache::JwksCache;
use jwt_auth::middleware::{jwt_auth, JwtState};
use std::sync::Arc;
use tower::limit::ConcurrencyLimitLayer;
use tower_http::cors::{AllowOrigin, CorsLayer};
use tower_http::trace::{DefaultMakeSpan, DefaultOnResponse, TraceLayer};
use tracing::Level;

/// Builds a CORS layer based on the provided configuration.
fn build_cors(config: &Config) -> CorsLayer {
    let mut cors = CorsLayer::new();

    // Origins
    if config.cors_allowed_origins.contains(&"*".to_string()) {
        cors = cors.allow_origin(AllowOrigin::any());
    } else {
        let origins: Vec<AllowOrigin> = config
            .cors_allowed_origins
            .iter()
            .map(|o| AllowOrigin::exact(HeaderValue::from_str(o).unwrap()))
            .collect();
        for origin in origins {
            cors = cors.allow_origin(origin);
        }
    }

    // Methods
    let methods: Vec<Method> = config
        .cors_allowed_methods
        .iter()
        .filter_map(|m| m.parse().ok())
        .collect();
    cors = cors.allow_methods(methods);

    // Headers
    let headers: Vec<HeaderName> = config
        .cors_allowed_headers
        .iter()
        .filter_map(|h| h.parse().ok())
        .collect();
    cors = cors.allow_headers(headers);

    cors
}

/// Applies middleware layers to the given router based on the configuration.
fn apply_middleware(router: Router, config: &Config) -> Router {
    let cors = build_cors(&config);
    let mut r = router.layer(cors);

    let level = config.logging_level.parse::<Level>().unwrap_or(Level::INFO);

    let trace_layer = TraceLayer::new_for_http()
        .make_span_with(DefaultMakeSpan::new().level(level))
        .on_response(DefaultOnResponse::new().level(level));

    let concurrency_limit = ConcurrencyLimitLayer::new(config.max_concurrent_connections);

    if config.enable_logging {
        r = r.layer(trace_layer);
    }
    if config.max_concurrent_connections > 0 {
        r = r.layer(concurrency_limit);
    }
    r
}

/// Builds the Axum application router with all routes and middleware.
pub async fn build_app(service: Arc<dyn BlobService + Send + Sync>, config: &Config) -> Router {
    let jwt_cfg = config.jwt.as_ref().expect("Jwt config must be set");

    let jwks_cache = JwksCache::new(
        jwt_cfg.authority.clone(),
        Some(jwt_cfg.rotation_interval_secs),
        Some(jwt_cfg.auto_refresh_interval_secs),
    )
    .await
    .expect("Failed to create JWKS cache");

    let jwt_state = JwtState {
        jwks_cache,
        audience: jwt_cfg.audience.clone(),
        issuer: jwt_cfg.authority.clone(),
    };

    let api_v1 = Router::new()
        .route("/upload/{key}", post(upload_handler))
        .route("/download/{key}", get(download_handler))
        .route("/delete/{key}", delete(delete_handler))
        .route("/list", get(list_handler))
        .route("/exists/{key}", get(exists_handler))
        .route("/move/{from}/{to}", put(move_handler))
        .route("/copy/{from}/{to}", post(copy_handler))
        .with_state(service.clone())
        .layer(middleware::from_fn_with_state(jwt_state.clone(), jwt_auth));

    let app = Router::new()
        .nest("/api/v1", api_v1)
        .fallback(|| async { StatusCode::NOT_FOUND });

    apply_middleware(app, &config)
}
