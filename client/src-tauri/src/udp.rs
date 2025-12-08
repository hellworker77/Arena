use std::net::UdpSocket;
use std::sync::{Arc, Mutex};
use std::thread;

use tauri::{AppHandle, Emitter, Manager};

pub struct UdpClient {
    socket: UdpSocket,
}

impl UdpClient {
    pub fn new() -> Self {
        let socket = UdpSocket::bind("0.0.0.0:0").expect("bind failed");
        socket
            .set_nonblocking(true)
            .expect("set_nonblocking failed");

        UdpClient { socket }
    }

    pub fn connect_to(&self, server: &str) {
        self.socket
            .connect(server)
            .expect("connect failed");
    }

    pub fn send(&self, data: &[u8]) {
        let _ = self.socket.send(data);
    }

    pub fn start_listener(self: Arc<Self>, app: AppHandle) {
        thread::spawn(move || {
            let mut buf = [0u8; 2048];

            loop {
                if let Ok(n) = self.socket.recv(&mut buf) {
                    let packet = buf[..n].to_vec();

                    let _ = app.emit("udp://packet", packet);
                }

                thread::sleep(std::time::Duration::from_millis(1));
            }
        });
    }
}