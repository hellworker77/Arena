use crate::codec::zlib_aes_codec::ZlibAesCodec;
use anyhow::Result;
use application::codec::blob_codec::BlobCodec;
use application::feature::builder_metadata_builder::BlobMetadataBuilder;
use application::repository::blob_repository::BlobRepository;
use async_trait::async_trait;
use chrono::Utc;
use domain::models::blob_metadata::BlobMetadata;
use domain::models::compression_kind::CompressionKind;
use domain::models::encryption_kind::EncryptionKind;
use domain::models::pak_index_entry::PakIndexEntry;
use mime_guess::MimeGuess;
use sha2::{Digest, Sha256};
use std::collections::HashMap;
use std::fs::{create_dir_all, File};
use std::io::{Read, Seek, SeekFrom, Write};
use std::path::{Path, PathBuf};
use uuid::Uuid;

/// Repository for managing PAK files and their entries.
/// Allows retrieval of blob entries by key.
#[derive(Debug)]
pub struct PakRepository {
    pub base_dir: PathBuf,
    pub codec: ZlibAesCodec,
}

impl PakRepository {
    pub fn new(base_dir: impl Into<PathBuf>, key: [u8; 32]) -> Result<Self> {
        let p = base_dir.into();
        create_dir_all(&p)?;
        Ok(Self {
            base_dir: p,
            codec: ZlibAesCodec::new(key),
        })
    }

    fn file_path(&self, key: &str) -> PathBuf {
        let mut p = self.base_dir.clone();
        p.push(format!("{key}.pak"));
        p
    }

    /// Read and parse PAK file
    fn read_pak(&self, path: &Path) -> Result<(Vec<PakIndexEntry>, Vec<u8>, Vec<BlobMetadata>)> {
        if !path.exists() {
            return Ok((vec![], vec![], vec![]));
        }

        let mut f = File::open(path)?;

        // magic
        let mut magic = [0u8; 4];
        f.read_exact(&mut magic)?;
        if &magic != b"PAK1" {
            anyhow::bail!("bad pak magic");
        }

        let mut count_b = [0u8; 4];
        f.read_exact(&mut count_b)?;
        let count = u32::from_le_bytes(count_b);

        let mut entries = Vec::with_capacity(count as usize);

        for _ in 0..count {
            let mut len_b = [0u8; 2];
            f.read_exact(&mut len_b)?;
            let key_len = u16::from_le_bytes(len_b) as usize;

            let mut key_buf = vec![0u8; key_len];
            f.read_exact(&mut key_buf)?;
            let key = String::from_utf8(key_buf)?;

            let mut version_b = [0u8; 4];
            f.read_exact(&mut version_b)?;
            let version = u32::from_le_bytes(version_b);

            let mut offset_b = [0u8; 8];
            f.read_exact(&mut offset_b)?;
            let offset = u64::from_le_bytes(offset_b);

            let mut orig_b = [0u8; 8];
            f.read_exact(&mut orig_b)?;
            let size_original = u64::from_le_bytes(orig_b);

            let mut comp_b = [0u8; 8];
            f.read_exact(&mut comp_b)?;
            let size_compressed = u64::from_le_bytes(comp_b);

            let mut nonce = [0u8; 12];
            f.read_exact(&mut nonce)?;

            let mut uuid_b = [0u8; 16];
            f.read_exact(&mut uuid_b)?;
            let blob_id = Uuid::from_bytes(uuid_b);

            entries.push(PakIndexEntry {
                key,
                version,
                offset,
                size_original,
                size_compressed,
                nonce,
                blob_id,
            });
        }

        // metadata length (last 8 bytes)
        let file_len = f.seek(SeekFrom::End(0))?;
        f.seek(SeekFrom::End(-8))?;
        let mut ml = [0u8; 8];
        f.read_exact(&mut ml)?;
        let meta_len = u64::from_le_bytes(ml);

        if meta_len + 8 > file_len {
            anyhow::bail!("invalid metadata length");
        }

        let meta_start = file_len - 8 - meta_len;
        f.seek(SeekFrom::Start(meta_start))?;
        let mut meta_buf = vec![0u8; meta_len as usize];
        f.read_exact(&mut meta_buf)?;
        let metadata: Vec<BlobMetadata> = serde_json::from_slice(&meta_buf)?;

        // blob data is everything between index and metadata
        let index_end = 4
            + 4
            + entries
                .iter()
                .map(|e| 2 + e.key.len() + 4 + 8 + 8 + 8 + 12 + 16) // corrected: size_compressed = u64
                .sum::<usize>() as u64;

        f.seek(SeekFrom::Start(index_end))?;
        let blob_len = (meta_start - index_end) as usize;
        let mut blob_data = vec![0u8; blob_len];
        f.read_exact(&mut blob_data)?;

        Ok((entries, blob_data, metadata))
    }

