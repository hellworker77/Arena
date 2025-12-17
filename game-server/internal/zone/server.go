package zone

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"game-server/internal/shared"
	"game-server/internal/shared/wire"
	"game-server/internal/zone/spatial"
)

type Server struct {
	cfg Config

	mu      sync.Mutex
	players map[shared.SessionID]*player
	ents    map[shared.EntityID]*entity

	nextEID     shared.EntityID
	serverTick  uint32
	grid        *spatial.Grid
}

type entity struct {
	EID shared.EntityID
	X, Y int16
	VX, VY int16

	HP uint16 // toy state
}

type player struct {
	SID shared.SessionID
	CID shared.CharacterID
	EID shared.EntityID

	nextClientTick uint32

	known       map[shared.EntityID]struct{}
	lastSentPos map[shared.EntityID][2]int16
	lastSentHP  map[shared.EntityID]uint16

	// event queue (toy)
	pendingEvents []string
}

func New(cfg Config) *Server {
	if cfg.TickHz <= 0 { cfg.TickHz = 20 }
	if cfg.AOIRadius <= 0 { cfg.AOIRadius = 25 }
	if cfg.CellSize <= 0 { cfg.CellSize = 8 }
	if cfg.MaxMoveEvents <= 0 { cfg.MaxMoveEvents = 256 }
	if cfg.MaxStateEvents <= 0 { cfg.MaxStateEvents = 64 }
	if cfg.MaxEventEvents <= 0 { cfg.MaxEventEvents = 64 }
	if cfg.BudgetBytes <= 0 { cfg.BudgetBytes = 900 }
	if cfg.StateEveryTicks <= 0 { cfg.StateEveryTicks = 5 }

	return &Server{
		cfg:     cfg,
		players: make(map[shared.SessionID]*player),
		ents:    make(map[shared.EntityID]*entity),
		nextEID:  1,
		grid:     spatial.New(cfg.CellSize),
	}
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil { return err }
	defer ln.Close()

	log.Printf("zone up: zone=%d listen=%s aoi=%d cell=%d budget=%d stateEvery=%d",
		s.cfg.ZoneID, s.cfg.ListenAddr, s.cfg.AOIRadius, s.cfg.CellSize, s.cfg.BudgetBytes, s.cfg.StateEveryTicks)

	c, err := ln.Accept()
	if err != nil { return err }
	defer c.Close()

	r := bufio.NewReaderSize(c, 64*1024)
	w := bufio.NewWriterSize(c, 64*1024)

	inbound := make(chan wire.Frame, 512)
	go func() {
		defer close(inbound)
		for {
			fr, err := wire.ReadFrame(r)
			if err != nil { return }
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
			if !ok { return errors.New("gateway link closed") }
			s.handleFrame(w, fr)
		case <-ticker.C:
			s.stepAndReplicate(w)
		}
	}
}

func (s *Server) allocEntity(x, y int16) shared.EntityID {
	eid := s.nextEID
	s.nextEID++
	s.ents[eid] = &entity{EID: eid, X: x, Y: y, HP: 100}
	return eid
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
			eid := s.allocEntity(int16(int(s.nextEID%50)), int16(int((s.nextEID*3)%50)))
			s.players[sid] = &player{
				SID: sid, CID: cid, EID: eid,
				known: make(map[shared.EntityID]struct{}),
				lastSentPos: make(map[shared.EntityID][2]int16),
				lastSentHP: make(map[shared.EntityID]uint16),
				pendingEvents: []string{"welcome"},
			}
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
		p := s.players[sid]
		if p != nil { delete(s.ents, p.EID) }
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
		if tick < p.nextClientTick || tick > p.nextClientTick+64 {
			s.mu.Unlock()
			return
		}
		p.nextClientTick = tick + 1
		e := s.ents[p.EID]
		if e != nil {
			e.VX = mx
			e.VY = my
			// toy: taking input drains HP a bit (to show state changes)
			if e.HP > 0 { e.HP-- }
		}
		s.mu.Unlock()
	default:
		_ = wire.WriteFrame(w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "unknown msg type"))
	}
}

func (s *Server) rebuildGridLocked() {
	s.grid.Clear()
	for eid, e := range s.ents {
		s.grid.Insert(uint32(eid), e.X, e.Y)
	}
}

func (s *Server) stepPhysicsLocked() {
	s.serverTick++
	for _, e := range s.ents {
		if e.VX != 0 || e.VY != 0 {
			e.X += e.VX
			e.Y += e.VY
		}
	}
}

type eidDist struct {
	eid shared.EntityID
	d2  int32
}

