package protocol

import (
	"bytes"
	"encoding/binary"
)

// ---- Legacy snapshot (kept for compatibility inside the server codebase) ----

// EntitySnapshot is the old minimal snapshot format.
// NOTE: New code should prefer EntityState + MarshalSnapshotFull/Delta.
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

// ---- New snapshot format (versioned payload, supports full + delta) ----

type SnapshotMode uint8

const (
	SnapshotFull  SnapshotMode = 0
	SnapshotDelta SnapshotMode = 1
)

// EntityState is the state replicated to clients.
// Keep this small; everything here is paid per entity per snapshot.
type EntityState struct {
	ID uint32
	X  float32
	Y  float32
	VX float32
	VY float32
}

// Payload format (little-endian):
//   mode:uint8
//   serverTick:uint32
//   snapshotSeq:uint32
//   if mode==Full:
//       count:uint16
//       repeated EntityState
//   if mode==Delta:
//       upserts:uint16
//       repeated EntityState
//       removes:uint16
//       repeated id:uint32

func MarshalSnapshotFull(serverTick, snapshotSeq uint32, states []EntityState) []byte {
	buf := new(bytes.Buffer)
	_ = buf.WriteByte(byte(SnapshotFull))
	_ = binary.Write(buf, binary.LittleEndian, serverTick)
	_ = binary.Write(buf, binary.LittleEndian, snapshotSeq)
	_ = binary.Write(buf, binary.LittleEndian, uint16(len(states)))
	for _, s := range states {
		_ = binary.Write(buf, binary.LittleEndian, s.ID)
		_ = binary.Write(buf, binary.LittleEndian, s.X)
		_ = binary.Write(buf, binary.LittleEndian, s.Y)
		_ = binary.Write(buf, binary.LittleEndian, s.VX)
		_ = binary.Write(buf, binary.LittleEndian, s.VY)
	}
	return buf.Bytes()
}

func MarshalSnapshotDelta(serverTick, snapshotSeq uint32, upserts []EntityState, removes []uint32) []byte {
	buf := new(bytes.Buffer)
	_ = buf.WriteByte(byte(SnapshotDelta))
	_ = binary.Write(buf, binary.LittleEndian, serverTick)
	_ = binary.Write(buf, binary.LittleEndian, snapshotSeq)
	_ = binary.Write(buf, binary.LittleEndian, uint16(len(upserts)))
	for _, s := range upserts {
		_ = binary.Write(buf, binary.LittleEndian, s.ID)
		_ = binary.Write(buf, binary.LittleEndian, s.X)
		_ = binary.Write(buf, binary.LittleEndian, s.Y)
		_ = binary.Write(buf, binary.LittleEndian, s.VX)
		_ = binary.Write(buf, binary.LittleEndian, s.VY)
	}
	_ = binary.Write(buf, binary.LittleEndian, uint16(len(removes)))
	for _, id := range removes {
		_ = binary.Write(buf, binary.LittleEndian, id)
	}
	return buf.Bytes()
}
