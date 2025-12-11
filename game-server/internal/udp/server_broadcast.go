package udp

import (
	"fmt"
	ecs2 "game-server/internal/ecs"
	"game-server/pkg/protocol"
)

func (s *Server) broadcastSnapshot() {
	s.playerMu.RLock()

	snapshots := make([]protocol.EntitySnapshot, 0, len(s.playerEntities))

	for _, ent := range s.playerEntities {
		pos := ecs2.Get(s.world, ent, ecs2.Position)

		snapshots = append(snapshots, protocol.EntitySnapshot{
			ID:        uint32(ent),
			PositionX: pos.X,
			PositionY: pos.Y,
		})
	}

	payload := protocol.MarshalSnapshotPayload(snapshots)

	for addrStr := range s.playerConnections {
		if err := s.sendEncryptedPacket(addrStr, protocol.PTSnapshot, payload); err != nil {
			fmt.Println("Failed to send encrypted snapshot to ", addrStr, "\n[ERROR]:\t", err)
		}
	}

	s.playerMu.RUnlock()
}
