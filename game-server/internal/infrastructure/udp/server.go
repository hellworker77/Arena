package udp

import (
	"game-server/internal/infrastructure/udp/udp_types"
	"net"
	"sync"
	"time"
)

const (
	tickRate = 20
) // 20 ticks per second

type Server struct {
	conn   *net.UDPConn
	Inputs chan udp_types.InputPacket

	playerMu sync.RWMutex
	players  map[string]*net.UDPAddr
}

func NewServer(addr string) (*Server, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:    conn,
		Inputs:  make(chan udp_types.InputPacket, 1024),
		players: make(map[string]*net.UDPAddr),
	}, nil
}

func (s *Server) Listen() {
	buf := make([]byte, 1024)

	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		s.playerMu.Lock()
		if _, ok := s.players[addr.String()]; !ok {
			s.players[addr.String()] = addr
		}

		s.playerMu.Unlock()

		s.Inputs <- udp_types.InputPacket{From: addr, Data: buf[:n]}
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
			s.handleInput(packet)
		default:
			return
		}
	}
}

func (s *Server) handleInput(packet udp_types.InputPacket) {
	_, _ = s.conn.WriteToUDP(packet.Data, packet.From)
}
