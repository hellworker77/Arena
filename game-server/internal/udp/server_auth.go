package udp

import (
	"fmt"
	"game-server/internal/udp/udp_types"
	"game-server/pkg/auth"
	"game-server/pkg/protocol"
	"net"
	"time"
)

func (s *Server) handleAuth(addr *net.UDPAddr, data []byte) {
	token := string(data)
	addrStr := addr.String()

	if !s.allowAuthAttempt(addrStr) {
		s.sendEncryptedPacket(addrStr, protocol.PTAuthResp, []byte("TOO_MANY_ATTEMPTS"))
		return
	}

	claims, err := s.jwtValidator.ValidateToken(token)
	if err != nil {
		s.sendEncryptedPacket(addrStr, protocol.PTAuthResp, []byte("INVALID_TOKEN"))
		return
	}

	session, err := auth.NewSessionFromToken(token, addrStr)
	if err != nil {
		s.sendEncryptedPacket(addrStr, protocol.PTAuthResp, []byte("SESSION_CREATION_FAILED"))
		return
	}

	s.playerMu.Lock()
	defer s.playerMu.Unlock()

	s.playerConnections[addrStr] = addr
	s.playerSessions[addrStr] = session
	s.state[addrStr] = udp_types.NewClientState()
	s.playerClaims[addrStr] = claims

	e := s.createPlayerEntity(addrStr)
	s.playerEntities[addrStr] = e

	meta := fmt.Sprintf("OK|ent:%d", e) // simple example
	if err := s.sendEncryptedPacket(addrStr, protocol.PTAuthResp, []byte(meta)); err != nil {
		fmt.Println("Failed to send PTAuthOk to", addrStr, ":", err)
	}
}

func (s *Server) allowAuthAttempt(addrStr string) bool {
	s.authAttemptsMu.Lock()
	defer s.authAttemptsMu.Unlock()

	now := time.Now()
	windowStart := now.Add(-s.authLimitWindow)

	arr := s.authAttempts[addrStr]

	i := 0

	for ; i < len(arr); i++ {
		if arr[i].After(windowStart) {
			break
		}
	}

	arr = arr[i:]
	arr = append(arr, now)
	s.authAttempts[addrStr] = arr

	return len(arr) <= s.authLimitN
}
