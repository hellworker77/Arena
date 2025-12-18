use tokio::signal;

/// Waits for a shutdown signal (SIGTERM or Ctrl+C).
///
/// Contract:
/// - Resolves exactly once.
/// - Does NOT perform cleanup itself.
/// - Safe to await from multiple places (cloned future not required).
///
/// Unix:
/// - SIGTERM (Kubernetes / Docker)
/// - Ctrl+C (local dev)
///
/// Non-Unix:
/// - Ctrl+C only
pub async fn shutdown_signal() {
    #[cfg(unix)]
    {
        use tokio::signal::unix::{signal, SignalKind};

        let mut sigterm =
            signal(SignalKind::terminate()).expect("failed to install SIGTERM handler");

        tokio::select! {
            _ = signal::ctrl_c() => {},
            _ = sigterm.recv() => {},
        }
    }

    #[cfg(not(unix))]
    {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    }

    // Intentionally no logging here.
    // Logging is responsibility of the caller.
}