    fn write_pak(
        &self,
        path: &Path,
        entries: &[PakIndexEntry],
        blob_data: &[u8],
        metadata: &[BlobMetadata],
    ) -> Result<()> {
        let tmp = path.with_extension("tmp");
        if let Some(parent) = tmp.parent() {
            create_dir_all(parent)?;
        }

        let mut out = File::create(&tmp)?;
        out.write_all(b"PAK1")?;
        out.write_all(&(entries.len() as u32).to_le_bytes())?;

        for e in entries {
            let kb = e.key.as_bytes();
            out.write_all(&(kb.len() as u16).to_le_bytes())?;
            out.write_all(kb)?;
            out.write_all(&e.version.to_le_bytes())?;
            out.write_all(&e.offset.to_le_bytes())?;
            out.write_all(&e.size_original.to_le_bytes())?;
            out.write_all(&e.size_compressed.to_le_bytes())?;
            out.write_all(&e.nonce)?;
            out.write_all(e.blob_id.as_bytes())?;
        }

        out.write_all(blob_data)?;

        let meta_json = serde_json::to_vec_pretty(metadata)?;
        out.write_all(&meta_json)?;
        out.write_all(&(meta_json.len() as u64).to_le_bytes())?;

        out.sync_all()?;
        std::fs::rename(&tmp, &path)?;

        Ok(())
    }
}

#[async_trait]
impl BlobRepository for PakRepository {
    async fn put(&self, key: &str, content: &[u8]) -> Result<()> {
        let path = self.file_path(key);

        // Read existing PAK file
        let (mut entries, mut blob_data, mut metadata) = self.read_pak(&path)?;

        // Detect current version
        let version = entries.iter().map(|e| e.version).max().unwrap_or(0) + 1;

        // Generate blob_id
        let blob_id = Uuid::new_v4();

        // Mime (path -> content -> fallback
        let content_type = MimeGuess::from_path(key)
            .first()
            .map(|m| m.essence_str().to_string())
            .or_else(|| infer::get(content).map(|t| t.mime_type().to_string()))
            .unwrap_or_else(|| "application/octet-stream".to_string());

        // SHA256 from plaintext
        let mut hasher = Sha256::new();
        hasher.update(content);
        let sha256_plain = hex::encode(hasher.finalize());

        // Encode (compress + encrypt)
        let (cipher_buf, nonce, compressed_len, encrypted_len) = self.codec.encode(content)?;

        // Offset inside PAK
        let pak_offset = blob_data.len() as u64;

        // Append blob_data
        blob_data.extend_from_slice(&cipher_buf);

        // Index entry (minimal without less)
        let entry = PakIndexEntry {
            key: key.to_string(),
            version,
            offset: pak_offset,
            size_original: content.len() as u64,
            size_compressed: encrypted_len as u64,
            nonce,
            blob_id,
        };

        entries.push(entry);

        // Etag (simple version-based)
        let etag = format!("\"sha256:{}\"", sha256_plain);

        // Metadata (extensible, API- oriented)
        let meta = BlobMetadataBuilder::new()
            .with_blob_id(blob_id)
            .with_version(version)
            .with_key(key.to_string())
            .with_blob_type("default".to_string())
            .with_created_at_unix(Utc::now().timestamp())
            .with_updated_at_unix(Utc::now().timestamp())
            .with_size_original(content.len() as u64)
            .with_size_compressed(compressed_len as u64)
            .with_size_encrypted(encrypted_len as u64)
            .with_pak_offset(pak_offset)
            .with_encryption_nonce(nonce)
            .with_compression(CompressionKind::Zlib)
            .with_encryption(EncryptionKind::Aes256Gcm)
            .with_content_type(content_type)
            .with_content_checksum_sha256(sha256_plain)
            .with_content_etag(etag)
            .build()
            .map_err(|e| anyhow::anyhow!(e))?;

        // Metadata append-only
        metadata.push(meta);

        self.write_pak(&path, &entries, &blob_data, &metadata)?;

        Ok(())
    }

