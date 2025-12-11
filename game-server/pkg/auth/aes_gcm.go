package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type PlayerSession struct {
	Key []byte
}

func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		Key: generateAESKey(),
	}
}

func NewSessionFromToken(token, addr string) (*PlayerSession, error) {
	key, err := DeriveHKDFKey(token, addr)
	if err != nil {
		return nil, err
	}

	return &PlayerSession{Key: key}, nil
}

func generateAESKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	return key
}

func encryptAESGCM(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptAESGCM(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < aesGcm.NonceSize() {
		return nil, fmt.Errorf("malformed ciphertext")
	}

	nonce := data[:aesGcm.NonceSize()]
	ciphertext := data[aesGcm.NonceSize():]

	return aesGcm.Open(nil, nonce, ciphertext, nil)
}

func (ps *PlayerSession) EncryptPacket(plaintext []byte) ([]byte, error) {
	return encryptAESGCM(ps.Key, plaintext)
}

func (ps *PlayerSession) DecryptPacket(plaintext []byte) ([]byte, error) {
	return decryptAESGCM(ps.Key, plaintext)
}
