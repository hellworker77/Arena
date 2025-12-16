use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::{create_dir_all, File};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use crate::error::error::StoreResult;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum KeyRecord {
    Put { key: String, version: u32, hash: [u8; 32], ts: u64 },
    Delete { key: String, version: u32, ts: u64 },
}

fn write_len_prefixed(file: &mut File, bytes: &[u8]) -> StoreResult<()> {
    file.write_all(&(bytes.len() as u32).to_le_bytes())?;
    file.write_all(bytes)?;
    Ok(())
}

fn read_len_prefixed(file: &mut File) -> StoreResult<Option<Vec<u8>>> {
    let mut len = [0u8; 4];
    if file.read_exact(&mut len).is_err() {
        return Ok(None);
    }
    let n = u32::from_le_bytes(len) as usize;
    let mut buf = vec![0u8; n];
    file.read_exact(&mut buf)?;
    Ok(Some(buf))
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

    /// Minimal (but correct) latest lookup: mem first, then SST newest->oldest linear scan.
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
}