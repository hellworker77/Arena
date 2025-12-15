use std::fs::File;
use std::io::Write;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub enum CasRecord {
    Add {
        hash: [u8; 32],
        segment_id: u64,
        offset: u64,
        size: u64,
    },
    RefInc {
        hash: [u8; 32],
    },
    RefDec {
        hash: [u8; 32],
    },
}

pub struct CasIndexStore {
    mem: Vec<CasRecord>,
    pub(crate) sstables: Vec<PathBuf>,
}

impl CasIndexStore {
    pub fn new() -> Self {
        Self { mem: vec![], sstables: vec![] }
    }


    pub fn apply(&mut self, rec: CasRecord) {
        self.mem.push(rec);
    }


    pub fn flush(&mut self, path: &Path) -> Result<(), bincode::Error> {
        let mut f = File::create(path)?;
        f.write_all(b"CAS1")?;
        f.write_all(&(self.mem.len() as u32).to_le_bytes())?;
        for rec in &self.mem {
            let buf = bincode::serialize(rec)?;
            f.write_all(&(buf.len() as u32).to_le_bytes())?;
            f.write_all(&buf)?;
        }
        f.sync_all()?;
        self.sstables.push(path.to_path_buf());
        self.mem.clear();
        Ok(())
    }
}