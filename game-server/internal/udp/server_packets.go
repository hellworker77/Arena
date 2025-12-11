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
	// TODO: разобрать тело (input)
	// пока просто отправим echo с протоколом

	s.sendEncryptedPacket(addrStr, protocol.PTInput, body)
}
