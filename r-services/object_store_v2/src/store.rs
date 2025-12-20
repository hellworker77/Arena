use crate::error::error::{StoreError, StoreResult};
use crate::index::cas::{CasIndexStore, CasRecord};
use crate::index::key::{KeyIndexStore, KeyRecord};
use crate::manifest::{Manifest, ManifestRecord};
use crate::segment::{read_object, SegmentWriter};
use crate::wal::{Wal, WalRecord};
use sha2::{Digest, Sha256};
use std::collections::HashMap;
use std::path::PathBuf;
use crate::gc::config::GcConfig;
use crate::gc::{executor, planner, snapshot};

pub const SEG_OBJ_HDR_LEN: u64 = 32 + 12 + 8 + 8; // hash+nonce+size_plain+size_cipher

#[derive(Debug, Clone)]
pub struct CasEntry {
    pub segment_id: u64,
    pub offset: u64,
    pub size: u64,
    pub refcount: i64,
}

#[derive(Debug)]
pub struct CasIndex {
    pub map: HashMap<[u8; 32], CasEntry>,
}

fn sha256(data: &[u8]) -> [u8; 32] {
    let mut h = Sha256::new();
    h.update(data);
    h.finalize().into()
}

fn now_ts() -> u64 {
    use std::time::{SystemTime, UNIX_EPOCH};
    SystemTime::now().duration_since(UNIX_EPOCH).unwrap_or_default().as_secs()
}

pub struct ObjectStore {
    pub wal: Wal,
    pub manifest: Manifest,

    pub key_store: KeyIndexStore,
    pub cas_store: CasIndexStore,

    pub cas: CasIndex,

    pub segment: SegmentWriter,
    pub segments: HashMap<u64, PathBuf>,

    pub dir_index: PathBuf,

    pub max_segment_bytes: u64,
    pub max_segment_objects: u32,
}

impl ObjectStore {

    /// Best-effort GC + compaction.
    /// - Never touches active segment
    /// - Never panics
    /// - Safe to call periodically
    /// - If thresholds are not met, does nothing
    pub fn try_gc_compact(&mut self, cfg: &GcConfig) -> StoreResult<()> {
        // 1) Build snapshot (read-only, no side effects)
        let snap = snapshot::build_snapshot(
            self.cas.map.clone(),
            &self.key_store,
            self.segments.clone(),
            self.active_segment_id(),
            self.sealed_segment_ids(),
        )?;

        // 2) Build GC plan based on thresholds
        let plan = planner::build_plan(cfg, &snap);

        // 3) Nothing to do — exit fast
        if plan.actions.is_empty() {
            return Ok(());
        }

        // 4) Execute plan (manifest-driven, crash-safe)
        let segments_dir = self.segments_dir();

        executor::execute_plan(
            &plan,
            &snap,
            &mut self.manifest,
            &segments_dir,
        )?;

        Ok(())
    }


    pub fn put(&mut self, key: String, plain: &[u8]) -> StoreResult<()> {
        self.rotate_if_needed()?;
        let hash = sha256(plain);
        let ts = now_ts();

        // versioning: simplistic (mem-only); good enough for now
        let version = match self.key_store.get_latest(&key) {
            Some(KeyRecord::Put { version, .. }) | Some(KeyRecord::Delete { version, .. }) => version + 1,
            None => 1,
        };

        // WAL intent
        self.wal.append(&WalRecord::Put {
            key: key.clone(),
            version,
            hash,
            size: plain.len() as u64,
            ts,
        })?;

        // dedup
        if let Some(e) = self.cas.map.get_mut(&hash) {
            e.refcount += 1;
            self.cas_store.apply(CasRecord::RefInc { hash });
            self.key_store.apply(KeyRecord::Put { key, version, hash, ts });
            self.wal.append(&WalRecord::Commit)?;
            return Ok(());
        }

        // codec omitted: cipher == plain in this minimal version
        let nonce = rand::random::<[u8; 12]>();
        let cipher = plain;

        let offset = self.segment.write_object(hash, nonce, cipher, plain.len() as u64)?;
        self.segment.flush_data()?;

        // update CAS runtime + store
        self.cas.map.insert(
            hash,
            CasEntry {
                segment_id: self.segment.segment_id,
                offset,
                size: cipher.len() as u64,
                refcount: 1,
            },
        );
        self.cas_store.apply(CasRecord::Add {
            hash,
            segment_id: self.segment.segment_id,
            offset,
            size: cipher.len() as u64,
        });

        // update key store
        self.key_store.apply(KeyRecord::Put { key, version, hash, ts });

        // commit barrier
        self.wal.append(&WalRecord::Commit)?;
        Ok(())
    }

