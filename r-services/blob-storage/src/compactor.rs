use crate::cas_index::{CasEntry, CasIndex};
use crate::segment_writer::SegmentWriter;
use std::collections::HashSet;
use std::fs::File;
use std::io::{Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};

pub struct Compactor;

impl Compactor {
    pub fn compact_segments(
        live: &HashSet<[u8; 32]>,
        cas: &mut CasIndex,
        old_segments: &[PathBuf],
        new_segment_path: &Path,
        new_segment_id: u64,
    ) -> Result<(), std::io::Error> {
        let mut writer = SegmentWriter::create(new_segment_path, new_segment_id)?;
        let mut count = 0u32;

        for (hash, entry) in cas.map.clone() {
            if !live.contains(&hash) {
                continue;
            }

            let mut f = File::open(&old_segments[entry.segment_id as usize])?;
            f.seek(SeekFrom::Start(entry.offset))?;

            let mut buf = vec![0u8; entry.size as usize];
            f.read_exact(&mut buf)?;

            //write to new segment
            let new_offset = writer.write_object(hash, [0u8; 12], &buf, entry.size)?;

            //update CAS
            cas.insert(
                hash,
                CasEntry {
                    segment_id: new_segment_id,
                    offset: new_offset,
                    size: entry.size,
                    refcount: entry.refcount,
                },
            );
            count += 1;
        }

        writer.finalize(count)?;
        Ok(())
    }
}
