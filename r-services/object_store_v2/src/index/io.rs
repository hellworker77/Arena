use std::fs::File;
use std::io::{Read, Write};
use crate::error::error::StoreResult;

/// Writes a length-prefixed byte slice.
/// Format: [u32 len][bytes...]
///
/// This is a shared on-disk primitive used by all index SSTables.
pub fn write_len_prefixed(file: &mut File, bytes: &[u8]) -> StoreResult<()> {
    let len = bytes.len() as u32;
    file.write_all(&len.to_le_bytes())?;
    file.write_all(bytes)?;
    Ok(())
}

/// Reads a length-prefixed byte slice.
/// Returns Ok(None) on clean EOF.
///
/// Format: [u32 len][bytes...]
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