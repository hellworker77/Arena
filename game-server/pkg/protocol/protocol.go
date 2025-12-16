package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

type PacketType uint8

const (
	// Handshake packets are sent in plaintext (pre-session).
	PTHello       PacketType = 8
	PTChallenge   PacketType = 9
	PTInput       PacketType = 1
	PTSnapshot    PacketType = 2
	PTReliableCmd PacketType = 3
	PTAuth        PacketType = 10
	PTAuthResp    PacketType = 11
)

type PacketHeader struct {
	Ver        uint8
	Type       uint8
	Connection uint32
	Seq        uint32
	AckLatest  uint32
	AckBitmap  uint64
}

func WritePacketHeader(buf *bytes.Buffer, h PacketHeader) {
	_ = buf.WriteByte(h.Ver)
	_ = buf.WriteByte(h.Type)
	binary.Write(buf, binary.LittleEndian, h.Connection)
	binary.Write(buf, binary.LittleEndian, h.Seq)
	binary.Write(buf, binary.LittleEndian, h.AckLatest)
	binary.Write(buf, binary.LittleEndian, h.AckBitmap)
}

func readHeader(r io.Reader) (PacketHeader, error) {
	var h PacketHeader
	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return h, err
	}
	h.Ver = b[0]
	h.Type = b[1]
	if err := binary.Read(r, binary.LittleEndian, &h.Connection); err != nil {
		return h, err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Seq); err != nil {
		return h, err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.AckLatest); err != nil {
		return h, err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.AckBitmap); err != nil {
		return h, err
	}
	return h, nil
}

func ParsePacket(data []byte) (PacketHeader, []byte, error) {
	r := bytes.NewReader(data)

	h, err := readHeader(r)
	if err != nil {
		return h, nil, err
	}

	body, _ := io.ReadAll(r)
	return h, body, nil
}
