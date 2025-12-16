use crate::index::cas::{CasIndexStore, CasRecord};
use crate::index::key::{KeyIndexStore, KeyRecord};
use crate::manifest::{ManifestRecord, ManifestState};
use crate::wal::WalRecord;
use std::path::Path;
use crate::error::error::StoreResult;

/// Commit-aware recovery: applies only WAL records that end with Commit.
/// Also applies Delete with RefDec based on current key state during replay.
pub struct Recovery;

impl Recovery {
    pub fn recover(
        manifest_path: &Path,
        wal_path: &Path,
    ) -> StoreResult<(ManifestState, KeyIndexStore, CasIndexStore, Vec<WalRecord>)> {
        let manifest_records = crate::manifest::Manifest::read_all(manifest_path)?;
        let st = ManifestState::from_records(&manifest_records);

        let mut key_store = KeyIndexStore::new();
        key_store.sstables = st.key_sst.clone();

        let mut cas_store = CasIndexStore::new();
        cas_store.sstables = st.cas_sst.clone();

        let wal_records = crate::wal::Wal::read_all(wal_path)?;

        // replay WAL with commit barrier
        let mut pending: Vec<WalRecord> = Vec::new();
        for rec in &wal_records {
            match rec {
                WalRecord::Commit => {
                    for p in pending.drain(..) {
                        match p {
                            WalRecord::Put { key, version, hash, ts, .. } => {
                                key_store.apply(KeyRecord::Put { key, version, hash, ts });
                                cas_store.apply(CasRecord::RefInc { hash });
                            }
                            WalRecord::Delete { key, version, ts } => {
                                if let Some(KeyRecord::Put { hash, .. }) = key_store.mem.get(&key) {
                                    cas_store.apply(CasRecord::RefDec { hash: *hash });
                                }
                                key_store.apply(KeyRecord::Delete { key, version, ts });
                            }
                            WalRecord::Commit => {}
                        }
                    }
                }
                other => pending.push(other.clone()),
            }
        }

        Ok((st, key_store, cas_store, wal_records))
    }
}