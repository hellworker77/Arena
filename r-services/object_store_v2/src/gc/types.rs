use std::collections::{HashMap, HashSet};
use std::path::PathBuf;
use crate::segment::SegmentObjectInfo;

pub type Hash32 = [u8; 32];

#[derive(Debug, Clone)]
pub struct SegmentStats {
    pub segment_id: u64,
    pub total_bytes: u64,
    pub live_bytes: u64,
    pub dead_bytes: u64,
    pub dead_ratio: f64,
}

#[derive(Debug, Clone)]
pub struct GcSnapshot {
    /// Live hashes referenced by the latest KeyIndex view.
    pub live_hashes: HashSet<Hash32>,

    /// CAS map materialized (hash -> entry).
    pub cas_entries: HashMap<Hash32, crate::store::CasEntry>,

    /// Sealed segment scan (segment_id -> hash -> object info).
    pub sealed_scan: HashMap<u64, HashMap<Hash32, SegmentObjectInfo>>,

    /// All segments from manifest (active + sealed).
    pub segments: HashMap<u64, PathBuf>,

    pub active_id: u64,
    pub sealed_ids: HashSet<u64>,
}

#[derive(Debug, Clone)]
pub enum GcAction {
    /// Drop a fully-dead sealed segment.
    DropSegment { segment_id: u64 },

    /// Rewrite a sealed segment into a new sealed segment.
    RewriteSegment { segment_id: u64 },
}

#[derive(Debug, Clone)]
pub struct GcPlan {
    pub global_dead_ratio: f64,
    pub global_dead_bytes: u64,
    pub per_segment: Vec<SegmentStats>,
    pub actions: Vec<GcAction>,
}