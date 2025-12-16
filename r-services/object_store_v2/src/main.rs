use std::net::SocketAddr;
use std::sync::Arc;
use std::time::Duration;
use tokio::time::interval;
use crate::bootstrap::bootstrap;
use crate::rest::SharedStore;

mod error;
mod bootstrap;
mod wal;
mod manifest;
mod segment;
mod index;
mod recovery;
mod store;
mod gc;
mod compaction;
mod rest;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let store = bootstrap("./repo")?;
    let shared: SharedStore = Arc::new(tokio::sync::Mutex::new(store));

    // background maintenance
    {
        let shared = shared.clone();
        tokio::spawn(async move {
            let mut tick = interval(Duration::from_secs(300));
            loop {
                tick.tick().await;
                let mut s = shared.lock().await;
                let _ = s.checkpoint();
            }
        });
    }

    let app = rest::router_v1(shared);

    let addr: SocketAddr = "0.0.0.0:8080".parse()?;
    println!("listening on http://{}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await?, app).await?;

    Ok(())
}