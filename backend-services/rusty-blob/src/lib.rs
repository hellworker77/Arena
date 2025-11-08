use crate::core::application::contracts::data::config::Config;
use crate::inra::infrastructure::factories::build_repositories::build_repositories;
use crate::inra::infrastructure::services::blob_service::BlobServiceImpl;
use crate::presentation::startup::startup;
use std::sync::Arc;
use tracing_subscriber::{EnvFilter};
use crate::inra::infrastructure::factories::build_services::build_services;

pub mod core;
mod inra;
mod presentation;

pub async fn run() {
    let config = Config::load("rusty-blob/config.toml").unwrap();

    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::new(&config.logging_level))
        .init();

    let repository = build_repositories(&config);

    let service =  build_services(repository, &config);
    
    match config.server_mode.as_str() {
        "rest" => {
            startup(service, &config).await;
        }
        _ => {
            eprintln!("Unsupported server mode: {}", config.server_mode);
        }
    }
}