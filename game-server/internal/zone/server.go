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

	"game-server/internal/persist"
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

	// gateway link
	w *bufio.Writer
}

type entity struct {
	EID shared.EntityID
	X, Y int16
	VX, VY int16
	HP uint16

	Owner shared.CharacterID
	Dirty bool
}

type player struct {
	SID shared.SessionID
	CID shared.CharacterID
	EID shared.EntityID

	nextClientTick uint32

	known map[shared.EntityID]struct{}
	lastSentPos map[shared.EntityID][2]int16
	lastSentHP map[shared.EntityID]uint16

	pendingEvents []string

	Transferring bool
}

func New(cfg Config) *Server {
	if cfg.TickHz <= 0 { cfg.TickHz = 20 }
	if cfg.AOIRadius <= 0 { cfg.AOIRadius = 25 }
	if cfg.CellSize <= 0 { cfg.CellSize = 8 }
	if cfg.BudgetBytes <= 0 { cfg.BudgetBytes = 900 }
	if cfg.StateEveryTicks <= 0 { cfg.StateEveryTicks = 5 }
	if cfg.SaveEveryTicks <= 0 { cfg.SaveEveryTicks = 20 }
	if cfg.TransferTargetZone == 0 {
		panic("zone: TransferTargetZone must be set (strict)")
	}
	if cfg.Store == nil || cfg.SaveQ == nil {
		panic("zone: Store and SaveQ must be provided (strict)")
	}

	return &Server{
		cfg: cfg,
		players: make(map[shared.SessionID]*player),
		ents: make(map[shared.EntityID]*entity),
		nextEID: 1,
		grid: spatial.New(cfg.CellSize),
	}
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil { return err }
	defer ln.Close()

	log.Printf("zone up: zone=%d listen=%s xferTarget=%d boundaryX=%d",
		s.cfg.ZoneID, s.cfg.ListenAddr, s.cfg.TransferTargetZone, s.cfg.TransferBoundaryX)

	c, err := ln.Accept()
	if err != nil { return err }
	defer c.Close()

	r := bufio.NewReaderSize(c, 64*1024)
	s.w = bufio.NewWriterSize(c, 64*1024)

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
			s.mu.Lock()
			s.enqueueDirtyLocked()
			s.mu.Unlock()
			return nil
		case fr, ok := <-inbound:
			if !ok { return errors.New("gateway link closed") }
			s.handleFrame(ctx, fr)
		case <-ticker.C:
			s.stepAndSend(ctx)
		}
	}
}

func (s *Server) allocEntity(cid shared.CharacterID, base persist.CharacterState) shared.EntityID {
	eid := s.nextEID
	s.nextEID++
	hp := base.HP
	if hp == 0 { hp = 100 }
	s.ents[eid] = &entity{
		EID: eid, X: base.X, Y: base.Y, HP: hp,
		Owner: cid, Dirty: true,
	}
	return eid
}

func (s *Server) handleFrame(ctx context.Context, fr wire.Frame) {
	switch fr.Type {
	case wire.MsgAttachPlayer:
		sid, cid, zid, err := wire.DecodeAttachPlayer(fr.Payload)
		if err != nil || uint32(zid) != s.cfg.ZoneID {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad attach"))
			return
		}
		s.attachFromStore(ctx, sid, cid)
		_ = wire.WriteFrame(s.w, wire.MsgAttachAck, nil)

	case wire.MsgAttachWithState:
		sid, cid, zid, x, y, hp, err := wire.DecodeAttachWithState(fr.Payload)
		if err != nil || uint32(zid) != s.cfg.ZoneID {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad attach-with-state"))
			return
		}
		s.mu.Lock()
		if _, ok := s.players[sid]; !ok {
			base := persist.CharacterState{CharacterID: cid, ZoneID: shared.ZoneID(s.cfg.ZoneID), X: x, Y: y, HP: hp}
			eid := s.allocEntity(cid, base)
			s.players[sid] = &player{
				SID: sid, CID: cid, EID: eid,
				known: make(map[shared.EntityID]struct{}),
				lastSentPos: make(map[shared.EntityID][2]int16),
				lastSentHP: make(map[shared.EntityID]uint16),
				pendingEvents: []string{"entered zone"},
			}
		}
		s.mu.Unlock()
		_ = wire.WriteFrame(s.w, wire.MsgAttachAck, nil)

	case wire.MsgDetachPlayer:
		sid, err := wire.DecodeDetachPlayer(fr.Payload)
		if err != nil {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad detach"))
			return
		}
		s.mu.Lock()
		p := s.players[sid]
		if p != nil {
			if e := s.ents[p.EID]; e != nil {
				e.Dirty = true
				s.enqueueCharacterLocked(e, p.CID)
				delete(s.ents, p.EID)
			}
			delete(s.players, sid)
		}
		s.mu.Unlock()

	case wire.MsgPlayerInput:
		sid, tick, mx, my, err := wire.DecodePlayerInput(fr.Payload)
		if err != nil {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad input"))
			return
		}
		s.mu.Lock()
		p := s.players[sid]
		if p == nil || p.Transferring {
			s.mu.Unlock()
			return
		}
		if tick < p.nextClientTick || tick > p.nextClientTick+64 {
			s.mu.Unlock()
			return
		}
		p.nextClientTick = tick + 1
		e := s.ents[p.EID]
		if e != nil {
			e.VX = mx; e.VY = my
			if e.HP > 0 { e.HP-- }
			e.Dirty = true
		}
		s.mu.Unlock()

	default:
		_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "unknown msg type"))
	}
}

