mod wal;
mod segment_writer;
mod cas_index;
mod key_index;
mod gc_worker;
mod compactor;
mod utils;
mod object_store;
mod manifest;
mod key_index_store;
mod cas_index_store;
mod recovery;

#[tokio::main]
async fn main() {
}


/*
src/
├── lib.rs
├── error.rs
├── config.rs
├── metrics.rs
src/storage/
├── mod.rs
├── engine.rs
├── manifest.rs
├── layout.rs
src/wal/
├── mod.rs
├── record.rs
├── writer.rs
├── reader.rs
├── replay.rs
src/segment/
├── mod.rs
├── format.rs
├── writer.rs
├── reader.rs
├── manager.rs
src/cas/
├── mod.rs
├── index.rs
├── entry.rs
├── persistent.rs
src/index/
├── mod.rs
├── entry.rs
├── memtable.rs
├── sstable.rs
├── compaction.rs
src/gc/
├── mod.rs
├── marker.rs
├── sweeper.rs
├── policy.rs
src/compaction/
├── mod.rs
├── planner.rs
├── executor.rs
src/crypto/
├── mod.rs
├── aes_gcm.rs
├── nonce.rs
src/codec/
├── mod.rs
├── zlib.rs
├── encoder.rs
├── decoder.rs
src/api/
├── mod.rs
├── blob_repository.rs
├── s3_like.rs

*/