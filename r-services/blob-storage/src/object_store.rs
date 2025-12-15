use crate::cas_index::{CasEntry, CasIndex};
use crate::key_index::KeyIndex;
use crate::segment_writer::SegmentWriter;
use crate::utils::sha256;
use crate::wal::{Wal, WalRecord};

pub struct ObjectStore {
    pub wal: Wal,
    pub cas: CasIndex,
    pub keys: KeyIndex,
    pub segment: SegmentWriter
}

impl ObjectStore {
    pub fn put(&mut self, key: String, plain: &[u8]) -> Result<(), Box<dyn std::error::Error>> {
        // 1. hash (CAS identity)
        let hash = sha256(plain);

        // 2. versioning
        let version = self.keys.put(key.clone(), hash);

        // 3. WAL intent
        self.wal.append(&WalRecord::Put {
            key: key.clone(),
            version,
            hash,
            size: plain.len() as u64,
        })?;

        //4. CAS check(dedup)
        if let Some(entry) = self.cas.map.get_mut(&hash) {
            entry.refcount += 1;
            self.wal.append(&WalRecord::Commit)?;
            return Ok(());
        }

        // 5. encode + write segment
        let nonce = rand::random::<[u8; 12]>();
        let cipher = plain;
        let offset = self
            .segment
            .write_object(hash, nonce, cipher, plain.len() as u64)?;

        // 6. update CAS
        self.cas.insert(
            hash,
            CasEntry {
                segment_id: self.segment.segment_id,
                offset,
                size: cipher.len() as u64,
                refcount: 1,
            },
        );
        
        // 7. Commit
        self.wal.append(&WalRecord::Commit)?;
        Ok(())
    }
}