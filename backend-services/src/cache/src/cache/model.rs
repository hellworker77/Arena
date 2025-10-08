use std::{
    collections::HashMap,
    sync::{Arc, Mutex},
    time::{Duration, Instant},
};

#[derive(Clone, Default)]
pub struct CacheStore {
    pub store: Arc<Mutex<HashMap<String, (String, Option<Instant>)>>>
}

impl CacheStore {
    pub fn insert(&self, key: String, value: String, ttl_ms: u64){
        let exp = if ttl_ms > 0 {
            Some(Instant::now() + Duration::from_millis(ttl_ms))
        } else {
            None
        };

        let mut store = self.store.lock().unwrap();
        store.insert(key, (value, exp));
    }

    pub fn get(&self, key: &str) -> Option<String> {
        let mut store = self.store.lock().unwrap();
        if let Some((val, exp)) = store.get(key) {
            if exp.map_or(true, |t| t > Instant::now()){
                return Some(val.clone())
            } else {
                store.remove(key);
            }
        }

        None
    }

    pub fn delete(&self, key: &str) -> bool {
        let mut store = self.store.lock().unwrap();
        
        store.remove(key).is_some()
    }

    pub fn clear(&self) {
        let mut store = self.store.lock().unwrap();
        store.clear();
    }

    pub fn retain_valid(&self) {
        let mut store = self.store.lock().unwrap();
        let now = Instant::now();

        store.retain(|_, (_, exp)| exp.map_or(true, |t| t > now));
    }
}