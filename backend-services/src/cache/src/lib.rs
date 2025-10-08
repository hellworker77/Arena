pub mod cache;
pub mod tasks;

pub mod proto {
    tonic::include_proto!("cache");
}

use tonic::transport::Server;
use cache::service::CacheService;
use tasks::ttl_cleaner::start_ttl_cleaner;
use tokio::net::TcpListener;
use tonic::codegen::tokio_stream;

pub async fn run(addr: String) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let service = CacheService::new();
    let store = service.store.clone();

    tokio::spawn(start_ttl_cleaner(store));

    // Привязываем TcpListener явно
    let listener = TcpListener::bind(&addr).await?;
    let addr = listener.local_addr()?;
    println!("gRPC cache server running on {}", addr);

    Server::builder()
        .add_service(service.into_server())
        .serve_with_incoming(tokio_stream::wrappers::TcpListenerStream::new(listener))
        .await?;

    Ok(())
}