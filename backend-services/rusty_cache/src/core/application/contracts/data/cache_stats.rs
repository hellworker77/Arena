#[derive(Debug, Clone)]
pub struct CacheStats {
    pub items: usize,

    pub total_size_bytes: usize,

    pub hits: u64,

    pub misses: u64,
}