func (s *Server) attachFromStore(ctx context.Context, sid shared.SessionID, cid shared.CharacterID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.players[sid]; ok { return }

	base, found, _ := s.cfg.Store.LoadCharacter(ctx, cid)
	if !found {
		base = persist.CharacterState{
			CharacterID: cid, ZoneID: shared.ZoneID(s.cfg.ZoneID),
			X: int16(int(cid%50)), Y: int16(int((cid*3)%50)), HP: 100,
		}
	} else {
		base.ZoneID = shared.ZoneID(s.cfg.ZoneID)
	}

	eid := s.allocEntity(cid, base)
	s.players[sid] = &player{
		SID: sid, CID: cid, EID: eid,
		known: make(map[shared.EntityID]struct{}),
		lastSentPos: make(map[shared.EntityID][2]int16),
		lastSentHP: make(map[shared.EntityID]uint16),
		pendingEvents: []string{"welcome"},
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
			e.Dirty = true
		}
	}
}

type eidDist struct {
	eid shared.EntityID
	d2  int32
}

func (s *Server) stepAndSend(ctx context.Context) {
	s.mu.Lock()
	s.stepPhysicsLocked()
	s.rebuildGridLocked()

	tick := s.serverTick
	sendState := (tick % uint32(s.cfg.StateEveryTicks)) == 0
	doSave := (tick % uint32(s.cfg.SaveEveryTicks)) == 0

	// transfer detection (toy X boundary)
	var transfers []struct{
		sid shared.SessionID
		cid shared.CharacterID
		st persist.CharacterState
	}
	for _, p := range s.players {
		if p.Transferring { continue }
		e := s.ents[p.EID]
		if e == nil { continue }
		if s.shouldTransfer(e.X) {
			p.Transferring = true
			st := persist.CharacterState{
				CharacterID: p.CID,
				ZoneID: shared.ZoneID(s.cfg.ZoneID),
				X: e.X, Y: e.Y, HP: e.HP,
				ServerTick: tick,
			}
			// enqueue save now (best effort)
			s.cfg.SaveQ.Enqueue(st)
			transfers = append(transfers, struct{
				sid shared.SessionID; cid shared.CharacterID; st persist.CharacterState
			}{sid: p.SID, cid: p.CID, st: st})
		}
	}

	// prepare per-session replication (minimal: move + state + event) within AOI
	type perOut struct {
		sid shared.SessionID
		ev []wire.RepEvent
		move []wire.RepEvent
		state []wire.RepEvent
	}
	out := make([]perOut, 0, len(s.players))
	tmp := make([]uint32, 0, 256)

	for _, p := range s.players {
		if p.Transferring { continue }
		pe := s.ents[p.EID]
		if pe == nil { continue }

		tmp = tmp[:0]
		cands := s.grid.QueryCircle(pe.X, pe.Y, s.cfg.AOIRadius, tmp)

		newSet := make(map[shared.EntityID]struct{}, len(cands))
		dists := make([]eidDist, 0, len(cands))
		for _, eidU := range cands {
			eid := shared.EntityID(eidU)
			newSet[eid] = struct{}{}
			ex, ey, ok := s.grid.GetPos(eidU)
			if !ok { continue }
			dx := int32(ex) - int32(pe.X); dy := int32(ey) - int32(pe.Y)
			dists = append(dists, eidDist{eid: eid, d2: dx*dx+dy*dy})
		}
		sort.Slice(dists, func(i,j int) bool {
			if dists[i].d2 == dists[j].d2 { return dists[i].eid < dists[j].eid }
			return dists[i].d2 < dists[j].d2
		})

		ev := make([]wire.RepEvent, 0, 8)
		if len(p.pendingEvents) > 0 {
			ev = append(ev, wire.RepEvent{Op: wire.RepEventText, Text: p.pendingEvents[0]})
			p.pendingEvents = p.pendingEvents[1:]
		}

		move := make([]wire.RepEvent, 0, 64)
		state := make([]wire.RepEvent, 0, 16)

		for eid := range p.known {
			if _, ok := newSet[eid]; !ok {
				move = append(move, wire.RepEvent{Op: wire.RepDespawn, EID: eid})
				delete(p.known, eid)
				delete(p.lastSentPos, eid)
				delete(p.lastSentHP, eid)
			}
		}
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
			if len(move) >= 256 { break }
		}

		if sendState {
			for _, ed := range dists {
				eid := ed.eid
				e := s.ents[eid]
				if e == nil { continue }
				if _, ok := p.known[eid]; !ok { continue }
				if p.lastSentHP[eid] != e.HP {
					state = append(state, wire.RepEvent{Op: wire.RepStateHP, EID: eid, Val: e.HP})
					p.lastSentHP[eid] = e.HP
				}
				if len(state) >= 64 { break }
			}
		}

		// budget (simple): ev->move->state
		b := s.cfg.BudgetBytes
		ev = trim(ev, b); b -= est(ev)
		move = trim(move, b); b -= est(move)
		state = trim(state, b)

		if len(ev)+len(move)+len(state) > 0 {
			out = append(out, perOut{sid: p.SID, ev: ev, move: move, state: state})
		}
	}

	if doSave {
		s.enqueueDirtyLocked()
	}
	s.mu.Unlock()

	// send replications
	for _, m := range out {
		if len(m.ev) > 0 { _ = wire.WriteFrame(s.w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanEvent, m.ev)) }
		if len(m.move) > 0 { _ = wire.WriteFrame(s.w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanMove, m.move)) }
		if len(m.state) > 0 { _ = wire.WriteFrame(s.w, wire.MsgReplicate, wire.EncodeReplicate(m.sid, tick, wire.ChanState, m.state)) }
	}

	// send transfers (after tick outputs)
	for _, t := range transfers {
		payload := wire.EncodeTransfer(t.sid, t.cid, shared.ZoneID(s.cfg.TransferTargetZone), t.st)
		_ = wire.WriteFrame(s.w, wire.MsgTransfer, payload)

		// NOTE (strict): this skeleton immediately removes the entity/player after requesting transfer.
		// In production: require ACK from gateway + target zone before finalizing.
		s.mu.Lock()
		p := s.players[t.sid]
		if p != nil {
			delete(s.ents, p.EID)
			delete(s.players, t.sid)
		}
		s.mu.Unlock()
	}

	_ = ctx
}

