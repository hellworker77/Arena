import Phaser from "phaser";
import {decodePacket, encodePacket, PacketType} from "@/utils/protocol.ts";
import {tauriUDP} from "@/utils/tauriUDP.ts";

export class Play extends Phaser.Scene {
    //private player!: Phaser.GameObjects.Image;
    private seq = 1;
    private shouldCleanup: (() => void)[] = [];

    constructor() {
        super({ key: "play" });
    }

    init() {
        this.events.once("destroy", async () => this.destroy())
    }

    async create() {
        //this.player = this.add.image(100, 100, "player1");

        await tauriUDP.connect("localhost:16054");

        const unlisten = await tauriUDP.onPacket(data => {
            this.handlePacket(data);
        })

        this.shouldCleanup.push(unlisten)

        this.time.addEvent({
            delay: 50,
            loop: true,
            callback: () => this.sendInput(),
        });
    }

    update(_time: number, _dt: number) {}

    private handlePacket(data: Uint8Array) {
        const packet = decodePacket(data)
    }

    private async sendInput() {
        const left = this.input.keyboard?.addKey("A")?.isDown ? 1 : 0
        const right = this.input.keyboard?.addKey("D")?.isDown ? 1 : 0
        const up = this.input.keyboard?.addKey("W")?.isDown ? 1 : 0
        const down = this.input.keyboard?.addKey("S")?.isDown ? 1 : 0

        const body = new Uint8Array([left, right, up, down])

        const packet = encodePacket({
            ver: 1,
            type: PacketType.PTInput,
            connection: 0,
            seq: this.seq++,
            ackLatest: 0,
            ackBitmap: 0n
        }, body)

        await tauriUDP.send(packet)
    }

    private async destroy() {
        this.shouldCleanup.forEach(fn => fn())
        this.shouldCleanup.length = 0
    }
}