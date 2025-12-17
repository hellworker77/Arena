package wire

import (
	"encoding/binary"
	"errors"

	"game-server/internal/persist"
	"game-server/internal/shared"
)

// Attach: [sid:16][cid:u64][zid:u32][interest:u32]
func EncodeAttachPlayer(sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID, interest InterestMask) []byte {
	b := make([]byte, 32)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint64(b[16:24], uint64(cid))
	binary.LittleEndian.PutUint32(b[24:28], uint32(zid))
	binary.LittleEndian.PutUint32(b[28:32], uint32(interest))
	return b
}

func DecodeAttachPlayer(b []byte) (sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID, interest InterestMask, err error) {
	if len(b) != 32 {
		return sid, 0, 0, 0, errors.New("bad attach payload")
	}
	copy(sid[:], b[0:16])
	cid = shared.CharacterID(binary.LittleEndian.Uint64(b[16:24]))
	zid = shared.ZoneID(binary.LittleEndian.Uint32(b[24:28]))
	interest = InterestMask(binary.LittleEndian.Uint32(b[28:32]))
	return
}

// AttachWithState: [sid:16][cid:u64][zid:u32][interest:u32][x:i16][y:i16][hp:u16]
func EncodeAttachWithState(sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID, interest InterestMask, x, y int16, hp uint16) []byte {
	b := make([]byte, 32+2+2+2)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint64(b[16:24], uint64(cid))
	binary.LittleEndian.PutUint32(b[24:28], uint32(zid))
	binary.LittleEndian.PutUint32(b[28:32], uint32(interest))
	binary.LittleEndian.PutUint16(b[32:34], uint16(x))
	binary.LittleEndian.PutUint16(b[34:36], uint16(y))
	binary.LittleEndian.PutUint16(b[36:38], hp)
	return b
}

func DecodeAttachWithState(b []byte) (sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID, interest InterestMask, x, y int16, hp uint16, err error) {
	if len(b) != 38 {
		return sid, 0, 0, 0, 0, 0, 0, errors.New("bad attach-with-state payload")
	}
	copy(sid[:], b[0:16])
	cid = shared.CharacterID(binary.LittleEndian.Uint64(b[16:24]))
	zid = shared.ZoneID(binary.LittleEndian.Uint32(b[24:28]))
	interest = InterestMask(binary.LittleEndian.Uint32(b[28:32]))
	x = int16(binary.LittleEndian.Uint16(b[32:34]))
	y = int16(binary.LittleEndian.Uint16(b[34:36]))
	hp = binary.LittleEndian.Uint16(b[36:38])
	return
}

// Detach: [sid:16]
func EncodeDetachPlayer(sid shared.SessionID) []byte {
	b := make([]byte, 16)
	copy(b, sid[:])
	return b
}
func DecodeDetachPlayer(b []byte) (sid shared.SessionID, err error) {
	if len(b) != 16 { return sid, errors.New("bad detach payload") }
	copy(sid[:], b)
	return
}

// Input: [sid:16][tick:u32][mx:i16][my:i16]
func EncodePlayerInput(sid shared.SessionID, clientTick uint32, mx, my int16) []byte {
	b := make([]byte, 24)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint32(b[16:20], clientTick)
	binary.LittleEndian.PutUint16(b[20:22], uint16(mx))
	binary.LittleEndian.PutUint16(b[22:24], uint16(my))
	return b
}
func DecodePlayerInput(b []byte) (sid shared.SessionID, tick uint32, mx, my int16, err error) {
	if len(b) != 24 { return sid, 0, 0, 0, errors.New("bad input payload") }
	copy(sid[:], b[0:16])
	tick = binary.LittleEndian.Uint32(b[16:20])
	mx = int16(binary.LittleEndian.Uint16(b[20:22]))
	my = int16(binary.LittleEndian.Uint16(b[22:24]))
	return
}

// Action: [sid:16][tick:u32][skill:u16][targetEID:u32]
func EncodePlayerAction(sid shared.SessionID, tick uint32, skill uint16, target shared.EntityID) []byte {
	b := make([]byte, 16+4+2+4)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint32(b[16:20], tick)
	binary.LittleEndian.PutUint16(b[20:22], skill)
	binary.LittleEndian.PutUint32(b[22:26], uint32(target))
	return b
}
func DecodePlayerAction(b []byte) (sid shared.SessionID, tick uint32, skill uint16, target shared.EntityID, err error) {
	if len(b) != 26 { return sid, 0, 0, 0, errors.New("bad action payload") }
	copy(sid[:], b[0:16])
	tick = binary.LittleEndian.Uint32(b[16:20])
	skill = binary.LittleEndian.Uint16(b[20:22])
	target = shared.EntityID(binary.LittleEndian.Uint32(b[22:26]))
	return
}

