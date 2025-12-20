#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    object_store_v2::startup().await
}
