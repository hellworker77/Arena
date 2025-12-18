use crate::bootstrap::bootstrap;
use crate::rest::SharedStore;
use crate::shutdown::shutdown_signal;
use std::net::SocketAddr;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::broadcast;
use tokio::time::interval;

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
mod shutdown;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let store = bootstrap("./repo")?;
    let shared: SharedStore = Arc::new(tokio::sync::Mutex::new(store));

    // Readiness flag:
    // true  => accept traffic
    // false => draining / shutting down
    let ready = Arc::new(AtomicBool::new(true));

    // shutdown channel
    let (shutdown_tx, _) = broadcast::channel::<()>(1);

    // background maintenance
    {
        let shared = shared.clone();
        let mut shutdown_rx = shutdown_tx.subscribe();

        tokio::spawn(async move {
            let mut tick = interval(Duration::from_secs(300));

            loop {
                tokio::select! {
                    _ = tick.tick() => {
                        let mut s = shared.lock().await;
                        let _ = s.checkpoint();
                    }
                    _ = shutdown_rx.recv() => {
                        break;
                    }
                }
            }

            // final checkpoint
            let mut s = shared.lock().await;
            let _ = s.checkpoint();
            println!("background worker stopped");
        });
    }

    let app = rest::router_v1(shared, ready.clone());

    let addr: SocketAddr = "0.0.0.0:8080".parse()?;
    let listener = tokio::net::TcpListener::bind(addr).await?;

    println!("listening on http://{}", addr);

    axum::serve(listener, app)
        .with_graceful_shutdown(async move {
            shutdown_signal().await;

            // Stop accepting traffic
            ready.store(false, Ordering::Relaxed);

            // Notify background workers
            let _ = shutdown_tx.send(());
        })
        .await?;

    println!("server exited cleanly");
    Ok(())
}