use std::process::Command;

fn clear_console() {
    if cfg!(target_os = "windows") {
        Command::new("cmd").args(&["/C", "cls"]).status().unwrap();
    } else {
        Command::new("clear").status().unwrap();
    }
}

#[tokio::main]
async fn main() {
    clear_console();
    rusty_blob::run().await;
}