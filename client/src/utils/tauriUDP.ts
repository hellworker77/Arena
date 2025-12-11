import {listen} from "@tauri-apps/api/event";
import {invoke} from "@tauri-apps/api/core";

class TauriUDP {
    public async connect(addr: string) {
        await invoke("udp_connect", {addr: addr});
    }

    public async send(data: Uint8Array) {
        await invoke("udp_send", {
            data: Array.from(data)
        });
    }

    public async onPacket(cb: (data: Uint8Array) => void) {
        return await listen<number[]>("udp://packet", (event) => {
            const data = new Uint8Array(event.payload);
            cb(data);
        })
    }
}

export const tauriUDP = new TauriUDP();