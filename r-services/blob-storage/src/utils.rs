use sha2::Digest;

pub fn sha256(data: &[u8]) -> [u8; 32] {
    let mut h = sha2::Sha256::new();
    h.update(data);
    h.finalize().into()
}