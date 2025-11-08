use tauri::{AppHandle, WebviewUrl, WebviewWindowBuilder};

// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .setup(|app| {

            let handle = app.handle();

            let splashscreen = WebviewWindowBuilder::new(handle, "splashscreen", WebviewUrl::App("splash.html".into()))
                .build()?;

            let main_window = WebviewWindowBuilder::new(handle, "main", WebviewUrl::App("index.html".into()))
                .visible(false)
                .build()?;

            std::thread::spawn(move || {
                std::thread::sleep(std::time::Duration::from_secs(10));
                splashscreen.close().unwrap();
                main_window.show().unwrap();
            });

            Ok(())
        })
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