func (s *Server) stepAndReplicate(w *bufio.Writer) {
	s.mu.Lock()
	s.stepPhysicsLocked()
	s.rebuildGridLocked()
	tick := s.serverTick
	sendState := (tick % uint32(s.cfg.StateEveryTicks)) == 0

	type perOut struct {
		sid shared.SessionID
		move []wire.RepEvent
		state []wire.RepEvent
		ev []wire.RepEvent
	}
	out := make([]perOut, 0, len(s.players))

	tmp := make([]uint32, 0, 256)
	for _, p := range s.players {
		pe := s.ents[p.EID]
		if pe == nil { continue }

		// Build candidate list within AOI
		tmp = tmp[:0]
		cands := s.grid.QueryCircle(pe.X, pe.Y, s.cfg.AOIRadius, tmp)

		newSet := make(map[shared.EntityID]struct{}, len(cands))
		dists := make([]eidDist, 0, len(cands))
		for _, eidU := range cands {
			eid := shared.EntityID(eidU)
			newSet[eid] = struct{}{}
			ex, ey, ok := s.grid.GetPos(eidU)
			if !ok { continue }
			dx := int32(ex) - int32(pe.X)
			dy := int32(ey) - int32(pe.Y)
			dists = append(dists, eidDist{eid: eid, d2: dx*dx + dy*dy})
		}
		sort.Slice(dists, func(i, j int) bool {
			if dists[i].d2 == dists[j].d2 { return dists[i].eid < dists[j].eid }
			return dists[i].d2 < dists[j].d2
		})

		// Queues (priority already by construction)
		move := make([]wire.RepEvent, 0, 64)
		state := make([]wire.RepEvent, 0, 16)
		ev := make([]wire.RepEvent, 0, 8)

		// 1) event channel: flush queued events first (within cap)
		if len(p.pendingEvents) > 0 {
			limit := s.cfg.MaxEventEvents
			if limit > len(p.pendingEvents) { limit = len(p.pendingEvents) }
			for i := 0; i < limit; i++ {
				ev = append(ev, wire.RepEvent{Op: wire.RepEventText, Text: p.pendingEvents[i]})
			}
			p.pendingEvents = p.pendingEvents[limit:]
		}

		// 2) despawn (highest move priority)
		for eid := range p.known {
			if _, ok := newSet[eid]; !ok {
				move = append(move, wire.RepEvent{Op: wire.RepDespawn, EID: eid})
				delete(p.known, eid)
				delete(p.lastSentPos, eid)
				delete(p.lastSentHP, eid)
				if len(move) >= s.cfg.MaxMoveEvents { break }
			}
		}

		// 3) spawn/move (nearest first)
		if len(move) < s.cfg.MaxMoveEvents {
			for _, ed := range dists {
				eid := ed.eid
				e := s.ents[eid]
				if e == nil { continue }
				if _, ok := p.known[eid]; !ok {
					move = append(move, wire.RepEvent{Op: wire.RepSpawn, EID: eid, X: e.X, Y: e.Y})
					p.known[eid] = struct{}{}
					p.lastSentPos[eid] = [2]int16{e.X, e.Y}
					p.lastSentHP[eid] = e.HP
				} else {
					prev := p.lastSentPos[eid]
					if prev[0] != e.X || prev[1] != e.Y {
						move = append(move, wire.RepEvent{Op: wire.RepMove, EID: eid, X: e.X, Y: e.Y})
						p.lastSentPos[eid] = [2]int16{e.X, e.Y}
					}
				}
				if len(move) >= s.cfg.MaxMoveEvents { break }
			}
		}

		// 4) state channel (less frequent)
		if sendState && len(state) < s.cfg.MaxStateEvents {
			for _, ed := range dists {
				eid := ed.eid
				e := s.ents[eid]
				if e == nil { continue }
				// only for known entities (don’t leak state without spawn)
				if _, ok := p.known[eid]; !ok { continue }
				prev := p.lastSentHP[eid]
				if prev != e.HP {
					state = append(state, wire.RepEvent{Op: wire.RepStateHP, EID: eid, Val: e.HP})
					p.lastSentHP[eid] = e.HP
				}
				if len(state) >= s.cfg.MaxStateEvents { break }
			}
		}

		// Per-session scheduler: strict byte budget per tick.
		// Priority: ev -> move -> state.
		budget := s.cfg.BudgetBytes
		ev = trimToBudgetEvents(ev, budget, wire.ChanEvent)
		budget -= estimateEncodedSize(wire.ChanEvent, ev)

		move = trimToBudgetEvents(move, budget, wire.ChanMove)
		budget -= estimateEncodedSize(wire.ChanMove, move)

		state = trimToBudgetEvents(state, budget, wire.ChanState)

		if len(ev) > 0 || len(move) > 0 || len(state) > 0 {
			out = append(out, perOut{sid: p.SID, move: move, state: state, ev: ev})
		}
	}
	s.mu.Unlock()

	for _, m := range out {
		if len(m.ev) > 0 {
			_ = wire.WriteFrame(w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanEvent, m.ev))
		}
		if len(m.move) > 0 {
			_ = wire.WriteFrame(w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanMove, m.move))
		}
		if len(m.state) > 0 {
			_ = wire.WriteFrame(w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanState, m.state))
		}
	}
}

// Size estimation (strict upper-ish bound) to enforce budget before encoding.
// payload overhead per replicate: 16+4+1+2 = 23 bytes
func estimateEncodedSize(ch wire.RepChannel, evs []wire.RepEvent) int {
	_ = ch
	sz := 23
	for _, e := range evs {
		switch e.Op {
		case wire.RepSpawn, wire.RepMove:
			sz += 1 + 4 + 4
		case wire.RepDespawn:
			sz += 1 + 4
		case wire.RepStateHP:
			sz += 1 + 4 + 2
		case wire.RepEventText:
			if len(e.Text) > 65535 {
				sz += 1 + 2 + 65535
			} else {
				sz += 1 + 2 + len(e.Text)
			}
		}
	}
	return sz
}

func trimToBudgetEvents(evs []wire.RepEvent, budget int, ch wire.RepChannel) []wire.RepEvent {
	if budget <= 23 { // cannot even fit header
		return evs[:0]
	}
	// greedy: keep prefix until size fits (events are already priority-ordered)
	out := evs[:0]
	for _, e := range evs {
		out = append(out, e)
		if estimateEncodedSize(ch, out) > budget {
			out = out[:len(out)-1]
			break
		}
	}
	return out
}
