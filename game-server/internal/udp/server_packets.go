package udp

import (
	"game-server/pkg/protocol"
	"net"
)

func (s *Server) dispatchPacket(addr *net.UDPAddr, header protocol.PacketHeader, data []byte) {
	addrStr := addr.String()

	switch protocol.PacketType(header.Type) {
	case protocol.PTInput:
		s.applyPlayerInput(addrStr, data)
	case protocol.PTReliableCmd:
	case protocol.PTSnapshot:
	default:
	}
}

func (s *Server) applyPlayerInput(addrStr string, body []byte) {
	// Lightweight per-client rate limit (anti-spam). This is a transport-level guard;
	// the simulation loop still validates tick ordering.
	b := s.inputLimiters[addrStr]
	if b == nil {
		b = newTokenBucket(s.inputLimitBurst, s.inputLimitRate)
		s.inputLimiters[addrStr] = b
	}
	if !b.allow(1) {
		return
	}

	in, err := protocol.UnmarshalInput(body)
	if err != nil {
		return
	}
	// Authoritative server: enqueue input; simulation loop consumes it.
	s.loop.QueueInput(addrStr, in)
}
