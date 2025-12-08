package udp

import (
	"game-server/internal/ecs"
	"game-server/internal/ecs/ecs_signatures/static"
	ecs_systems2 "game-server/internal/ecs/ecs_systems"
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

	playerMu          sync.RWMutex
	playerConnections map[string]*net.UDPAddr
	playerEntities    map[string]static.EntityID
	playerSessions    map[string]*auth.PlayerSession
	state             map[string]*udp_types2.ClientState
	playerClaims      map[string]*auth.JWTClaims

	authAttemptsMu  sync.Mutex
	authAttempts    map[string][]time.Time
	authLimitN      int
	authLimitWindow time.Duration

	jwtValidator *auth.JWTValidator

	world   *ecs.World
	systems []ecs_systems2.ECSSystem

	Inputs chan udp_types2.InputPacket
}

func NewServer(addr string, jwtCfg auth.JwtCfg) (*Server, error) {
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

	return &Server{
		conn:              conn,
		Inputs:            make(chan udp_types2.InputPacket, 1024),
		playerConnections: make(map[string]*net.UDPAddr),
		playerEntities:    make(map[string]static.EntityID),
		playerSessions:    make(map[string]*auth.PlayerSession),
		playerClaims:      make(map[string]*auth.JWTClaims),
		state:             make(map[string]*udp_types2.ClientState),
		authAttempts:      make(map[string][]time.Time),
		authLimitN:        maxAuthAttempts,
		authLimitWindow:   authAttemptWindow,
		world:             ecs.NewWorld(),
		jwtValidator:      jwtValidator,
		systems: []ecs_systems2.ECSSystem{
			ecs_systems2.MovementSystem{},
		},
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
	s.conn.Close()
}

func (s *Server) Startup() {
	tickDuration := time.Duration(float64(time.Second) / float64(tickRate))
	ticker := time.NewTicker(tickDuration)

	defer ticker.Stop()

	for range ticker.C {
		s.update()
	}
}

func (s *Server) update() {
	for {
		select {
		case packet := <-s.Inputs:
			s.handleEncryptedPacket(packet)
		default:
			goto SYSTEMS
		}
	}

SYSTEMS:
	dt := float32(1.0 / tickRate)

	for _, sys := range s.systems {
		sys.Run(s.world, dt)
	}
	s.broadcastSnapshot()
}
