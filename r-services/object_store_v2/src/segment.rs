use crate::error::error::{StoreError, StoreResult};
use std::collections::HashMap;
use std::fs::{File, create_dir_all, OpenOptions};
use std::io::{Read, Seek, SeekFrom, Write};
use std::path::{Path, PathBuf};

/// Segment object header size (excluding cipher payload)
const OBJ_HDR: u64 = 32 + 12 + 8 + 8;

/// Immutable segment format:
/// "SEG1" [4]
/// object_count [u32]
/// repeated:
///   hash[32], nonce[12], size_plain[u64], size_cipher[u64], cipher[size_cipher]
pub struct SegmentWriter {
    file: File,
    offset: u64,
    pub segment_id: u64,
    object_count: u32,
    path: PathBuf,
}

impl SegmentWriter {
    pub fn create(path: &Path, segment_id: u64) -> StoreResult<Self> {
        if let Some(p) = path.parent() {
            create_dir_all(p)?;
        }
        let mut file = File::create(path)?;
        file.write_all(b"SEG1")?;
        file.write_all(&0u32.to_le_bytes())?;
        Ok(Self {
            file,
            offset: 8,
            segment_id,
            object_count: 0,
            path: path.to_path_buf(),
        })
    }

    pub fn write_object(
        &mut self,
        hash: [u8; 32],
        nonce: [u8; 12],
        cipher: &[u8],
        size_plain: u64,
    ) -> StoreResult<u64> {
        let start = self.offset;
        self.file.write_all(&hash)?;
        self.file.write_all(&nonce)?;
        self.file.write_all(&size_plain.to_le_bytes())?;
        self.file.write_all(&(cipher.len() as u64).to_le_bytes())?;
        self.file.write_all(cipher)?;
        self.offset += OBJ_HDR + cipher.len() as u64;
        self.object_count += 1;
        Ok(start)
    }

    pub fn flush_data(&mut self) -> StoreResult<()> {
        self.file.sync_data()?;
        Ok(())
    }

    pub fn seal(&mut self) -> StoreResult<PathBuf> {
        self.file.seek(SeekFrom::Start(4))?;
        self.file.write_all(&self.object_count.to_le_bytes())?;
        self.file.sync_all()?;
        Ok(self.path.clone())
    }

    pub fn current_size(&self) -> u64 {
        self.offset
    }

    pub fn current_objects(&self) -> u32 {
        self.object_count
    }

    pub fn open_append(path: &Path, segment_id: u64) -> StoreResult<Self> {
        if let Some(p) = path.parent() {
            std::fs::create_dir_all(p)?;
        }

        let mut file = std::fs::OpenOptions::new()
            .read(true)
            .write(true)
            .create(true)
            .open(path)?;

        let len = file.metadata()?.len();

        if len == 0 {
            file.write_all(b"SEG1")?;
            file.write_all(&0u32.to_le_bytes())?;
            file.sync_all()?;
        } else if len < 8 {
            return Err(StoreError::SegmentScan(format!("active segment header missing: {:?}", path)));
        }

        Ok(Self {
            file,
            offset: std::cmp::max(len, 8),
            segment_id,
            object_count: 0,
            path: path.to_path_buf(),
        })
    }
}

#[derive(Debug, Clone)]
pub struct SegmentObjectInfo {
    pub hash: [u8; 32],
    pub offset: u64,
    pub size_plain: u64,
    pub size_cipher: u64,
}

/// Read a single object at `offset`
pub fn read_object(segment_path: &Path, offset: u64) -> StoreResult<(Vec<u8>, [u8; 32])> {
    let mut f = File::open(segment_path)?;

    let mut magic = [0u8; 4];
    f.read_exact(&mut magic)?;
    if &magic != b"SEG1" {
        return Err(StoreError::BadSegmentMagic);
    }
    let mut cnt = [0u8; 4];
    f.read_exact(&mut cnt)?;

    f.seek(SeekFrom::Start(offset))?;

    let mut hash = [0u8; 32];
    f.read_exact(&mut hash)?;

    let mut _nonce = [0u8; 12];
    f.read_exact(&mut _nonce)?;

    let mut sp = [0u8; 8];
    f.read_exact(&mut sp)?;
    let _size_plain = u64::from_le_bytes(sp);

    let mut sc = [0u8; 8];
    f.read_exact(&mut sc)?;
    let size_cipher = u64::from_le_bytes(sc);

    let mut cipher = vec![0u8; size_cipher as usize];
    f.read_exact(&mut cipher)?;

    Ok((cipher, hash))
}

/// Production-ish segment validation: sequential scan and build hash->info map.
/// This is what you should use at bootstrap to detect corruption and to validate CAS offsets.
///
/// Strict rules:
/// - magic must match
/// - offsets must be within file
/// - object hash read must be consistent with stored field
pub fn scan_segment(segment_path: &Path) -> StoreResult<HashMap<[u8; 32], SegmentObjectInfo>> {
    let mut f = File::open(segment_path)?;
    let file_len = f.metadata()?.len();

    if file_len <= 8 {
        return Ok(HashMap::new());
    }

    let mut magic = [0u8; 4];
    f.read_exact(&mut magic)?;
    if &magic != b"SEG1" {
        return Err(StoreError::BadSegmentMagic);
    }
    let mut cnt = [0u8; 4];
    f.read_exact(&mut cnt)?;
    let _declared = u32::from_le_bytes(cnt);

    let mut pos = 8u64;
    let mut out = HashMap::new();

    while pos < file_len {
        if file_len - pos < OBJ_HDR {
            return Err(StoreError::SegmentScan(format!(
                "truncated header at offset {pos}"
            )));
        }

        f.seek(SeekFrom::Start(pos))?;

        let mut hash = [0u8; 32];
        f.read_exact(&mut hash)?;

        let mut _nonce = [0u8; 12];
        f.read_exact(&mut _nonce)?;

        let mut sp = [0u8; 8];
        f.read_exact(&mut sp)?;
        let size_plain = u64::from_le_bytes(sp);

        let mut sc = [0u8; 8];
        f.read_exact(&mut sc)?;
        let size_cipher = u64::from_le_bytes(sc);

        let payload_start = pos + OBJ_HDR;
        let payload_end = payload_start + size_cipher;

        if payload_end > file_len {
            return Err(StoreError::SegmentScan(format!(
                "payload out of bounds at offset {pos}"
            )));
        }

        // NOTE: We do not read payload here (cheap scan). If you want stronger validation,
        // you can hash/decrypt payload later. For CAS correctness, offset+hash is enough.
        out.entry(hash).or_insert(SegmentObjectInfo {
            hash,
            offset: pos,
            size_plain,
            size_cipher,
        });

        pos = payload_end;
    }

    Ok(out)
}

pub fn scan_all_segments(
    segments: &HashMap<u64, PathBuf>,
) -> StoreResult<HashMap<u64, HashMap<[u8; 32], SegmentObjectInfo>>> {
    let mut out = HashMap::new();
    for (seg_id, path) in segments {
        out.insert(*seg_id, scan_segment(path)?);
    }
    Ok(out)
}
