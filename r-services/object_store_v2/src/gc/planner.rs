use crate::gc::config::GcConfig;
use crate::gc::types::{GcAction, GcPlan, GcSnapshot, SegmentStats};

pub fn build_plan(cfg: &GcConfig, snap: &GcSnapshot) -> GcPlan {
    let mut per_segment: Vec<SegmentStats> = Vec::new();

    // Compute per-segment total bytes from sealed scan (source of truth).
    // We count cipher payload only; header bytes can be added if you want exact disk usage.
    for (segment_id, scan) in &snap.sealed_scan {
        let mut total_bytes = 0u64;
        let mut live_bytes = 0u64;

        for (hash, info) in scan.iter() {
            total_bytes = total_bytes.saturating_add(info.size_cipher);
            if snap.live_hashes.contains(hash) {
                live_bytes = live_bytes.saturating_add(info.size_cipher);
            }
        }

        let dead_bytes = total_bytes.saturating_sub(live_bytes);
        let dead_ratio = if total_bytes == 0 {
            0.0
        } else {
            dead_bytes as f64 / total_bytes as f64
        };

        per_segment.push(SegmentStats {
            segment_id: *segment_id,
            total_bytes,
            live_bytes,
            dead_bytes,
            dead_ratio,
        });
    }

    // Global threshold trigger.
    let global_total: u64 = per_segment.iter().map(|s| s.total_bytes).sum();
    let global_dead: u64 = per_segment.iter().map(|s| s.dead_bytes).sum();
    let global_dead_ratio = if global_total == 0 {
        0.0
    } else {
        global_dead as f64 / global_total as f64
    };

    let trigger = global_dead_ratio >= cfg.min_dead_ratio || global_dead >= cfg.min_dead_bytes;

    // If not triggered, return stats only (no actions).
    if !trigger {
        return GcPlan {
            global_dead_ratio,
            global_dead_bytes: global_dead,
            per_segment,
            actions: Vec::new(),
        }
    }

    // Decide actions by segment thresholds.
    // Prefer dropping heavily-dead segments, then rewriting moderately-dead segments.
    per_segment.sort_by(|a, b| b.dead_ratio.partial_cmp(&a.dead_ratio).unwrap());

    let mut actions: Vec<GcAction> = Vec::new();

    let mut drops = 0usize;
    let mut rewrites = 0usize;

    for s in &per_segment {
        if drops < cfg.max_drop_segments && s.dead_ratio >= cfg.segment_drop_dead_ratio {
            actions.push(GcAction::DropSegment { segment_id: s.segment_id });
            drops += 1;
            continue;
        }

        if rewrites < cfg.max_rewrite_segments && s.dead_ratio >= cfg.segment_rewrite_dead_ratio {
            actions.push(GcAction::RewriteSegment { segment_id: s.segment_id });
            rewrites += 1;
        }
    }

    GcPlan {
        global_dead_ratio,
        global_dead_bytes: global_dead,
        per_segment,
        actions,
    }
}
