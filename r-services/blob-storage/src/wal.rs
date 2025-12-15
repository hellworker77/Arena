use std::error::Error;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::{Seek, SeekFrom, Write};
use std::path::Path;

#[derive(Debug, Serialize, Deserialize)]
pub enum WalRecord {
    Put {
        key: String,
        version: u32,
        hash: [u8; 32],
        size: u64,
    },
    Delete {
        key: String,
        version: u32,
    },
    Commit,
}

pub struct Wal {
    file: File,
}

impl Wal {
    pub fn open(path: &Path) -> Result<Self, std::io::Error> {
        let file = File::open(path)?;
        Ok(Wal { file })
    }

    pub fn append(&mut self, rec: &WalRecord) -> Result<(), bincode::Error> {
        let buf = bincode::serialize(rec)?;
        self.file.write_all(&(buf.len() as u32).to_le_bytes())?;
        self.file.write_all(&buf)?;
        self.file.sync_data()?;
        Ok(())
    }
}
