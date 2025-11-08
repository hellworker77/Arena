use std::collections::HashMap;
use serde::Serialize;

#[derive(Serialize)]
pub struct BlobVersionInfo {
    pub size: u64
}

#[derive(Serialize)]
pub struct BlobInfo {
    pub versions: HashMap<String, BlobVersionInfo>,
}