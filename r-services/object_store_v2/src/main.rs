use std::{
    net::SocketAddr,
    sync::{
        Arc,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};

use hyper::server::conn::http1;
use hyper_util::rt::TokioIo;
use hyper_util::service::TowerToHyperService;
use tokio::sync::broadcast;
use tokio::time::interval;
use tower::ServiceExt;
use crate::bootstrap::bootstrap;
use crate::rest::SharedStore;
use crate::shutdown::shutdown_signal;

mod bootstrap;
mod compaction;
mod error;
mod gc;
mod index;
mod manifest;
mod recovery;
mod rest;
mod segment;
mod shutdown;
mod store;
mod wal;

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
        });
    }

    let app = rest::router_v1(shared, ready.clone());

    let addr: SocketAddr = "0.0.0.0:8080".parse()?;
    let listener = tokio::net::TcpListener::bind(addr).await?;
    println!("listening on http://{}", addr);

    // Shutdown gate for accept loop
    let (accept_stop_tx, mut accept_stop_rx) = tokio::sync::oneshot::channel::<()>();

    // Accept loop: stop accepting new TCP connections once shutdown begins.
    let server_task = tokio::spawn(async move {
        loop {
            tokio::select! {
                _ = &mut accept_stop_rx => {
                    // Stop accepting new connections.
                    break;
                }
                res = listener.accept() => {
                    let (stream, _peer) = match res {
                        Ok(x) => x,
                        Err(_) => continue,
                    };

                    let io = TokioIo::new(stream);
                    let svc = app.clone().into_service();
                    let hyper_svc = TowerToHyperService::new(svc);

                    // HTTP/1 connection builder:
                    // - Disable keep-alive for tighter draining behavior.
                    // - This affects only new connections; existing ones will finish in-flight requests.
                    tokio::spawn(async move {
                        let conn = http1::Builder::new()
                            .keep_alive(false) // limit keep-alive
                            .serve_connection(io, hyper_svc);

                        let _ = conn.await;
                    });
                }
            }
        }
    });

    // Wait for shutdown signal
    shutdown_signal().await;
    println!("shutdown signal received");

    // Enter draining mode: stop writes + mark not ready
    ready.store(false, Ordering::Relaxed);

    // Notify background workers
    let _ = shutdown_tx.send(());

    // Stop accepting new connections
    let _ = accept_stop_tx.send(());

    // Wait for accept loop to stop
    let _ = server_task.await;

    println!("server exited cleanly");
    Ok(())
}