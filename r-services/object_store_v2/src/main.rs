use std::net::SocketAddr;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::broadcast;
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

async fn shutdown_signal() {
    #[cfg(unix)]
    {
        use tokio::signal::unix::{signal, SignalKind};

        let mut sigterm = signal(SignalKind::terminate())
            .expect("failed to install SIGTERM handler");

        tokio::select! {
            _ = tokio::signal::ctrl_c() => {}
            _ = sigterm.recv() => {}
        }
    }

    #[cfg(not(unix))]
    {
        tokio::signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    }

    println!("shutdown signal received");

    println!("shutdown signal received");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let store = bootstrap("./repo")?;
    let shared: SharedStore = Arc::new(tokio::sync::Mutex::new(store));

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

    let app = rest::router_v1(shared);

    let addr: SocketAddr = "0.0.0.0:8080".parse()?;
    let listener = tokio::net::TcpListener::bind(addr).await?;

    println!("listening on http://{}", addr);

    axum::serve(listener, app)
        .with_graceful_shutdown(async move {
            shutdown_signal().await;
            let _ = shutdown_tx.send(()); // notify workers
        })
        .await?;

    println!("server exited cleanly");
    Ok(())
}