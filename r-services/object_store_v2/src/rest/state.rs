use crate::store::ObjectStore;
use std::sync::{Arc, atomic::AtomicBool};
use tokio::sync::Mutex;

use super::Metrics;

/// Shared application state.
///
/// This is the ONLY axum state object.
#[derive(Clone)]
pub struct AppState {
    /// Object storage core.
    pub store: Arc<Mutex<ObjectStore>>,

    /// Prometheus-style counters.
    pub metrics: Arc<Metrics>,

    /// Readiness flag (false during shutdown).
    pub ready: Arc<AtomicBool>,
}
