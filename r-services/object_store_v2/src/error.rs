pub mod error {
    use thiserror::Error;

    #[derive(Debug, Error)]
    pub enum StoreError {
        #[error("io error: {0}")]
        Io(#[from] std::io::Error),

        #[error("serialization error: {0}")]
        Serde(#[from] Box<bincode::ErrorKind>),

        #[error("not found")]
        NotFound,

        #[error("object deleted")]
        Deleted,

        #[error("cas miss")]
        CasMiss,

        #[error("segment missing")]
        SegmentMissing,

        #[error("hash mismatch (corruption)")]
        HashMismatch,

        #[error("bad segment magic")]
        BadSegmentMagic,

        #[error("bad sstable magic")]
        BadSstMagic,

        #[error("invalid segment offset")]
        InvalidSegmentOffset,

        #[error("segment scan failed: {0}")]
        SegmentScan(String),

        #[error("manifest missing segment for id {0}")]
        ManifestMissingSegment(u64),

        #[error("cas references unknown object hash")]
        CasDanglingObject,

        #[error("cas references unknown segment id {0}")]
        CasDanglingSegment(u64),
    }

    pub type StoreResult<T> = Result<T, StoreError>;
}