// Error: [code:u16][msgLen:u16][msgBytes...]
func EncodeError(code ErrCode, msg string) []byte {
	if len(msg) > 65535 { msg = msg[:65535] }
	b := make([]byte, 4+len(msg))
	binary.LittleEndian.PutUint16(b[0:2], uint16(code))
	binary.LittleEndian.PutUint16(b[2:4], uint16(len(msg)))
	copy(b[4:], []byte(msg))
	return b
}
func DecodeError(b []byte) (code ErrCode, msg string, err error) {
	if len(b) < 4 { return ErrUnknown, "", errors.New("bad error payload") }
	code = ErrCode(binary.LittleEndian.Uint16(b[0:2]))
	n := int(binary.LittleEndian.Uint16(b[2:4]))
	if len(b) != 4+n { return ErrUnknown, "", errors.New("bad error payload length") }
	return code, string(b[4:]), nil
}

type RepEvent struct {
	Op  RepOp
	EID shared.EntityID
	X, Y int16
	Val  uint16
	Text string
	Kind EntityKind
	Mask InterestMask
}

// Replicate: [sid:16][serverTick:u32][chan:u8][n:u16] events...
//
// event encodings by op:
// - RepSpawn: [op:u8][eid:u32][kind:u8][mask:u32][x:i16][y:i16]
// - RepMove:  [op:u8][eid:u32][x:i16][y:i16]
// - RepDespawn: [op:u8][eid:u32]
// - RepStateHP: [op:u8][eid:u32][hp:u16]
// - RepEventText: [op:u8][len:u16][bytes...]
func EncodeReplicate(sid shared.SessionID, serverTick uint32, ch RepChannel, events []RepEvent) []byte {
	if len(events) > 65535 { events = events[:65535] }
	sz := 16 + 4 + 1 + 2
	for _, e := range events {
		switch e.Op {
		case RepSpawn:
			sz += 1 + 4 + 1 + 4 + 4
		case RepMove:
			sz += 1 + 4 + 4
		case RepDespawn:
			sz += 1 + 4
		case RepStateHP:
			sz += 1 + 4 + 2
		case RepEventText:
			txt := e.Text
			if len(txt) > 65535 { txt = txt[:65535] }
			sz += 1 + 2 + len(txt)
		}
	}
	b := make([]byte, sz)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint32(b[16:20], serverTick)
	b[20] = byte(ch)
	binary.LittleEndian.PutUint16(b[21:23], uint16(len(events)))
	off := 23
	for _, e := range events {
		switch e.Op {
		case RepSpawn:
			b[off] = byte(e.Op); off++
			binary.LittleEndian.PutUint32(b[off:off+4], uint32(e.EID)); off += 4
			b[off] = byte(e.Kind); off++
			binary.LittleEndian.PutUint32(b[off:off+4], uint32(e.Mask)); off += 4
			binary.LittleEndian.PutUint16(b[off:off+2], uint16(e.X))
			binary.LittleEndian.PutUint16(b[off+2:off+4], uint16(e.Y))
			off += 4
		case RepMove:
			b[off] = byte(e.Op); off++
			binary.LittleEndian.PutUint32(b[off:off+4], uint32(e.EID)); off += 4
			binary.LittleEndian.PutUint16(b[off:off+2], uint16(e.X))
			binary.LittleEndian.PutUint16(b[off+2:off+4], uint16(e.Y))
			off += 4
		case RepDespawn:
			b[off] = byte(e.Op); off++
			binary.LittleEndian.PutUint32(b[off:off+4], uint32(e.EID)); off += 4
		case RepStateHP:
			b[off] = byte(e.Op); off++
			binary.LittleEndian.PutUint32(b[off:off+4], uint32(e.EID)); off += 4
			binary.LittleEndian.PutUint16(b[off:off+2], e.Val); off += 2
		case RepEventText:
			txt := e.Text
			if len(txt) > 65535 { txt = txt[:65535] }
			b[off] = byte(e.Op); off++
			binary.LittleEndian.PutUint16(b[off:off+2], uint16(len(txt))); off += 2
			copy(b[off:], []byte(txt)); off += len(txt)
		}
	}
	return b
}

