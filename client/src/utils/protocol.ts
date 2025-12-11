export enum PacketType {
    PTInput = 1,
    PTSnapshot,
    PTReliableCmd
}

export const HEADER_SIZE = 22

export interface PacketHeader {
    ver: number;
    type: PacketType;
    connection: number;
    seq: number;
    ackLatest: number;
    ackBitmap: bigint;
}

export function encodeHeader(h: PacketHeader): Uint8Array {
    const buf = new ArrayBuffer(HEADER_SIZE)

    const dView = new DataView(buf);

    let offset = 0;

    dView.setUint8(offset ++, h.ver)
    dView.setUint8(offset ++, h.type)

    dView.setUint32(offset, h.connection, true)
    offset += 4

    dView.setUint32(offset, h.seq, true)
    offset += 4

    dView.setUint32(offset, h.ackLatest, true)
    offset += 4

    dView.setBigUint64(offset, h.ackBitmap, true)

    return new Uint8Array(buf)
}

export function encodePacket(header: PacketHeader, body: Uint8Array): Uint8Array {
    const headerBuf = encodeHeader(header)
    const out = new Uint8Array(headerBuf.length + body.length)
    out.set(headerBuf)
    out.set(body, headerBuf.length)
    return out
}

export function decodeHeader(buf: Uint8Array): PacketHeader {
    const view = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
    let offset = 0;

    const ver = view.getUint8(offset++);
    const type = view.getUint8(offset++) as PacketType;

    const connection = view.getUint32(offset, true);
    offset += 4;

    const seq = view.getUint32(offset, true);
    offset += 4;

    const ackLatest = view.getUint32(offset, true);
    offset += 4;

    const ackBitmap = view.getBigUint64(offset, true);

    return { ver, type, connection, seq, ackLatest, ackBitmap };
}

export function decodePacket(data: Uint8Array) {
    const header = decodeHeader(data.subarray(0, HEADER_SIZE));
    const body = data.subarray(HEADER_SIZE);
    return { header, body };
}