#[derive(Clone, Debug)]
pub struct GcConfig {
    /// Minimum dead ratio across sealed segments to trigger GC/compaction.
    pub min_dead_ratio: f64,

    /// Minimum dead bytes across sealed segments to trigger GC/compaction.
    pub min_dead_bytes: u64,

    /// Per-segment dead ratio threshold to rewrite a segment (compaction).
    pub segment_rewrite_dead_ratio: f64,

    /// Per-segment dead ratio threshold to drop a segment entirely.
    pub segment_drop_dead_ratio: f64,

    /// Hard cap on how many segments may be rewritten in one run.
    pub max_rewrite_segments: usize,

    /// Hard cap on how many segments may be dropped in one run.
    pub max_drop_segments: usize,
}

impl Default for GcConfig {
    fn default() -> Self {
        Self {
            min_dead_ratio: 0.3,
            min_dead_bytes: 1 * 1024 * 1024 * 1024, // 1 GiB
            segment_rewrite_dead_ratio: 0.35,
            segment_drop_dead_ratio: 0.95,
            max_rewrite_segments: 4,
            max_drop_segments: 16,
        }
    }
}