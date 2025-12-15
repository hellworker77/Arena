use std::fs::{File, OpenOptions};
use std::io::Write;
use std::path::{Path, PathBuf};

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub enum ManifestRecord {
    NewSegment { segment_id: u64, path: PathBuf },
    SealSegment { segment_id: u64 },
    NewKeySst { path: PathBuf },
    NewCasSst { path: PathBuf },
    DropSegment { segment_id: u64 },
    Checkpoint { wal_seq: u64 },
}

pub struct Manifest {
    file: File
}

impl Manifest {
    pub fn open(path: &Path) -> Result<Self, std::io::Error> {
        let file = OpenOptions::new().create(true).append(true).open(path)?;
        Ok(Self { file })
    }
    
    pub fn append(&mut self, rec: &ManifestRecord) -> Result<(), bincode::Error> {
        let buf = bincode::serialize(rec)?;
        self.file.write_all(&(buf.len() as u32).to_le_bytes())?;
        self.file.write_all(&buf)?;
        self.file.sync_all()?;
        Ok(())
    }
}
