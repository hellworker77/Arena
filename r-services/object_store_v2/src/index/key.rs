use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::{create_dir_all, File};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use crate::error::error::StoreResult;
use crate::index::io::{read_len_prefixed, write_len_prefixed};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum KeyRecord {
    Put { key: String, version: u32, hash: [u8; 32], ts: u64 },
    Delete { key: String, version: u32, ts: u64 },
}

pub struct KeyIndexStore {
    pub mem: HashMap<String, KeyRecord>,
    pub sstables: Vec<PathBuf>,
}

impl KeyIndexStore {
    pub fn new() -> Self {
        Self { mem: HashMap::new(), sstables: vec![] }
    }

    pub fn apply(&mut self, rec: KeyRecord) {
        match &rec {
            KeyRecord::Put { key, .. } | KeyRecord::Delete { key, .. } => {
                self.mem.insert(key.clone(), rec);
            }
        }
    }

    pub fn flush(&mut self, path: &Path) -> StoreResult<()> {
        if let Some(p) = path.parent() {
            create_dir_all(p)?;
        }
        let mut f = File::create(path)?;
        f.write_all(b"KEY1")?;
        f.write_all(&(self.mem.len() as u32).to_le_bytes())?;
        for rec in self.mem.values() {
            let buf = bincode::serialize(rec)?;
            write_len_prefixed(&mut f, &buf)?;
        }
        f.sync_all()?;
        self.sstables.push(path.to_path_buf());
        self.mem.clear();
        Ok(())
    }

    /// Minimal (but correct) latest lookup: mem first, then SST newest->the oldest linear scan.
    /// Production: add index/bloom/binary search later.
    pub fn get_latest(&self, key: &str) -> Option<KeyRecord> {
        if let Some(r) = self.mem.get(key) {
            return Some(r.clone());
        }
        for sst in self.sstables.iter().rev() {
            let mut f = File::open(sst).ok()?;
            let mut magic = [0u8; 4];
            f.read_exact(&mut magic).ok()?;
            if &magic != b"KEY1" {
                continue;
            }
            let mut cnt = [0u8; 4];
            f.read_exact(&mut cnt).ok()?;
            let n = u32::from_le_bytes(cnt);
            for _ in 0..n {
                let buf = read_len_prefixed(&mut f).ok()??;
                let rec = bincode::deserialize::<KeyRecord>(&buf).ok()?;
                match &rec {
                    KeyRecord::Put { key: k, .. } | KeyRecord::Delete { key: k, .. } if k == key => {
                        return Some(rec);
                    }
                    _ => {}
                }
            }
        }
        None
    }

    /// Returns the latest visible KeyRecord for every key.
    ///
    /// Semantics:
    /// - memtable overrides SSTables
    /// - SSTables are scanned from newest to oldest
    /// - the first occurrence of a key wins
    /// - Delete records are included (tombstones)
    ///
    /// Complexity:
    /// - O(mem + total_sst_records)
    /// - Acceptable for GC / maintenance paths
    pub fn iter_latest(&self) -> StoreResult<Vec<KeyRecord>> {
        let mut seen: HashMap<String, KeyRecord> = HashMap::new();

        // 1) Memtable has the highest priority.
        for (key, rec) in self.mem.iter() {
            seen.insert(key.clone(), rec.clone());
        }

        // 2) SSTables: newest -> oldest.
        for sst in self.sstables.iter().rev() {
            let mut f = File::open(sst)?;

            let mut magic = [0u8; 4];
            f.read_exact(&mut magic)?;
            if &magic != b"KEY1" {
                // Corrupt or foreign file; skip defensively.
                continue;
            }

            let mut cnt = [0u8; 4];
            f.read_exact(&mut cnt)?;
            let n = u32::from_le_bytes(cnt);

            for _ in 0..n {
                let buf = match read_len_prefixed(&mut f)? {
                    Some(b) => b,
                    None => break,
                };

                let rec: KeyRecord = match bincode::deserialize(&buf) {
                    Ok(r) => r,
                    Err(_) => continue,
                };

                let key = match &rec {
                    KeyRecord::Put { key, ..} => key,
                    KeyRecord::Delete { key, ..} => key,
                };

                // First occurrence wins (newest).
                if !seen.contains_key(key) {
                    seen.insert(key.clone(), rec);
                }
            }
        }

        Ok(seen.into_values().collect())
    }
}