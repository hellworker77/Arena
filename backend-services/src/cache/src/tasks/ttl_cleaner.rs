use tokio::time::{interval, Duration};
use crate::cache::model::CacheStore;

pub async fn start_ttl_cleaner(store: CacheStore) {
    let mut ticker = interval(Duration::from_secs(5));

    loop {
        ticker.tick().await;
        store.retain_valid();
    }
}