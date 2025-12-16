package udp

import (
	"fmt"
	"game-server/pkg/protocol"
)

func (s *Server) broadcastSnapshot(payload []byte) {
	s.playerMu.RLock()
	for addrStr := range s.playerConnections {
		if err := s.sendEncryptedPacket(addrStr, protocol.PTSnapshot, payload); err != nil {
			fmt.Println("Failed to send encrypted snapshot to ", addrStr, "\n[ERROR]:\t", err)
		}
	}

	s.playerMu.RUnlock()
}
