package udp

import (
	"fmt"
	"game-server/internal/game"
	"game-server/pkg/protocol"
)

func (s *Server) broadcastSnapshotFrame(frame game.SnapshotFrame) {
	s.playerMu.RLock()
	for addrStr, payload := range frame.Payloads {
		// Only send to authenticated connections we still track.
		if s.playerConnections[addrStr] == nil || s.playerSessions[addrStr] == nil {
			continue
		}
		if err := s.sendEncryptedPacket(addrStr, protocol.PTSnapshot, payload); err != nil {
			fmt.Println("Failed to send encrypted snapshot to ", addrStr, "\n[ERROR]:\t", err)
		}
	}
	s.playerMu.RUnlock()
}
