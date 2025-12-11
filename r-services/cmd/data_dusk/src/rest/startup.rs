use std::env;
use std::net::SocketAddr;
use tokio::net::TcpListener;
use crate::rest::build_app::build_app;
use axum::serve;

pub async fn startup() {
    let app = build_app().await;

    let ip = env::var("HOST").unwrap_or("0.0.0.0".into());
    let port: u16 = env::var("PORT")
        .unwrap_or("3000".into())
        .parse()
        .expect("PORT must be valid u16");

    let addr = SocketAddr::new(ip.parse().unwrap(), port);

    let listener = TcpListener::bind(addr).await.unwrap();
    println!("🚀 Blob server running on http://{}", addr);
    serve(listener, app).await.expect("Server crashed");
}