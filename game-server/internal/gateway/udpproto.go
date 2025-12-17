package gateway

import (
	"encoding/binary"
	"errors"
)

const (
	UdpMagic uint16 = 0x4D4D // 'MM'
)

// Channels
const (
	ChanUnreliable uint8 = 0
	ChanReliable   uint8 = 1
)

// Payload types
const (
	PHello  uint8 = 1
	PInput  uint8 = 2
	PAction uint8 = 3
	PText   uint8 = 4
	PRep    uint8 = 5 // replicate line (demo)
)

// Packet:
// [magic:u16][proto:u16][chan:u8][ptype:u8][seq:u32][ack:u32][ackBits:u32][payload...]
const HeaderLen = 2+2+1+1+4+4+4

type Packet struct {
	Proto   uint16
	Chan    uint8
	PType   uint8
	Seq     uint32
	Ack     uint32
	AckBits uint32
	Payload []byte
}

func EncodePacket(p Packet, dst []byte) []byte {
	n := HeaderLen + len(p.Payload)
	if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}
	binary.LittleEndian.PutUint16(dst[0:2], UdpMagic)
	binary.LittleEndian.PutUint16(dst[2:4], p.Proto)
	dst[4] = p.Chan
	dst[5] = p.PType
	binary.LittleEndian.PutUint32(dst[6:10], p.Seq)
	binary.LittleEndian.PutUint32(dst[10:14], p.Ack)
	binary.LittleEndian.PutUint32(dst[14:18], p.AckBits)
	copy(dst[18:], p.Payload)
	return dst
}

func DecodePacket(b []byte) (Packet, error) {
	if len(b) < HeaderLen {
		return Packet{}, errors.New("short packet")
	}
	if binary.LittleEndian.Uint16(b[0:2]) != UdpMagic {
		return Packet{}, errors.New("bad magic")
	}
	p := Packet{
		Proto: binary.LittleEndian.Uint16(b[2:4]),
		Chan: b[4],
		PType: b[5],
		Seq: binary.LittleEndian.Uint32(b[6:10]),
		Ack: binary.LittleEndian.Uint32(b[10:14]),
		AckBits: binary.LittleEndian.Uint32(b[14:18]),
		Payload: b[18:],
	}
	return p, nil
}