    pub fn get(&self, key: &str) -> StoreResult<Vec<u8>> {
        let rec = self.key_store.get_latest(key).ok_or(StoreError::NotFound)?;

        let hash = match rec {
            KeyRecord::Put { hash, .. } => hash, // ✅ no deref (it’s a value)
            KeyRecord::Delete { .. } => return Err(StoreError::Deleted),
        };

        let entry = self.cas.map.get(&hash).ok_or(StoreError::CasMiss)?;
        let seg_path = self.segments.get(&entry.segment_id).ok_or(StoreError::SegmentMissing)?;

        let (cipher, read_hash) = read_object(seg_path, entry.offset)?;
        if read_hash != hash {
            return Err(StoreError::HashMismatch);
        }
        Ok(cipher)
    }

    pub fn delete(&mut self, key: String) -> StoreResult<()> {
        let ts = now_ts();
        let version = match self.key_store.get_latest(&key) {
            Some(KeyRecord::Put { version, .. }) | Some(KeyRecord::Delete { version, .. }) => version + 1,
            None => 1,
        };

        // find old hash
        let old_hash = match self.key_store.get_latest(&key) {
            Some(KeyRecord::Put { hash, .. }) => Some(hash),
            _ => None,
        };

        self.wal.append(&WalRecord::Delete { key: key.clone(), version, ts })?;

        if let Some(h) = old_hash {
            if let Some(e) = self.cas.map.get_mut(&h) {
                e.refcount -= 1;
            }
            self.cas_store.apply(CasRecord::RefDec { hash: h });
        }

        self.key_store.apply(KeyRecord::Delete { key, version, ts });
        self.wal.append(&WalRecord::Commit)?;
        Ok(())
    }

    pub fn checkpoint(&mut self) -> StoreResult<()> {
        let key_sst = self.dir_index.join(format!("key-{}.sst", now_ts()));
        let cas_sst = self.dir_index.join(format!("cas-{}.sst", now_ts()));

        self.key_store.flush(&key_sst)?;
        self.manifest.append(&ManifestRecord::NewKeySst { path: key_sst })?;

        self.cas_store.flush(&cas_sst)?;
        self.manifest.append(&ManifestRecord::NewCasSst { path: cas_sst })?;

        self.manifest.append(&ManifestRecord::Checkpoint { wal_seq: 0 })?;
        Ok(())
    }

    pub fn locate_for_read(&self, key: &str) -> StoreResult<(std::path::PathBuf, u64, u64, String)> {
        let rec = self.key_store.get_latest(key).ok_or(StoreError::NotFound)?;

        let hash = match rec {
            KeyRecord::Put { hash, .. } => hash,
            KeyRecord::Delete { .. } => return Err(StoreError::Deleted),
        };

        let entry = self.cas.map.get(&hash).ok_or(StoreError::CasMiss)?;
        let seg_path = self.segments.get(&entry.segment_id).ok_or(StoreError::SegmentMissing)?.clone();

        // payload starts after object header
        let payload_offset = entry.offset + SEG_OBJ_HDR_LEN;
        let payload_len = entry.size;

        // ETag from hash (stable)
        let etag = format!("\"sha256:{}\"", hex::encode(hash));

        Ok((seg_path, payload_offset, payload_len, etag))
    }

    fn rotate_if_needed(&mut self) -> StoreResult<()> {
        let has_data = self.segment.current_objects() > 0;

        let need = has_data && (
            self.segment.current_size() >= self.max_segment_bytes
                || self.segment.current_objects() >= self.max_segment_objects
        );

        if !need {
            return Ok(());
        }

        // 1) seal current segment (fsync + finalize)
        let old_id = self.segment.segment_id;
        let old_path = self.segments.get(&old_id).cloned().ok_or(StoreError::SegmentMissing)?;

        let _sealed_path = self.segment.seal()?; // writes object_count + sync_all
        self.manifest.append(&ManifestRecord::SealSegment { segment_id: old_id })?;

        // 2) create new active segment
        let new_id = self.segments.keys().max().copied().unwrap_or(0) + 1;
        let new_path = old_path
            .parent()
            .unwrap()
            .join(format!("seg-{new_id:05}.seg"));

        self.manifest.append(&ManifestRecord::NewSegment { segment_id: new_id, path: new_path.clone() })?;
        self.manifest.append(&ManifestRecord::ActiveSegment { segment_id: new_id })?;

        self.segments.insert(new_id, new_path.clone());
        self.segment = SegmentWriter::create(&new_path, new_id)?;

        Ok(())
    }

    fn active_segment_id(&self) -> u64 {
        self.segment.segment_id
    }

    fn sealed_segment_ids(&self) -> std::collections::HashSet<u64> {
        self.segments
            .keys()
            .copied()
            .filter(|id| *id != self.active_segment_id())
            .collect()
    }

    fn segments_dir(&self) -> std::path::PathBuf {
        // All segment paths already live under one directory.
        // We reuse any existing segment path to infer the directory.
        self.segments
            .values()
            .next()
            .and_then(|p| p.parent())
            .map(|p| p.to_path_buf())
            .expect("segments directory must exist")
    }
}