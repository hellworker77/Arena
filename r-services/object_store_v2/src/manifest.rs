use crate::error::error::StoreResult;
use serde::{Deserialize, Serialize};
use std::fs::{create_dir_all, File, OpenOptions};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use crate::wal::{read_len_prefixed, write_len_prefixed};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ManifestRecord {
    NewSegment { segment_id: u64, path: PathBuf },
    SealSegment { segment_id: u64 },
    ActiveSegment { segment_id: u64 },
    NewKeySst { path: PathBuf },
    NewCasSst { path: PathBuf },
    DropSegment { segment_id: u64 },
    Checkpoint { wal_seq: u64 },
}

pub struct Manifest {
    file: File,
}

impl Manifest {
    pub fn open(path: &Path) -> StoreResult<Self> {
        if let Some(p) = path.parent() {
            create_dir_all(p)?;
        }
        let file = OpenOptions::new().create(true).append(true).open(path)?;
        Ok(Self { file })
    }

    pub fn append(&mut self, rec: &ManifestRecord) -> StoreResult<()> {
        let buf = bincode::serialize(rec)?;
        write_len_prefixed(&mut self.file, &buf)?;
        self.file.sync_all()?;
        Ok(())
    }

    pub fn read_all(path: &Path) -> StoreResult<Vec<ManifestRecord>> {
        if !path.exists() {
            return Ok(vec![]);
        }

        let mut f = File::open(path)?;
        let mut out = Vec::new();
        while let Some(buf) = read_len_prefixed(&mut f)? {
            out.push(bincode::deserialize::<ManifestRecord>(&buf)?);
        }
        Ok(out)
    }
}

/// Aggregated manifest state for bootstrap/recovery
#[derive(Debug, Default)]
pub struct ManifestState {
    pub segments: std::collections::HashMap<u64, PathBuf>,
    pub sealed: std::collections::HashSet<u64>,
    pub active: Option<u64>,
    pub key_sst: Vec<PathBuf>,
    pub cas_sst: Vec<PathBuf>,
}

impl ManifestState {
    pub fn from_records(records: &[ManifestRecord]) -> Self {
        let mut st = ManifestState::default();
        for r in records {
            match r {
                ManifestRecord::NewSegment { segment_id, path } => {
                    st.segments.insert(*segment_id, path.clone());
                }
                ManifestRecord::SealSegment { segment_id } => {
                    st.sealed.insert(*segment_id);
                    if st.active == Some(*segment_id) {
                        st.active = None; // active cannot remain sealed
                    }
                }
                ManifestRecord::ActiveSegment { segment_id } => {
                    st.active = Some(*segment_id);
                }
                ManifestRecord::NewKeySst { path } => st.key_sst.push(path.clone()),
                ManifestRecord::NewCasSst { path } => st.cas_sst.push(path.clone()),
                ManifestRecord::DropSegment { segment_id } => {
                    st.segments.remove(segment_id);
                    st.sealed.remove(segment_id);
                    if st.active == Some(*segment_id) {
                        st.active = None;
                    }
                }
                ManifestRecord::Checkpoint { .. } => {}
            }
        }
        st
    }
}