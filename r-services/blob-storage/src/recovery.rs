use std::fs::File;
use std::io::Read;
use std::path::Path;
use crate::cas_index_store::{CasIndexStore, CasRecord};
use crate::key_index_store::{KeyIndexStore, KeyRecord};
use crate::manifest::ManifestRecord;
use crate::wal::WalRecord;

pub struct Recovery;

impl Recovery {
    pub fn recover(
        manifest_path: &Path,
        wal_path: &Path,
    ) -> Result<(KeyIndexStore, CasIndexStore), Box<dyn std::error::Error>> {
        let mut key_index = KeyIndexStore::new();
        let mut cas_index = CasIndexStore::new();


        // load manifest
        if manifest_path.exists() {
            let mut f = File::open(manifest_path)?;
            loop {
                let mut len = [0u8; 4];
                if f.read_exact(&mut len).is_err() {
                    break;
                }
                let len = u32::from_le_bytes(len);
                let mut buf = vec![0u8; len as usize];
                f.read_exact(&mut buf)?;
                let rec: ManifestRecord = bincode::deserialize(&buf)?;
                match rec {
                    ManifestRecord::NewKeySst { path } => key_index.sstables.push(path),
                    ManifestRecord::NewCasSst { path } => cas_index.sstables.push(path),
                    _ => {}
                }
            }
        }


        // replay WAL
        if wal_path.exists() {
            let mut f = File::open(wal_path)?;
            loop {
                let mut len = [0u8; 4];
                if f.read_exact(&mut len).is_err() {
                    break;
                }
                let len = u32::from_le_bytes(len);
                let mut buf = vec![0u8; len as usize];
                f.read_exact(&mut buf)?;
                let rec: WalRecord = bincode::deserialize(&buf)?;
                match rec {
                    WalRecord::Put { key, version, hash, .. } => {
                        key_index.apply(KeyRecord::Put { key, version, hash, ts: 0 });
                        cas_index.apply(CasRecord::RefInc { hash });
                    }
                    WalRecord::Delete { key, version } => {
                        key_index.apply(KeyRecord::Delete { key, version, ts: 0 });
                    }
                    WalRecord::Commit => {}
                }
            }
        }


        Ok((key_index, cas_index))
    }
}