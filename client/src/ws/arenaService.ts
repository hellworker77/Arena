import * as signalR from "@microsoft/signalr";
import type { ArenaState } from "../@types/game/arenaState";

class ArenaService {
    private static instance: ArenaService | null = null;

    private constructor() {} // приватный конструктор

    public static getInstance(): ArenaService {
        if (!ArenaService.instance) {
            ArenaService.instance = new ArenaService();
        }
        return ArenaService.instance;
    }

    private connection: signalR.HubConnection | null = null;
    private listeners: ((state: ArenaState) => void)[] = [];

    public async connect() {
        if (this.connection) return; // уже подключились

        this.connection = new signalR.HubConnectionBuilder()
            .withUrl("http://localhost:5184/arenaHub")
            .withAutomaticReconnect()
            .build();

        this.connection.on("ReceiveSnapshot", (state: ArenaState) => {
            this.listeners.forEach(fn => fn(state));
        });

        await this.connection.start();
        console.log("Connected to ArenaHub");
    }


    public async moveTo(x: number, y: number) {
        if (!this.connection) return;
        try {
            await this.connection.invoke("MoveTo", x, y);
        } catch (err) {
            console.error("MoveTo error:", err);
        }
    }

    public onUpdate(callback: (state: ArenaState) => void) {
        this.listeners.push(callback);
    }

    public async getSnapshot(): Promise<ArenaState | null> {
        if (!this.connection) return null;
        return this.connection.invoke<ArenaState>("GetSnapshot");
    }
}

export const arenaService  = ArenaService.getInstance();