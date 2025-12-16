package udp

import (
	"crypto/subtle"
	"fmt"
	"game-server/internal/udp/udp_types"
	"game-server/pkg/auth"
	"game-server/pkg/protocol"
	"net"
	"time"
)

func (s *Server) handleHello(addr *net.UDPAddr, data []byte) {
	// Stateless cookie handshake:
	// client -> PTHello(nonce)
	// server -> PTChallenge(bucket,cookie)
	p, err := protocol.UnmarshalHello(data)
	if err != nil {
		return
	}
	bucket := s.timeBucketNow()
	cookie := s.mintCookie(addr.String(), p.ClientNonce, bucket)
	resp := protocol.MarshalChallenge(protocol.ChallengePayload{TimeBucket: bucket, Cookie: cookie})
	s.sendPlainPacket(addr, protocol.PTChallenge, resp)
}

func (s *Server) handleAuth(addr *net.UDPAddr, data []byte) {
	addrStr := addr.String()

	// Prefer the cookie-based handshake. It's stateless (no server allocations) and blocks spoofed-source floods.
	var token string
	if p, err := protocol.UnmarshalAuthWithCookie(data); err == nil {
		// Accept current bucket and neighbouring buckets to allow small clock skew / packet reordering.
		valid := false
		buckets := []uint32{p.TimeBucket}
		if p.TimeBucket > 0 {
			buckets = append(buckets, p.TimeBucket-1)
		}
		buckets = append(buckets, p.TimeBucket+1)
		for _, b := range buckets {
			expected := s.mintCookie(addrStr, p.ClientNonce, b)
			if subtle.ConstantTimeCompare(expected[:], p.Cookie[:]) == 1 {
				valid = true
				break
			}
		}
		if !valid {
			s.sendPlainPacket(addr, protocol.PTAuthResp, []byte("BAD_COOKIE"))
			return
		}
		token = p.Token
	} else {
		// Legacy mode: token-only PTAuth (no cookie). Keep off in production.
		if !s.allowLegacyAuth {
			s.sendPlainPacket(addr, protocol.PTAuthResp, []byte("NEED_HELLO"))
			return
		}
		token = string(data)
	}

	if !s.allowAuthAttempt(addrStr) {
		s.sendPlainPacket(addr, protocol.PTAuthResp, []byte("TOO_MANY_ATTEMPTS"))
		return
	}

	claims, err := s.jwtValidator.ValidateToken(token)
	if err != nil {
		s.sendPlainPacket(addr, protocol.PTAuthResp, []byte("INVALID_TOKEN"))
		return
	}

	session, err := auth.NewSessionFromToken(token, addrStr)
	if err != nil {
		s.sendPlainPacket(addr, protocol.PTAuthResp, []byte("SESSION_CREATION_FAILED"))
		return
	}

	s.playerMu.Lock()
	defer s.playerMu.Unlock()

	s.playerConnections[addrStr] = addr
	s.playerSessions[addrStr] = session
	s.state[addrStr] = udp_types.NewClientState()
	s.playerClaims[addrStr] = claims

	entID := s.loop.AddPlayer(addrStr)

	meta := fmt.Sprintf("OK|ent:%d", entID) // simple example
	// After auth, both sides can derive the session keys from the token, so the response is encrypted.
	_ = s.sendEncryptedPacket(addrStr, protocol.PTAuthResp, []byte(meta))
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
