use std::collections::HashMap;
use std::fs::create_dir_all;
use std::path::{Path, PathBuf};
use crate::error::error::{StoreError, StoreResult};
use crate::index::cas::{CasIndexStore, CasRecord};
use crate::manifest::{Manifest, ManifestRecord};
use crate::recovery::Recovery;
use crate::segment::{scan_all_segments, SegmentObjectInfo, SegmentWriter};
use crate::store::{CasEntry, CasIndex, ObjectStore};
use crate::wal::Wal;

#[derive(Clone, Copy)]
pub enum CasMaterializeMode {
    Strict,
    Permissive,
}

fn materialize_cas(
    cas_store: &CasIndexStore,
    segments: &HashMap<u64, PathBuf>, // all segments (active+sealed)
    sealed_scan: &HashMap<u64, HashMap<[u8; 32], SegmentObjectInfo>>, // only sealed scanned
    active_id: u64,
    mode: CasMaterializeMode,
) -> StoreResult<HashMap<[u8; 32], CasEntry>> {
    let mut map: HashMap<[u8; 32], CasEntry> = HashMap::new();

    let records = cas_store.iter_all()?;
    for rec in records {
        match rec {
            CasRecord::Add {
                hash,
                segment_id,
                offset,
                size,
            } => {
                // segment must exist
                if !segments.contains_key(&segment_id) {
                    match mode {
                        CasMaterializeMode::Strict => {
                            return Err(StoreError::CasDanglingSegment(segment_id))
                        }
                        CasMaterializeMode::Permissive => continue,
                    }
                }

                let is_sealed = sealed_scan.contains_key(&segment_id);
                if is_sealed {
                    // sealed must be validated by scan
                    let scan = sealed_scan.get(&segment_id).unwrap();
                    if !scan.contains_key(&hash) {
                        match mode {
                            CasMaterializeMode::Strict => return Err(StoreError::CasDanglingObject),
                            CasMaterializeMode::Permissive => continue,
                        }
                    }
                } else {
                    // not sealed => must be active
                    if segment_id != active_id {
                        match mode {
                            CasMaterializeMode::Strict => {
                                return Err(StoreError::CasDanglingSegment(segment_id))
                            }
                            CasMaterializeMode::Permissive => continue,
                        }
                    }
                    // active segment isn't scanned/validated here
                }

                map.entry(hash).or_insert(CasEntry {
                    segment_id,
                    offset,
                    size,
                    refcount: 0,
                });
            }

            CasRecord::RefInc { hash } => {
                if let Some(e) = map.get_mut(&hash) {
                    e.refcount += 1;
                } else if matches!(mode, CasMaterializeMode::Strict) {
                    return Err(StoreError::CasDanglingObject);
                }
            }

            CasRecord::RefDec { hash } => {
                if let Some(e) = map.get_mut(&hash) {
                    e.refcount -= 1;
                } else if matches!(mode, CasMaterializeMode::Strict) {
                    return Err(StoreError::CasDanglingObject);
                }
            }
        }
    }

    // Heal offsets ONLY for sealed segments (active isn't scanned)
    for (hash, e) in map.iter_mut() {
        if let Some(scan) = sealed_scan.get(&e.segment_id) {
            if let Some(info) = scan.get(hash) {
                e.offset = info.offset;
                e.size = info.size_cipher;
            }
        }
    }

    Ok(map)
}

pub fn bootstrap(base: impl AsRef<Path>) -> StoreResult<ObjectStore> {
    let base = base.as_ref();
    let dir_wal = base.join("wal");
    let dir_segments = base.join("segments");
    let dir_index = base.join("index");
    let dir_meta = base.join("meta");

    create_dir_all(&dir_wal)?;
    create_dir_all(&dir_segments)?;
    create_dir_all(&dir_index)?;
    create_dir_all(&dir_meta)?;

    let wal_path = dir_wal.join("00000001.wal");
    let manifest_path = dir_meta.join("manifest.log");

    // Recovery (manifest+wal, commit-aware)
    let (mst, key_store, cas_store, _wal_records) = Recovery::recover(&manifest_path, &wal_path)?;

    // 0) ВАЛИДАЦИЯ: SST из manifest должны существовать (иначе NotFound в iter_all / get_latest)
    // Если хочешь permissive — можно просто пропускать отсутствующие файлы.
    for p in mst.key_sst.iter().chain(mst.cas_sst.iter()) {
        if !p.exists() {
            // выбери одно:
            return Err(StoreError::SegmentScan(format!("SST file missing: {:?}", p)));
        }
    }

    // All segments from manifest (active + sealed)
    let mut segments = mst.segments.clone();
    let mut manifest = Manifest::open(&manifest_path)?;

    // Ensure active segment exists (create if missing)
    let active_id = if let Some(id) = mst.active {
        id
    } else {
        let new_id = if segments.is_empty() { 0 } else { segments.keys().max().copied().unwrap() + 1 };
        let new_path = dir_segments.join(format!("seg-{new_id:05}.seg"));

        manifest.append(&ManifestRecord::NewSegment { segment_id: new_id, path: new_path.clone() })?;
        manifest.append(&ManifestRecord::ActiveSegment { segment_id: new_id })?;

        segments.insert(new_id, new_path);
        new_id
    };

    let active_path = segments
        .get(&active_id)
        .ok_or(StoreError::ManifestMissingSegment(active_id))?
        .clone();

    // 1) sealed segments ONLY are scanned/validated — но только если файл реально существует
    let mut sealed_segments: HashMap<u64, PathBuf> = HashMap::new();
    for id in &mst.sealed {
        if let Some(p) = segments.get(id) {
            if p.exists() {
                sealed_segments.insert(*id, p.clone());
            } else {
                // sealed segment in manifest but file missing => это уже нарушение инварианта
                return Err(StoreError::ManifestMissingSegment(*id));
            }
        } else {
            return Err(StoreError::ManifestMissingSegment(*id));
        }
    }

    let seg_scan = scan_all_segments(&sealed_segments)?;

    // CAS materialization: validate sealed, tolerate active
    // ВАЖНО: твоя materialize_cas должна принимать active_id (у тебя это так)
    let cas_map = materialize_cas(
        &cas_store,
        &segments,
        &seg_scan,
        active_id,
        CasMaterializeMode::Strict,
    )?;
    let cas = CasIndex { map: cas_map };

    // Open WAL (creates file if missing)
    let wal = Wal::open(&wal_path)?;

    // Open active segment WITHOUT truncation.
    // open_append ДОЛЖЕН иметь create(true) + писать header если файл новый.
    let segment = SegmentWriter::open_append(&active_path, active_id)?;

    Ok(ObjectStore {
        wal,
        manifest,
        key_store,
        cas_store,
        cas,
        segment,
        segments,
        dir_index,
        max_segment_bytes: 64 * 1024 * 1024,
        max_segment_objects: 50_000,
    })
}