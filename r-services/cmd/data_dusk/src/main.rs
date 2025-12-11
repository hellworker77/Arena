use crate::rest::startup::startup;

mod http_server;
pub(crate) mod get_config;
mod rest;
pub(crate) mod types;

#[tokio::main]
async fn main() {
    startup().await
}