    async fn get(&self, key: &str) -> Result<Vec<u8>> {
        let path = self.file_path(key);
        let (entries, blob_data, _meta) = self.read_pak(&path)?;

        let e = entries
            .iter()
            .max_by_key(|e| e.version)
            .ok_or_else(|| anyhow::anyhow!("Not found"))?;
        let slice = &blob_data[e.offset as usize..(e.offset + e.size_compressed) as usize];
        let plain = self.codec.decode(&e.nonce, slice)?;

        Ok(plain)
    }

    async fn get_metadata(&self, key: &str) -> Result<BlobMetadata> {
        let path = self.file_path(key);
        let (entries, _blob_data, metadata) = self.read_pak(&path)?;
        let e = entries
            .iter()
            .max_by_key(|e| e.version)
            .ok_or_else(|| anyhow::anyhow!("Not found"))?;
        let m = metadata
            .iter()
            .find(|m| m.version == e.version)
            .ok_or_else(|| anyhow::anyhow!("Metadata missing"))?;
        Ok(m.clone())
    }

    async fn delete(&self, key: &str) -> Result<()> {
        anyhow::bail!("Delete not supported in PakRepository");
    }

    async fn exists(&self, key: &str) -> Result<bool> {
        Ok(self.file_path(key).exists())
    }

    async fn list(&self, prefix: Option<&str>) -> Result<HashMap<String, BlobMetadata>> {
        let mut out = HashMap::new();
        for entry in std::fs::read_dir(&self.base_dir)? {
            let p = entry?.path();
            if p.extension().and_then(|x| x.to_str()) != Some("pak") {
                continue;
            }
            let file_stem = p.file_stem().unwrap().to_string_lossy().to_string();
            if let Some(pref) = prefix {
                if !file_stem.starts_with(pref) {
                    continue;
                }
            }
            let (_entries, _data, meta) = self.read_pak(&p)?;
            if let Some(latest) = meta.iter().max_by_key(|m| m.version) {
                out.insert(file_stem, latest.clone());
            }
        }
        Ok(out)
    }

    async fn copy(&self, source_key: &str, destination_key: &str) -> Result<()> {
        anyhow::bail!("Copy not supported in PakRepository");
    }

    async fn r#move(&self, source_key: &str, destination_key: &str) -> Result<()> {
        anyhow::bail!("Move not supported in PakRepository");
    }
}

/*
┌──────────────────────────────┐
│  magic: "PAK1"               │ 4 bytes
├──────────────────────────────┤
│  version_count (u32)         │
├──────────────────────────────┤
│  index entries [...]         │ variable
│   {
│      key_len (u16)
│      key_bytes
│      version (u32)
│      offset  (u64)
│      size_original   (u64)
│      size_compressed (u64)
│      nonce[12]
│      uuid[16]
│   }
├──────────────────────────────┤
│  BLOB DATA (all encrypted versions)
├──────────────────────────────┤
│  metadata_json               │
├──────────────────────────────┤
│  metadata_json_length (u64)  │ last 8 bytes
└──────────────────────────────┘
*/
