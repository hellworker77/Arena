package udp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"game-server/internal/udp/udp_types"
	"game-server/pkg/protocol"
	"net"
	"time"
)

func (s *Server) sendPlainPacket(addr *net.UDPAddr, pType protocol.PacketType, payload []byte) {
	// Pre-auth packets are deliberately minimal: no netcode header ack state, no encryption.
	h := protocol.PacketHeader{Ver: 1, Type: uint8(pType)}
	var buf bytes.Buffer
	protocol.WritePacketHeader(&buf, h)
	buf.Write(payload)
	_, _ = s.conn.WriteToUDP(buf.Bytes(), addr)
}

func (s *Server) timeBucketNow() uint32 {
	sec := uint32(time.Now().Unix())
	if s.cookieBucketSec == 0 {
		return sec
	}
	return sec / s.cookieBucketSec
}

func (s *Server) mintCookie(addrStr string, nonce uint64, bucket uint32) [16]byte {
	mac := hmac.New(sha256.New, s.cookieSecret[:])
	mac.Write([]byte(addrStr))
	mac.Write([]byte("|"))
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], nonce)
	mac.Write(tmp[:])
	mac.Write([]byte("|"))
	var tb [4]byte
	binary.LittleEndian.PutUint32(tb[:], bucket)
	mac.Write(tb[:])
	sum := mac.Sum(nil)
	var out [16]byte
	copy(out[:], sum[:16])
	return out
}

func (s *Server) sendEncryptedPacket(addrStr string, pType protocol.PacketType, payload []byte) error {
	s.playerMu.RLock()
	st := s.state[addrStr]
	session := s.playerSessions[addrStr]
	to := s.playerConnections[addrStr]
	s.playerMu.RUnlock()

	if st == nil || session == nil || to == nil {
		return fmt.Errorf("invalid session")
	}

	seq, ackLatest, ackBits := st.PrepareHeaderOnSend()
	h := protocol.PacketHeader{
		Ver:        1,
		Type:       uint8(pType),
		Connection: 0,
		Seq:        seq,
		AckLatest:  ackLatest,
		AckBitmap:  ackBits,
	}

	var buf bytes.Buffer
	protocol.WritePacketHeader(&buf, h)
	buf.Write(payload)

	ciphertext, err := session.EncryptPacket(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = s.conn.WriteToUDP(ciphertext, to)
	return err
}

func (s *Server) handleEncryptedPacket(packet udp_types.InputPacket) {
	addrStr := packet.From.String()

	s.playerMu.RLock()
	session := s.playerSessions[addrStr]
	s.playerMu.RUnlock()

	if session == nil {
		// Pre-auth flood resistance: rate limit *any* plaintext traffic from unknown addresses.
		b := s.unauthLimiters[addrStr]
		if b == nil {
			b = newTokenBucket(s.unauthLimitBurst, s.unauthLimitRate)
			s.unauthLimiters[addrStr] = b
		}
		if !b.allow(1) {
			return
		}

		header, data, err := protocol.ParsePacket(packet.Data)
		if err != nil {
			return
		}

		switch protocol.PacketType(header.Type) {
		case protocol.PTHello:
			s.handleHello(packet.From, data)
		case protocol.PTAuth:
			s.handleAuth(packet.From, data)
		default:
			// Ignore anything else until authenticated.
		}
		return
	}

	plaintext, err := session.DecryptPacket(packet.Data)
	if err != nil {
		return
	}

	header, data, err := protocol.ParsePacket(plaintext)
	if err != nil {
		return
	}

	st := s.state[addrStr]
	st.UpdateOnReceive(header.Seq)

	s.dispatchPacket(packet.From, header, data)
}
