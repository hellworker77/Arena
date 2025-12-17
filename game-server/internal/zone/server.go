package zone

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"game-server/internal/shared"
	"game-server/internal/shared/wire"
)

type Server struct {
	cfg Config

	mu      sync.Mutex
	players map[shared.SessionID]*player

	serverTick uint32
}

type player struct {
	SID shared.SessionID
	CID shared.CharacterID

	X, Y   int16
	VX, VY int16

	nextClientTick uint32
}

func New(cfg Config) *Server {
	if cfg.TickHz <= 0 {
		cfg.TickHz = 20
	}
	return &Server{
		cfg:     cfg,
		players: make(map[shared.SessionID]*player),
	}
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Printf("zone up: zone=%d listen=%s", s.cfg.ZoneID, s.cfg.ListenAddr)

	c, err := ln.Accept()
	if err != nil {
		return err
	}
	defer c.Close()

	r := bufio.NewReaderSize(c, 64*1024)
	w := bufio.NewWriterSize(c, 64*1024)

	inbound := make(chan wire.Frame, 256)
	go func() {
		defer close(inbound)
		for {
			fr, err := wire.ReadFrame(r)
			if err != nil {
				return
			}
			inbound <- fr
		}
	}()

	ticker := time.NewTicker(time.Second / time.Duration(s.cfg.TickHz))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case fr, ok := <-inbound:
			if !ok {
				return errors.New("gateway link closed")
			}
			s.handleFrame(w, fr)

		case <-ticker.C:
			s.step(w)
		}
	}
}

func (s *Server) handleFrame(w *bufio.Writer, fr wire.Frame) {
	switch fr.Type {
	case wire.MsgAttachPlayer:
		sid, cid, zid, err := wire.DecodeAttachPlayer(fr.Payload)
		if err != nil || uint32(zid) != s.cfg.ZoneID {
			_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad attach"))
			return
		}
		s.mu.Lock()
		if _, ok := s.players[sid]; !ok {
			s.players[sid] = &player{SID: sid, CID: cid}
		}
		s.mu.Unlock()
		_ = wire.WriteFrame(w, wire.MsgAttachAck, nil)

	case wire.MsgDetachPlayer:
		sid, err := wire.DecodeDetachPlayer(fr.Payload)
		if err != nil {
			_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad detach"))
			return
		}
		s.mu.Lock()
		delete(s.players, sid)
		s.mu.Unlock()

	case wire.MsgPlayerInput:
		sid, tick, mx, my, err := wire.DecodePlayerInput(fr.Payload)
		if err != nil {
			_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad input"))
			return
		}
		s.mu.Lock()
		p := s.players[sid]
		if p == nil {
			s.mu.Unlock()
			_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrNoPlayer, "no player"))
			return
		}

		// Strict: accept only current tick, reject too old / too far ahead.
		if tick < p.nextClientTick {
			s.mu.Unlock()
			return
		}
		if tick > p.nextClientTick+64 {
			s.mu.Unlock()
			return
		}
		p.VX = mx
		p.VY = my
		p.nextClientTick = tick + 1
		s.mu.Unlock()

	default:
		_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "unknown msg type"))
	}
}

func (s *Server) step(w *bufio.Writer) {
	s.serverTick++

	s.mu.Lock()
	positions := make([][2]int16, 0, len(s.players))
	for _, p := range s.players {
		p.X += p.VX
		p.Y += p.VY
		positions = append(positions, [2]int16{p.X, p.Y})
	}
	s.mu.Unlock()

	_ = wire.WriteFrame(w, wire.MsgSnapshot, wire.EncodeSnapshot(s.serverTick, positions))
}
