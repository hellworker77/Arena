package shared

import (
	"crypto/rand"
	"encoding/binary"
)

type AccountID uint64
type CharacterID uint64
type SessionID [16]byte
type ZoneID uint32
type EntityID uint32

func NewSessionID() SessionID {
	var id SessionID
	_, _ = rand.Read(id[:])
	return id
}

func (s SessionID) String() string {
	const hex = "0123456789abcdef"
	out := make([]byte, 32)
	for i, b := range s[:] {
		out[i*2] = hex[b>>4]
		out[i*2+1] = hex[b&0x0f]
	}
	return string(out)
}

func PutU64(b []byte, v uint64) { binary.LittleEndian.PutUint64(b, v) }
func U64(b []byte) uint64       { return binary.LittleEndian.Uint64(b) }
func PutU32(b []byte, v uint32) { binary.LittleEndian.PutUint32(b, v) }
func U32(b []byte) uint32       { return binary.LittleEndian.Uint32(b) }
