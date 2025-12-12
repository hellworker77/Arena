use std::env;
use dotenv::dotenv;

pub fn get_config() -> (String, [u8; 32]) {
    dotenv().ok();

    let base_dir = env::var("BASE_DIR").unwrap_or_else(|_| "./data".to_string());

    let codec_key_str = env::var("AES_KEY").expect("AES_KEY must be set in .env");

    let codec_key_bytes = codec_key_str.as_bytes();

    if codec_key_bytes.len() != 32 {
        panic!(
            "AES_KEY must be exactly 32 bytes, found {} bytes",
            codec_key_bytes.len()
        );
    }

    let mut codec_key = [0u8; 32];
    codec_key.copy_from_slice(codec_key_bytes);

    (base_dir, codec_key)
}