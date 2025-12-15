use std::fs::File;
use std::io::{Seek, SeekFrom, Write};
use std::path::Path;

pub struct SegmentWriter {
    file: File,
    offset: u64,
    pub segment_id: u64,
}

impl SegmentWriter {
    pub fn create(path: &Path, segment_id: u64) -> Result<Self, std::io::Error> {
        let mut file = File::create(path)?;
        file.write_all(b"SEG1")?;
        file.write_all(&0u32.to_le_bytes())?;
        Ok(Self {
            file,
            offset: 8,
            segment_id,
        })
    }

    pub fn write_object(
        &mut self,
        hash: [u8; 32],
        nonce: [u8; 12],
        cipher: &[u8],
        size_plain: u64,
    ) -> Result<u64, std::io::Error> {
        let start = self.offset;
        self.file.write_all(&hash)?;
        self.file.write_all(&nonce)?;
        self.file.write_all(&size_plain.to_le_bytes())?;
        self.file.write_all(&(cipher.len() as u64).to_le_bytes())?;
        self.file.write_all(cipher)?;
        self.offset += 32 + 12 + 8 + 8 + cipher.len() as u64;
        Ok(start)
    }

    pub fn finalize(mut self, object_count: u32) -> Result<(), std::io::Error> {
        self.file.seek(SeekFrom::Start(4))?;
        self.file.write_all(&object_count.to_le_bytes())?;
        self.file.sync_all()?;
        Ok(())
    }
}
