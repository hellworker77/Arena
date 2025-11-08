use crate::core::application::contracts::data::config::Config;
use crate::core::application::services::blob_service::BlobService;
use crate::presentation::build_app::build_app;
use axum::serve;
use std::net::SocketAddr;
use std::sync::Arc;
use tokio::net::TcpListener;

// ---------- Server ----------
pub async fn startup(service: Arc<dyn BlobService + Send + Sync>, config: &Config) {
    let app = build_app(service, config).await;

    let addr: SocketAddr = config
        .server_address
        .parse()
        .expect("Invalid server address");

    let listener = TcpListener::bind(addr).await.unwrap();

    println!("🚀 Blob server running on http://{}", addr);

    serve(listener, app).await.expect("TODO: panic message");
}