func (s *Server) shouldTransfer(x int16) bool {
	b := s.cfg.TransferBoundaryX
	if b > 0 { return x > b }
	if b < 0 { return x < b }
	return false
}

// persistence helpers (call with lock)
func (s *Server) enqueueCharacterLocked(e *entity, cid shared.CharacterID) {
	st := persist.CharacterState{
		CharacterID: cid,
		ZoneID: shared.ZoneID(s.cfg.ZoneID),
		X: e.X, Y: e.Y, HP: e.HP,
		ServerTick: s.serverTick,
	}
	s.cfg.SaveQ.Enqueue(st)
	e.Dirty = false
}
func (s *Server) enqueueDirtyLocked() {
	for _, p := range s.players {
		e := s.ents[p.EID]
		if e != nil && e.Dirty {
			s.enqueueCharacterLocked(e, p.CID)
		}
	}
}

// budget helpers (coarse; same for all channels here)
func est(evs []wire.RepEvent) int {
	sz := 23
	for _, e := range evs {
		switch e.Op {
		case wire.RepSpawn, wire.RepMove:
			sz += 9
		case wire.RepDespawn:
			sz += 5
		case wire.RepStateHP:
			sz += 7
		case wire.RepEventText:
			l := len(e.Text); if l > 65535 { l = 65535 }
			sz += 3 + l
		}
	}
	return sz
}
func trim(evs []wire.RepEvent, budget int) []wire.RepEvent {
	if budget <= 23 { return evs[:0] }
	out := evs[:0]
	for _, e := range evs {
		out = append(out, e)
		if est(out) > budget {
			out = out[:len(out)-1]
			break
		}
	}
	return out
}
