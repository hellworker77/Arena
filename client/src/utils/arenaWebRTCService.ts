import * as signalR from "@microsoft/signalr";
import type { ArenaState } from "../@types/game/arenaState";

// ---------------------------
// Types
// ---------------------------

// Payload for "move" messages
type MovePayload = { x: number; y: number };

// Payload for "snapshotResponse" messages
type SnapshotResponsePayload = ArenaState;

// Union of all possible messages
type MessagePayload =
    | { type: "move"; payload: MovePayload }
    | { type: "getSnapshot" }
    | { type: "snapshotResponse"; payload: SnapshotResponsePayload };

// Messages that the client can send
type ClientMessage = Extract<MessagePayload, { type: "move" | "getSnapshot" }>;

// ---------------------------
// ArenaService
// ---------------------------
class ArenaService {
    private static instance: ArenaService | null = null;

    /** SignalR connection for WebRTC signaling */
    private connection: signalR.HubConnection | null = null;

    /** WebRTC peer connection */
    private peerConnection: RTCPeerConnection | null = null;

    /** DataChannel for sending/receiving game messages */
    private dataChannel: RTCDataChannel | null = null;

    /** Listeners for arena state updates */
    private listeners: ((state: ArenaState) => void)[] = [];

    /** Private constructor for singleton pattern */
    private constructor() {}

    /** Get singleton instance of ArenaService */
    public static getInstance(): ArenaService {
        if (!ArenaService.instance) {
            ArenaService.instance = new ArenaService();
        }
        return ArenaService.instance;
    }

    // ---------------------------
    // Connect to SignalR and initialize WebRTC
    // ---------------------------
    public async connect() {
        if (this.connection) return;

        // Initialize SignalR connection for signaling
        this.connection = new signalR.HubConnectionBuilder()
            .withUrl("http://localhost:5184/arenaSignalHub")
            .withAutomaticReconnect()
            .build();

        /** Handle SDP answer from server */
        this.connection.on("AnswerFromServer", async (sdpAnswer: string) => {
            if (!this.peerConnection) return;

            const answer: RTCSessionDescriptionInit = {
                type: "answer" as RTCSdpType,
                sdp: sdpAnswer,
            };

            await this.peerConnection.setRemoteDescription(answer);
        });

        /** Handle ICE candidates from server */
        this.connection.on("ServerIceCandidate", async (candidateJson: string) => {
            if (!this.peerConnection) return;
            const candidate = JSON.parse(candidateJson) as RTCIceCandidateInit;
            await this.peerConnection.addIceCandidate(candidate);
        });

        await this.connection.start();
        console.log("Connected to ArenaSignalHub");

        await this.initWebRtc();
    }

    // ---------------------------
    // Initialize WebRTC PeerConnection and DataChannel
    // ---------------------------
    private async initWebRtc() {
        this.peerConnection = new RTCPeerConnection({
            iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
        });

        // Create DataChannel for game messages
        this.dataChannel = this.peerConnection.createDataChannel("arena");
        this.dataChannel.onopen = () => console.log("DataChannel opened");
        this.dataChannel.onclose = () => console.log("DataChannel closed");

        /** Handle incoming messages from server */
        this.dataChannel.onmessage = (event) => {
            const msg = JSON.parse(event.data) as MessagePayload;

            if (msg.type === "snapshotResponse") {
                this.listeners.forEach((fn) => fn(msg.payload));
            }
        };

        /** Handle ICE candidates locally and send to server */
        this.peerConnection.onicecandidate = (event) => {
            if (event.candidate) {
                this.connection?.invoke(
                    "SendIceCandidateToServer",
                    JSON.stringify(event.candidate)
                );
            }
        };

        // Create SDP offer
        const offer = await this.peerConnection.createOffer();
        await this.peerConnection.setLocalDescription(offer);

        // Send SDP offer to server via SignalR
        await this.connection?.invoke("OfferToServer", offer.sdp!);
    }

    // ---------------------------
    // Send a message to server via DataChannel
    // ---------------------------
    public sendMessage(message: ClientMessage) {
        if (!this.dataChannel || this.dataChannel.readyState !== "open") return;
        this.dataChannel.send(JSON.stringify(message));
    }

    /** Move player to a new position */
    public moveTo(x: number, y: number) {
        this.sendMessage({ type: "move", payload: { x, y } });
    }

    /** Request arena snapshot from server */
    public getSnapshot() {
        this.sendMessage({ type: "getSnapshot" });
    }

    /** Subscribe to arena state updates */
    public onUpdate(callback: (state: ArenaState) => void) {
        this.listeners.push(callback);
    }
}

// ---------------------------
// Export singleton instance
// ---------------------------
export const arenaWebRTCService = ArenaService.getInstance();