func DecodeReplicate(b []byte) (sid shared.SessionID, serverTick uint32, ch RepChannel, events []RepEvent, err error) {
	if len(b) < 23 { return sid, 0, 0, nil, errors.New("bad replicate payload") }
	copy(sid[:], b[0:16])
	serverTick = binary.LittleEndian.Uint32(b[16:20])
	ch = RepChannel(b[20])
	n := int(binary.LittleEndian.Uint16(b[21:23]))
	off := 23
	events = make([]RepEvent, 0, n)
	for i := 0; i < n; i++ {
		if off+1 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
		op := RepOp(b[off]); off++
		switch op {
		case RepSpawn:
			if off+4+1+4+4 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			eid := shared.EntityID(binary.LittleEndian.Uint32(b[off:off+4])); off += 4
			kind := EntityKind(b[off]); off++
			mask := InterestMask(binary.LittleEndian.Uint32(b[off:off+4])); off += 4
			x := int16(binary.LittleEndian.Uint16(b[off:off+2]))
			y := int16(binary.LittleEndian.Uint16(b[off+2:off+4])); off += 4
			events = append(events, RepEvent{Op: op, EID: eid, Kind: kind, Mask: mask, X: x, Y: y})
		case RepMove:
			if off+8 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			eid := shared.EntityID(binary.LittleEndian.Uint32(b[off:off+4])); off += 4
			x := int16(binary.LittleEndian.Uint16(b[off:off+2]))
			y := int16(binary.LittleEndian.Uint16(b[off+2:off+4])); off += 4
			events = append(events, RepEvent{Op: op, EID: eid, X: x, Y: y})
		case RepDespawn:
			if off+4 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			eid := shared.EntityID(binary.LittleEndian.Uint32(b[off:off+4])); off += 4
			events = append(events, RepEvent{Op: op, EID: eid})
		case RepStateHP:
			if off+6 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			eid := shared.EntityID(binary.LittleEndian.Uint32(b[off:off+4])); off += 4
			hp := binary.LittleEndian.Uint16(b[off:off+2]); off += 2
			events = append(events, RepEvent{Op: op, EID: eid, Val: hp})
		case RepEventText:
			if off+2 > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			l := int(binary.LittleEndian.Uint16(b[off:off+2])); off += 2
			if off+l > len(b) { return sid, 0, 0, nil, errors.New("bad replicate payload length") }
			txt := string(b[off:off+l]); off += l
			events = append(events, RepEvent{Op: op, Text: txt})
		default:
			return sid, 0, 0, nil, errors.New("unknown replicate op")
		}
	}
	if off != len(b) { return sid, 0, 0, nil, errors.New("extra bytes in replicate payload") }
	return
}

// TransferPrepare: [sid:16][cid:u64][targetZone:u32][interest:u32][x:i16][y:i16][hp:u16]
func EncodeTransferPrepare(sid shared.SessionID, cid shared.CharacterID, target shared.ZoneID, interest InterestMask, st persist.CharacterState) []byte {
	b := make([]byte, 16+8+4+4+2+2+2)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint64(b[16:24], uint64(cid))
	binary.LittleEndian.PutUint32(b[24:28], uint32(target))
	binary.LittleEndian.PutUint32(b[28:32], uint32(interest))
	binary.LittleEndian.PutUint16(b[32:34], uint16(st.X))
	binary.LittleEndian.PutUint16(b[34:36], uint16(st.Y))
	binary.LittleEndian.PutUint16(b[36:38], st.HP)
	return b
}

func DecodeTransferPrepare(b []byte) (sid shared.SessionID, cid shared.CharacterID, target shared.ZoneID, interest InterestMask, x, y int16, hp uint16, err error) {
	if len(b) != 38 {
		return sid, 0, 0, 0, 0, 0, 0, errors.New("bad transfer-prepare payload")
	}
	copy(sid[:], b[0:16])
	cid = shared.CharacterID(binary.LittleEndian.Uint64(b[16:24]))
	target = shared.ZoneID(binary.LittleEndian.Uint32(b[24:28]))
	interest = InterestMask(binary.LittleEndian.Uint32(b[28:32]))
	x = int16(binary.LittleEndian.Uint16(b[32:34]))
	y = int16(binary.LittleEndian.Uint16(b[34:36]))
	hp = binary.LittleEndian.Uint16(b[36:38])
	return
}

// TransferCommit/Abort: [sid:16]
func EncodeTransferCommit(sid shared.SessionID) []byte { return EncodeDetachPlayer(sid) }
func DecodeTransferCommit(b []byte) (shared.SessionID, error) { return DecodeDetachPlayer(b) }
func EncodeTransferAbort(sid shared.SessionID) []byte { return EncodeDetachPlayer(sid) }
func DecodeTransferAbort(b []byte) (shared.SessionID, error) { return DecodeDetachPlayer(b) }
