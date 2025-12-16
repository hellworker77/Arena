package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// HelloPayload is sent by the client before authentication.
// It allows the server to mint a stateless cookie bound to the sender address.
type HelloPayload struct {
	ClientNonce uint64
}

// ChallengePayload is returned by the server in response to HelloPayload.
// TimeBucket is a coarse timestamp bucket used to limit cookie lifetime.
// Cookie is a truncated MAC.
type ChallengePayload struct {
	TimeBucket uint32
	Cookie     [16]byte
}

// AuthWithCookiePayload is sent by the client to authenticate.
// Token is the JWT access token.
type AuthWithCookiePayload struct {
	ClientNonce uint64
	TimeBucket  uint32
	Cookie      [16]byte
	Token       string
}

func MarshalHello(p HelloPayload) []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, p.ClientNonce)
	return buf.Bytes()
}

func UnmarshalHello(b []byte) (HelloPayload, error) {
	var p HelloPayload
	if len(b) < 8 {
		return p, errors.New("hello payload too short")
	}
	p.ClientNonce = binary.LittleEndian.Uint64(b[:8])
	return p, nil
}

func MarshalChallenge(p ChallengePayload) []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, p.TimeBucket)
	_, _ = buf.Write(p.Cookie[:])
	return buf.Bytes()
}

func UnmarshalChallenge(b []byte) (ChallengePayload, error) {
	var p ChallengePayload
	if len(b) < 4+16 {
		return p, errors.New("challenge payload too short")
	}
	p.TimeBucket = binary.LittleEndian.Uint32(b[:4])
	copy(p.Cookie[:], b[4:20])
	return p, nil
}

// MarshalAuthWithCookie encodes:
// [nonce:8][bucket:4][cookie:16][tokenLen:2][tokenBytes]
func MarshalAuthWithCookie(p AuthWithCookiePayload) ([]byte, error) {
	tok := []byte(p.Token)
	if len(tok) > 0xFFFF {
		return nil, errors.New("token too large")
	}
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, p.ClientNonce)
	_ = binary.Write(&buf, binary.LittleEndian, p.TimeBucket)
	_, _ = buf.Write(p.Cookie[:])
	_ = binary.Write(&buf, binary.LittleEndian, uint16(len(tok)))
	_, _ = buf.Write(tok)
	return buf.Bytes(), nil
}

func UnmarshalAuthWithCookie(b []byte) (AuthWithCookiePayload, error) {
	var p AuthWithCookiePayload
	if len(b) < 8+4+16+2 {
		return p, errors.New("auth payload too short")
	}
	o := 0
	p.ClientNonce = binary.LittleEndian.Uint64(b[o : o+8])
	o += 8
	p.TimeBucket = binary.LittleEndian.Uint32(b[o : o+4])
	o += 4
	copy(p.Cookie[:], b[o:o+16])
	o += 16
	tokLen := int(binary.LittleEndian.Uint16(b[o : o+2]))
	o += 2
	if tokLen < 0 || o+tokLen > len(b) {
		return p, errors.New("invalid token length")
	}
	p.Token = string(b[o : o+tokLen])
	return p, nil
}
