package wire

import (
	"encoding/binary"
	"errors"

	"game-server/internal/shared"
)

// MsgAttachPlayer payload: [sessionID:16][characterID:uint64][zoneID:uint32]
func EncodeAttachPlayer(sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID) []byte {
	b := make([]byte, 16+8+4)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint64(b[16:24], uint64(cid))
	binary.LittleEndian.PutUint32(b[24:28], uint32(zid))
	return b
}

func DecodeAttachPlayer(b []byte) (sid shared.SessionID, cid shared.CharacterID, zid shared.ZoneID, err error) {
	if len(b) != 28 {
		return sid, 0, 0, errors.New("bad attach payload")
	}
	copy(sid[:], b[0:16])
	cid = shared.CharacterID(binary.LittleEndian.Uint64(b[16:24]))
	zid = shared.ZoneID(binary.LittleEndian.Uint32(b[24:28]))
	return
}

// MsgDetachPlayer payload: [sessionID:16]
func EncodeDetachPlayer(sid shared.SessionID) []byte {
	b := make([]byte, 16)
	copy(b, sid[:])
	return b
}

func DecodeDetachPlayer(b []byte) (sid shared.SessionID, err error) {
	if len(b) != 16 {
		return sid, errors.New("bad detach payload")
	}
	copy(sid[:], b)
	return
}

// MsgPlayerInput payload: [sessionID:16][clientTick:uint32][mx:int16][my:int16]
func EncodePlayerInput(sid shared.SessionID, clientTick uint32, mx, my int16) []byte {
	b := make([]byte, 16+4+2+2)
	copy(b[0:16], sid[:])
	binary.LittleEndian.PutUint32(b[16:20], clientTick)
	binary.LittleEndian.PutUint16(b[20:22], uint16(mx))
	binary.LittleEndian.PutUint16(b[22:24], uint16(my))
	return b
}

func DecodePlayerInput(b []byte) (sid shared.SessionID, tick uint32, mx, my int16, err error) {
	if len(b) != 24 {
		return sid, 0, 0, 0, errors.New("bad input payload")
	}
	copy(sid[:], b[0:16])
	tick = binary.LittleEndian.Uint32(b[16:20])
	mx = int16(binary.LittleEndian.Uint16(b[20:22]))
	my = int16(binary.LittleEndian.Uint16(b[22:24]))
	return
}

// MsgError payload: [code:uint16][msgLen:uint16][msgBytes...]
func EncodeError(code ErrCode, msg string) []byte {
	if len(msg) > 65535 {
		msg = msg[:65535]
	}
	b := make([]byte, 2+2+len(msg))
	binary.LittleEndian.PutUint16(b[0:2], uint16(code))
	binary.LittleEndian.PutUint16(b[2:4], uint16(len(msg)))
	copy(b[4:], []byte(msg))
	return b
}

func DecodeError(b []byte) (code ErrCode, msg string, err error) {
	if len(b) < 4 {
		return ErrUnknown, "", errors.New("bad error payload")
	}
	code = ErrCode(binary.LittleEndian.Uint16(b[0:2]))
	n := int(binary.LittleEndian.Uint16(b[2:4]))
	if len(b) != 4+n {
		return ErrUnknown, "", errors.New("bad error payload length")
	}
	msg = string(b[4:])
	return
}

// MsgSnapshot payload (toy): [serverTick:uint32][entityCount:uint16] repeated: [x:int16][y:int16]
func EncodeSnapshot(serverTick uint32, positions [][2]int16) []byte {
	if len(positions) > 65535 {
		positions = positions[:65535]
	}
	b := make([]byte, 4+2+len(positions)*4)
	binary.LittleEndian.PutUint32(b[0:4], serverTick)
	binary.LittleEndian.PutUint16(b[4:6], uint16(len(positions)))
	off := 6
	for _, p := range positions {
		binary.LittleEndian.PutUint16(b[off:off+2], uint16(p[0]))
		binary.LittleEndian.PutUint16(b[off+2:off+4], uint16(p[1]))
		off += 4
	}
	return b
}

func DecodeSnapshot(b []byte) (serverTick uint32, positions [][2]int16, err error) {
	if len(b) < 6 {
		return 0, nil, errors.New("bad snapshot payload")
	}
	serverTick = binary.LittleEndian.Uint32(b[0:4])
	n := int(binary.LittleEndian.Uint16(b[4:6]))
	expect := 6 + n*4
	if len(b) != expect {
		return 0, nil, errors.New("bad snapshot payload length")
	}
	positions = make([][2]int16, n)
	off := 6
	for i := 0; i < n; i++ {
		x := int16(binary.LittleEndian.Uint16(b[off : off+2]))
		y := int16(binary.LittleEndian.Uint16(b[off+2 : off+4]))
		positions[i] = [2]int16{x, y}
		off += 4
	}
	return
}
