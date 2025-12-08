mod udp;
mod protocol;

use crate::udp::UdpClient;
use std::sync::{Arc, Mutex};
use tauri::{AppHandle, WebviewUrl, WebviewWindowBuilder};
use crate::protocol::PacketType;

#[tauri::command]
fn udp_connect(
    state: tauri::State<Arc<Mutex<Option<Arc<UdpClient>>>>>,
    addr: String,
    handle: AppHandle,
) {
    let client = Arc::new(UdpClient::new());
    client.connect_to(&addr);
    client.clone().start_listener(handle);

    *state.lock().unwrap() = Some(client);
}

#[tauri::command]
fn udp_send(state: tauri::State<Arc<Mutex<Option<Arc<UdpClient>>>>>,
            ptype_u8: u8,
            payload: Vec<u8>) {
    if let Some(client) = &*state.lock().unwrap() {
        let mut buf = Vec::new();
        let ptype: PacketType = PacketType::try_from(ptype_u8).expect("invalid ptype u8");
        let header = protocol::PacketHeader::new(ptype);
        protocol::write_packet(&mut buf, &header, &payload);
        client.send(&buf);
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .manage(Arc::new(Mutex::new(None::<Arc<UdpClient>>)))
        .setup(|app| {
            let handle = app.handle();

            let splashscreen = WebviewWindowBuilder::new(
                handle,
                "splashscreen",
                WebviewUrl::App("splash.html".into()),
            )
            .build()?;

            let main_window =
                WebviewWindowBuilder::new(handle, "main", WebviewUrl::App("index.html".into()))
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
        .invoke_handler(tauri::generate_handler![
            udp_connect,
            udp_send,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
