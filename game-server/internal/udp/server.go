package udp

import (
	"context"
	"crypto/rand"
	"game-server/internal/game"
	udp_types2 "game-server/internal/udp/udp_types"
	"game-server/pkg/auth"
	"net"
	"sync"
	"time"
)

const tickRate = 20
const maxAuthAttempts = 5
const authAttemptWindow = 10 * time.Second

type Server struct {
	conn *net.UDPConn

	ctx    context.Context
	cancel context.CancelFunc

	playerMu          sync.RWMutex
	playerConnections map[string]*net.UDPAddr
	playerSessions    map[string]*auth.PlayerSession
	state             map[string]*udp_types2.ClientState
	playerClaims      map[string]*auth.JWTClaims

	authAttemptsMu  sync.Mutex
	authAttempts    map[string][]time.Time
	authLimitN      int
	authLimitWindow time.Duration

	inputLimiters   map[string]*tokenBucket
	inputLimitBurst int
	inputLimitRate  int
	jwtValidator    *auth.JWTValidator

	// Handshake cookie (stateless challenge) for pre-auth flood resistance.
	cookieSecret     [32]byte
	cookieBucketSec  uint32
	unauthLimiters   map[string]*tokenBucket
	unauthLimitBurst int
	unauthLimitRate  int
	allowLegacyAuth  bool

	clientIdleTimeout time.Duration
	cleanupInterval   time.Duration
	unauthEntryTTL    time.Duration

	loop *game.Loop

	Inputs chan udp_types2.InputPacket
}

// ConfigureReplication tunes interest management and delta snapshot cadence.
// Call before Startup().
func (s *Server) ConfigureReplication(interestRadius float32, fullEveryTicks uint32) {
	if s.loop != nil {
		s.loop.ConfigureReplication(interestRadius, fullEveryTicks)
	}
}

// ConfigureReplicationAdvanced additionally sets spatial grid size and snapshot payload budget.
// Call before Startup().
func (s *Server) ConfigureReplicationAdvanced(interestRadius float32, fullEveryTicks uint32, gridCellSize float32, maxSnapshotBytes int) {
	if s.loop == nil {
		return
	}
	s.loop.ConfigureReplication(interestRadius, fullEveryTicks)
	if gridCellSize > 0 {
		s.loop.ConfigureSpatialGrid(gridCellSize)
	}
	if maxSnapshotBytes != 0 {
		s.loop.ConfigureSnapshotBudget(maxSnapshotBytes)
	}
}

func NewServer(addr string, jwtCfg auth.JwtCfg, allowLegacyAuth bool, clientIdleTimeout, cleanupInterval, unauthEntryTTL time.Duration) (*Server, error) {
	if clientIdleTimeout <= 0 {
		clientIdleTimeout = 30 * time.Second
	}
	if cleanupInterval <= 0 {
		cleanupInterval = 5 * time.Second
	}
	if unauthEntryTTL <= 0 {
		unauthEntryTTL = 60 * time.Second
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	jwks, _ := auth.NewJwksCache(jwtCfg.Authority, jwtCfg.RotationIntervalSec, jwtCfg.AutoRefreshIntervalSec)
	jwtValidator := &auth.JWTValidator{
		Cache:     jwks,
		Authority: jwtCfg.Authority,
		Audience:  jwtCfg.Audience,
	}

	ctx, cancel := context.WithCancel(context.Background())
	loop := game.NewLoop(game.NewEngine(), tickRate)

	var secret [32]byte
	_, _ = rand.Read(secret[:]) // best-effort; if it fails, secret stays weak but server still runs

	return &Server{
		conn:              conn,
		ctx:               ctx,
		cancel:            cancel,
		Inputs:            make(chan udp_types2.InputPacket, 1024),
		playerConnections: make(map[string]*net.UDPAddr),
		playerSessions:    make(map[string]*auth.PlayerSession),
		playerClaims:      make(map[string]*auth.JWTClaims),
		state:             make(map[string]*udp_types2.ClientState),
		inputLimiters:     make(map[string]*tokenBucket),
		inputLimitBurst:   30,
		inputLimitRate:    60,
		authAttempts:      make(map[string][]time.Time),
		authLimitN:        maxAuthAttempts,
		authLimitWindow:   authAttemptWindow,
		jwtValidator:      jwtValidator,
		loop:              loop,
		cookieSecret:      secret,
		cookieBucketSec:   10,
		unauthLimiters:    make(map[string]*tokenBucket),
		unauthLimitBurst:  20,
		unauthLimitRate:   20,
		allowLegacyAuth:   allowLegacyAuth,
		clientIdleTimeout: clientIdleTimeout,
		cleanupInterval:   cleanupInterval,
		unauthEntryTTL:    unauthEntryTTL,
	}, nil
}

func (s *Server) Listen() {
	buf := make([]byte, 2048)

	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		s.Inputs <- udp_types2.InputPacket{From: addr, Data: buf[:n]}
	}
}

func (s *Server) Close() {
	s.cancel()
	s.conn.Close()
}

func (s *Server) Startup() {
	// Packet processing is decoupled from the simulation.
	go s.processPackets()

	// Simulation loop owns the authoritative game state.
	s.loop.Run(s.ctx, func(frame game.SnapshotFrame) {
		s.broadcastSnapshotFrame(frame)
	})
}

func (s *Server) processPackets() {
	cleanupTick := time.NewTicker(s.cleanupInterval)
	defer cleanupTick.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-cleanupTick.C:
			s.cleanupOnce()
		case packet := <-s.Inputs:
			s.handleEncryptedPacket(packet)
		}
	}
}

func (s *Server) cleanupOnce() {
	now := time.Now()

	// 1) Disconnect idle authenticated clients.
	var toDrop []string
	s.playerMu.RLock()
	for addrStr, st := range s.state {
		// Only consider authenticated clients.
		if s.playerSessions[addrStr] == nil || st == nil {
			continue
		}
		if now.Sub(st.LastHeard) > s.clientIdleTimeout {
			toDrop = append(toDrop, addrStr)
		}
	}
	s.playerMu.RUnlock()
	for _, addrStr := range toDrop {
		s.disconnect(addrStr)
	}

	// 2) Prune pre-auth limiter state to avoid unbounded maps.
	for addrStr, b := range s.unauthLimiters {
		if b == nil {
			delete(s.unauthLimiters, addrStr)
			continue
		}
		if now.Sub(b.last) > s.unauthEntryTTL {
			delete(s.unauthLimiters, addrStr)
		}
	}

	// 3) Prune auth attempt tracking.
	s.authAttemptsMu.Lock()
	for addrStr, arr := range s.authAttempts {
		// Keep only recent attempts (authLimitWindow); drop map entry if empty.
		out := arr[:0]
		for _, t := range arr {
			if now.Sub(t) <= s.authLimitWindow {
				out = append(out, t)
			}
		}
		if len(out) == 0 {
			delete(s.authAttempts, addrStr)
		} else {
			s.authAttempts[addrStr] = out
		}
	}
	s.authAttemptsMu.Unlock()
}

func (s *Server) disconnect(addrStr string) {
	// Remove server-side network/session state.
	s.playerMu.Lock()
	delete(s.playerConnections, addrStr)
	delete(s.playerSessions, addrStr)
	delete(s.playerClaims, addrStr)
	delete(s.state, addrStr)
	delete(s.inputLimiters, addrStr)
	s.playerMu.Unlock()

	// Remove from simulation.
	s.loop.RemovePlayer(addrStr)
}
