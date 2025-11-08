use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Message {
    /// Unique message identifier
    pub id: String,
    
    /// Message payload
    pub payload: Vec<u8>,
    
    /// Message headers
    pub headers: std::collections::HashMap<String, String>,
    
    /// Unix timestamp in milliseconds
    pub timestamp: u128, 
    
    /// Number of delivery attempts
    pub attempts: u32,
}