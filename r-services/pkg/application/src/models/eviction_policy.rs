use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq)]
pub enum EvictionPolicy {
    NoEviction,
    LRU,
    LFU,
    Random,
    Custom,
}