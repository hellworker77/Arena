package udp

import (
	"bytes"
	"fmt"
	"game-server/internal/udp/udp_types"
	"game-server/pkg/protocol"
)

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
		header, data, err := protocol.ParsePacket(packet.Data)
		if err != nil {
			return
		}

		if protocol.PacketType(header.Type) == protocol.PTAuth {
			s.handleAuth(packet.From, data)
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
