package protocol

import (
	"bytes"
	"encoding/binary"
)

type EntitySnapshot struct {
	ID        uint32
	PositionX float32
	PositionY float32
}

func MarshalSnapshotPayload(snapshot []EntitySnapshot) []byte {
	buf := new(bytes.Buffer)

	count := uint16(len(snapshot))
	_ = binary.Write(buf, binary.LittleEndian, count)
	for _, ent := range snapshot {
		_ = binary.Write(buf, binary.LittleEndian, ent.ID)
		_ = binary.Write(buf, binary.LittleEndian, ent.PositionX)
		_ = binary.Write(buf, binary.LittleEndian, ent.PositionY)
	}

	return buf.Bytes()
}
