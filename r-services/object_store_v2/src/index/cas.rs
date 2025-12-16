use serde::{Deserialize, Serialize};
use std::fs::{create_dir_all, File};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use crate::error::error::{StoreError, StoreResult};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CasRecord {
    Add { hash: [u8; 32], segment_id: u64, offset: u64, size: u64 },
    RefInc { hash: [u8; 32] },
    RefDec { hash: [u8; 32] },
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

pub struct CasIndexStore {
    pub mem: Vec<CasRecord>,
    pub sstables: Vec<PathBuf>,
}

impl CasIndexStore {
    pub fn new() -> Self {
        Self { mem: vec![], sstables: vec![] }
    }

    pub fn apply(&mut self, rec: CasRecord) {
        self.mem.push(rec);
    }

    pub fn flush(&mut self, path: &Path) -> StoreResult<()> {
        if let Some(p) = path.parent() {
            create_dir_all(p)?;
        }
        let mut f = File::create(path)?;
        f.write_all(b"CAS1")?;
        f.write_all(&(self.mem.len() as u32).to_le_bytes())?;
        for rec in &self.mem {
            let buf = bincode::serialize(rec)?;
            write_len_prefixed(&mut f, &buf)?;
        }
        f.sync_all()?;
        self.sstables.push(path.to_path_buf());
        self.mem.clear();
        Ok(())
    }

    pub fn iter_all(&self) -> StoreResult<Vec<CasRecord>> {
        let mut out = Vec::new();

        for sst in &self.sstables {
            let mut f = File::open(sst)?;
            let mut magic = [0u8; 4];
            f.read_exact(&mut magic)?;
            if &magic != b"CAS1" {
                return Err(StoreError::BadSstMagic);
            }
            let mut cnt = [0u8; 4];
            f.read_exact(&mut cnt)?;
            let n = u32::from_le_bytes(cnt);
            for _ in 0..n {
                let buf = read_len_prefixed(&mut f)?.unwrap();
                out.push(bincode::deserialize::<CasRecord>(&buf)?);
            }
        }

        // include mem (if any)
        out.extend(self.mem.iter().cloned());

        Ok(out)
    }
}