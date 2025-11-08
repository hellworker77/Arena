use std::{sync::Arc, time::Duration};
use tokio::sync::RwLock;
use tokio::time::interval;
use crate::jwks::Jwks;

pub async fn start_jwks_refresher(jwks: Arc<RwLock<Jwks>>, authority: String) {
   tokio::spawn(async move {
       let mut ticker = interval(Duration::from_secs(10 * 60));
       loop{
           ticker.tick().await;
           match Jwks::fetch(&authority).await {
               Ok(new_jwks) => {
                   *jwks.write().await = new_jwks;
                   tracing::debug!("JWKS refreshed from {}", authority);
               }
               Err(err) => {
                   tracing::warn!("Failed to refresh JWKS from {}: {}", authority, err);
               }
           }
       }
   });
}