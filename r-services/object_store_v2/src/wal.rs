use serde::{Deserialize, Serialize};
use std::fs::{create_dir_all, File, OpenOptions};
use std::io::{Read, Write};
use std::path::Path;
use crate::error::error::StoreResult;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum WalRecord {
    Put {
        key: String,
        version: u32,
        hash: [u8; 32],
        size: u64,
        ts: u64,
    },
    Delete {
        key: String,
        version: u32,
        ts: u64,
    },
    Commit,
}

pub fn write_len_prefixed(file: &mut File, bytes: &[u8]) -> StoreResult<()> {
    file.write_all(&(bytes.len() as u32).to_le_bytes())?;
    file.write_all(bytes)?;
    Ok(())
}

pub fn read_len_prefixed(file: &mut File) -> StoreResult<Option<Vec<u8>>> {
    let mut len = [0u8; 4];
    if file.read_exact(&mut len).is_err() {
        return Ok(None);
    }
    let n = u32::from_le_bytes(len) as usize;
    let mut buf = vec![0u8; n];
    file.read_exact(&mut buf)?;
    Ok(Some(buf))
}

pub struct Wal {
    file: File,
}

impl Wal {
    pub fn open(path: &Path) -> StoreResult<Self> {
        if let Some(p) = path.parent() {
            create_dir_all(p)?;
        }
        let file = OpenOptions::new().create(true).append(true).open(path)?;
        Ok(Self { file })
    }

    pub fn append(&mut self, rec: &WalRecord) -> StoreResult<()> {
        let buf = bincode::serialize(rec)?;
        write_len_prefixed(&mut self.file, &buf)?;
        self.file.sync_data()?;
        Ok(())
    }

    pub fn read_all(path: &Path) -> StoreResult<Vec<WalRecord>> {
        if !path.exists() {
            return Ok(vec![]);
        }

        let mut f = File::open(path)?;
        let mut out = Vec::new();
        while let Some(buf) = read_len_prefixed(&mut f)? {
            out.push(bincode::deserialize::<WalRecord>(&buf)?);
        }
        Ok(out)
    }
}
