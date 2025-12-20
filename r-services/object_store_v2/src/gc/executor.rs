use crate::compaction::compact_one_segment;
use crate::error::error::{StoreError, StoreResult};
use crate::gc::types::{GcAction, GcPlan, GcSnapshot};
use crate::manifest::{Manifest, ManifestRecord};

pub fn execute_plan(
    plan: &GcPlan,
    snap: &GcSnapshot,
    manifest: &mut Manifest,
    dir_segments: &std::path::Path,
) -> StoreResult<()> {
    // Do nothing if there are no actions.
    if plan.actions.is_empty() {
        return Ok(());
    }

    // Never touch active segment.
    for a in &plan.actions {
        match a {
            GcAction::DropSegment { segment_id } | GcAction::RewriteSegment { segment_id } => {
                if *segment_id == snap.active_id {
                    return Err(StoreError::GcInvariantViolation);
                }
            }
        }
    }

    // Phase 1: apply drop actions (only fully-dead segments should be dropped).
    for a in &plan.actions {
        if let GcAction::DropSegment { segment_id } = a {
            drop_sealed_segment(*segment_id, snap, manifest)?;
        }
    }

    // Phase 2: rewrite actions (compaction).
    // Actual compaction writes new sealed segments and updates manifest.
    for a in &plan.actions {
        if let GcAction::RewriteSegment { segment_id } = a {
            compact_one_segment(
                *segment_id,
                snap,
                manifest,
                dir_segments,
            )?;
        }
    }

    Ok(())
}

fn drop_sealed_segment(
    segment_id: u64,
    snap: &GcSnapshot,
    manifest: &mut Manifest,
) -> StoreResult<()> {
    // Only sealed segments may be dropped.
    if !snap.sealed_ids.contains(&segment_id) {
        return Err(StoreError::GcInvariantViolation);
    }

    let path = snap
        .segments
        .get(&segment_id)
        .ok_or(StoreError::SegmentMissing)?
        .clone();

    // Persist decision before deleting files.
    manifest.append(&ManifestRecord::DropSegment { segment_id })?;

    // Best-effort file removal; recovery tolerates already-removed files if manifest says dropped.
    let _ = std::fs::remove_file(path);

    Ok(())
}