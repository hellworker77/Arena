use crate::jwks::Jwks;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::sync::RwLock;
use tracing::{info, warn};

#[derive(Clone)]
pub struct JwksCache {
    inner: Arc<RwLock<HashMap<String, (jsonwebtoken::DecodingKey, Instant)>>>,
    authority: String,
    jwks_rotation_interval_secs: u64,
    jwks_auto_refresh_interval_secs: u64,
}

impl JwksCache {
    pub async fn new(
        authority: String,
        jwks_rotation_interval_secs: Option<u64>,
        jwks_auto_refresh_interval_secs: Option<u64>,
    ) -> anyhow::Result<Self> {
        let cache = Self {
            inner: Arc::new(RwLock::new(HashMap::new())),
            authority: authority.clone(),
            jwks_rotation_interval_secs: jwks_rotation_interval_secs.unwrap_or(3600),
            jwks_auto_refresh_interval_secs: jwks_auto_refresh_interval_secs.unwrap_or(300),
        };

        cache.refresh().await?;
        cache.start_auto_refresh();
        Ok(cache)
    }

    /// Updates the JWKS cache by fetching the latest keys from the authority.
    pub async fn refresh(&self) -> anyhow::Result<()> {
        match Jwks::fetch(&self.authority).await {
            Ok(jwks) => {
                let mut map = self.inner.write().await;

                for key in jwks.keys {
                    match jsonwebtoken::DecodingKey::from_rsa_components(&key.n, &key.e) {
                        Ok(decoding_key) => {
                            map.insert(key.kid.clone(), (decoding_key, Instant::now()));
                        }
                        Err(e) => warn!("Failed to parse key {}: {:?}", key.kid, e),
                    }
                }

                info!("✅ JWKS refreshed successfully from {}", self.authority);
                Ok(())
            }
            Err(e) => {
                warn!("⚠️ Failed to refresh JWKS from {}: {:?}", self.authority, e);
                Err(e.into())
            }
        }
    }

    /// Get key by kid
    pub async fn get(&self, kid: &str) -> Option<jsonwebtoken::DecodingKey> {
        let map = self.inner.read().await;
        map.get(kid).map(|(key, _)| key.clone())
    }

    /// Remove expired keys and refresh periodically
    async fn cleanup(&self) {
        let ttl = Duration::from_secs(self.jwks_rotation_interval_secs);
        let mut map = self.inner.write().await;
        let before = map.len();
        map.retain(|_, (_, inserted)| inserted.elapsed() < ttl);
        let after = map.len();

        if after < before {
            info!("🧹 Cleaned up {} expired JWKS keys", before - after);
        }
    }

    /// Start automatic refresh task in the background
    fn start_auto_refresh(&self) {
        let this = self.clone();
        tokio::spawn(async move {
            info!(
                "🚀 Starting JWKS auto-refresh every {} seconds",
                this.jwks_auto_refresh_interval_secs
            );

            let mut interval =
                tokio::time::interval(Duration::from_secs(this.jwks_auto_refresh_interval_secs));

            loop {
                interval.tick().await;

                if let Err(e) = this.refresh().await {
                    warn!("JWKS auto-refresh failed: {:?}", e);
                }
                
                this.cleanup().await;
            }
        });
    }
}
