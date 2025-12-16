package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

const sessionKeyLength = 32

// DeriveHKDFKey derives a key from the given session key using HKDF.
func DeriveHKDFKey(token string, addr string) ([]byte, error) {
	ikm := []byte(token)
	salt := []byte(addr)
	info := []byte("game-session-key-v1")

	hk := hkdf.New(sha256.New, ikm, salt, info)
	key := make([]byte, sessionKeyLength)

	if _, err := io.ReadFull(hk, key); err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateNonce(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
