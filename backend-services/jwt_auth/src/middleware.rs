use crate::jwks::Jwks;
use crate::jwks_cache::JwksCache;
use axum::body::Body;
use axum::http::StatusCode;
use axum::{extract::State, http::Request, middleware::Next, response::Response};
use jsonwebtoken::{decode, Algorithm, Validation};
use std::sync::Arc;
use tokio::sync::RwLock;
use tracing::{error, info, warn};

#[derive(Clone)]
pub struct JwtState {
    pub jwks_cache: JwksCache,
    pub audience: String,
    pub issuer: String,
}
#[derive(Debug, serde::Deserialize, Clone)]
pub struct Claims {
    pub sub: String,
    pub aud: String,
    pub iss: String,
    pub exp: usize,
    pub jti: String,
}

pub async fn jwt_auth(
    State(state): State<JwtState>,
    req: Request<Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let auth_header = req
        .headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok())
        .filter(|h| h.starts_with("Bearer "))
        .ok_or(StatusCode::UNAUTHORIZED)?;

    let token = &auth_header[7..];
    let header = jsonwebtoken::decode_header(token).map_err(|_| StatusCode::UNAUTHORIZED)?;
    let kid = header.kid.ok_or(StatusCode::UNAUTHORIZED)?;

    let decoding_key = state
        .jwks_cache
        .get(&kid)
        .await
        .ok_or(StatusCode::UNAUTHORIZED)?;

    let mut validation = Validation::new(Algorithm::RS256);
    validation.set_audience(&[&state.audience]);
    validation.set_issuer(&[&state.issuer]);

    match decode::<Claims>(token, &decoding_key, &validation) {
        Ok(data) => {
            info!("JWT validated for sub={}", data.claims.sub);
            Ok(next.run(req).await)
        }
        Err(e) => {
            warn!("JWT decode/validation failed: {:?}", e);
            Err(StatusCode::UNAUTHORIZED)
        }
    }
}