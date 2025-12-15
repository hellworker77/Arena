use std::collections::HashMap;
use std::fs::File;
use std::io::Write;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub enum KeyRecord {
    Put {
        key: String,
        version: u32,
        hash: [u8; 32],
        ts: u64,
    },
    Delete {
        key: String,
        version: u32,
        ts: u64,
    },
}

pub struct KeyIndexStore {
    mem: HashMap<String, KeyRecord>,
    pub(crate) sstables: Vec<PathBuf>,
}

impl KeyIndexStore {
    pub fn new() -> Self {
        Self {
            mem: HashMap::new(),
            sstables: vec![],
        }
    }

    pub fn apply(&mut self, rec: KeyRecord) {
        match &rec {
            KeyRecord::Put { key, .. } | KeyRecord::Delete { key, .. } => {
                self.mem.insert(key.clone(), rec);
            }
        }
    }

    pub fn flush(&mut self, path: &Path) -> Result<(), bincode::Error> {
        let mut f = File::create(path)?;
        f.write_all(b"KEY1")?;
        f.write_all(&(self.mem.len() as u32).to_le_bytes())?;
        for rec in self.mem.values() {
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