use std::{
    net::SocketAddr,
    sync::{
        Arc,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};

use crate::bootstrap::bootstrap;
use crate::rest::SharedStore;
use crate::shutdown::shutdown_signal;
use hyper::server::conn::http1;
use hyper_util::rt::TokioIo;
use hyper_util::service::TowerToHyperService;
use tokio::sync::broadcast;
use tokio::task::JoinSet;
use tokio::time::{interval, timeout};

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

const DRAIN_TIMEOUT: Duration = Duration::from_secs(30);

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
                    _ = shutdown_rx.recv() => break,
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

    let app_for_accept = app.clone();
    let shutdown_for_accept = shutdown_tx.clone();

    // Accept loop
    let accept_task = tokio::spawn(async move {
        let mut conns: JoinSet<()> = JoinSet::new();

        loop {
            tokio::select! {
                _ = &mut accept_stop_rx => {
                    break;
                }
                res = listener.accept() => {
                    let (stream, _peer) = match res {
                        Ok(x) => x,
                        Err(_) => continue,
                    };

                    let app = app_for_accept.clone();
                    let mut shutdown_rx = shutdown_for_accept.subscribe();

                    conns.spawn(async move {
                        let io = TokioIo::new(stream);

                        // Adapt tower::Service (axum) into hyper::Service.
                        let svc = app.into_service();
                        let hyper_svc = TowerToHyperService::new(svc);

                        // HTTP/1 only here:
                        // - keep-alive disabled to reduce long-lived idle connections
                        // - draining is handled by:
                        //   - stop accepting new TCP connections
                        //   - `Connection: close` header during draining
                        //   - write rejection middleware
                        let conn_fut = http1::Builder::new()
                            .keep_alive(false)
                            .serve_connection(io, hyper_svc);

                        tokio::select! {
                            _ = shutdown_rx.recv() => {
                                // Best-effort: stop waiting; dropping future will close the socket.
                                // This is a hard stop for this connection.
                            }
                            _ = conn_fut => {}
                        }
                    });
                }
            }
        }

        // Drain existing connections spawned by this accept loop with a timeout.
        let _ = timeout(DRAIN_TIMEOUT, async {
            while let Some(_res) = conns.join_next().await {}
        })
        .await;
    });

    // Wait for shutdown signal
    shutdown_signal().await;
    println!("shutdown signal received");

    // Enter draining mode
    ready.store(false, Ordering::Relaxed);

    // Notify workers and connections
    let _ = shutdown_tx.send(());

    // Stop accepting new connections
    let _ = accept_stop_tx.send(());

    // Wait for accept loop and connection drain to finish (bounded)
    let _ = timeout(DRAIN_TIMEOUT, accept_task).await;

    println!("server exited cleanly");
    Ok(())
}
