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

	"game-server/internal/metrics"
	"game-server/internal/persist"
	"game-server/internal/shared"
	"game-server/internal/shared/wire"
	"game-server/internal/zone/spatial"
)

type Server struct {
	cfg Config

	mu sync.Mutex

	w *bufio.Writer // gateway link

	world *World
	grid *spatial.Grid
	serverTick uint32

	players map[shared.SessionID]*player

	// pending transfer prepare waiting for commit/abort (Step13)
	transferPending map[shared.SessionID]*pendingTransfer

	met *metrics.Counters
}

type player struct {
	SID shared.SessionID
	CID shared.CharacterID
	EID shared.EntityID

	Interest wire.InterestMask

	nextClientTick uint32

	known map[shared.EntityID]struct{}
	lastSentPos map[shared.EntityID][2]int16
	lastSentHP map[shared.EntityID]uint16

	pendingEvents []string
}

type posSample struct {
	Tick uint32
	X, Y int16
}

type posHistory struct {
	cap int
	s   []posSample
}

func newPosHistory(capacity int) *posHistory {
	if capacity <= 0 { capacity = 40 }
	return &posHistory{cap: capacity, s: make([]posSample, 0, capacity)}
}

func (h *posHistory) add(tick uint32, x, y int16) {
	// keep monotonic by tick
	if n := len(h.s); n > 0 && h.s[n-1].Tick == tick {
		h.s[n-1] = posSample{Tick: tick, X: x, Y: y}
		return
	}
	h.s = append(h.s, posSample{Tick: tick, X: x, Y: y})
	if len(h.s) > h.cap {
		h.s = h.s[len(h.s)-h.cap:]
	}
}

// sampleAt returns position at or before tick (nearest older). If no data, ok=false.
func (h *posHistory) sampleAt(tick uint32) (x, y int16, ok bool) {
	if len(h.s) == 0 { return 0,0,false }
	// if before earliest
	if tick <= h.s[0].Tick {
		return h.s[0].X, h.s[0].Y, true
	}
	// if after latest
	last := h.s[len(h.s)-1]
	if tick >= last.Tick {
		return last.X, last.Y, true
	}
	// binary search for rightmost <= tick
	lo, hi := 0, len(h.s)-1
	for lo <= hi {
		m := (lo+hi)/2
		if h.s[m].Tick == tick {
			return h.s[m].X, h.s[m].Y, true
		}
		if h.s[m].Tick < tick {
			lo = m+1
		} else {
			hi = m-1
		}
	}
	// hi is last < tick
	if hi >= 0 && hi < len(h.s) {
		return h.s[hi].X, h.s[hi].Y, true
	}
	return 0,0,false
}

type pendingTransfer struct {
	TargetZone shared.ZoneID
	StartedTick uint32

	// snapshot
	X, Y int16
	HP uint16
	Interest wire.InterestMask
	CID shared.CharacterID
	EID shared.EntityID

	// freeze movement
}

func New(cfg Config) *Server {
	if cfg.TickHz <= 0 { cfg.TickHz = 20 }
	if cfg.AOIRadius <= 0 { cfg.AOIRadius = 25 }
	if cfg.CellSize <= 0 { cfg.CellSize = 8 }
	if cfg.BudgetBytes <= 0 { cfg.BudgetBytes = 900 }
	if cfg.StateEveryTicks <= 0 { cfg.StateEveryTicks = 5 }
	if cfg.SaveEveryTicks <= 0 { cfg.SaveEveryTicks = 20 }
	if cfg.SnapshotEveryTicks <= 0 { cfg.SnapshotEveryTicks = 200 } // 10s at 20Hz
	if cfg.AIBudgetPerTick <= 0 { cfg.AIBudgetPerTick = 200 }
	if cfg.TransferTimeoutTicks == 0 { cfg.TransferTimeoutTicks = 60 } // 3s at 20Hz
	if cfg.HistoryTicks <= 0 { cfg.HistoryTicks = 40 }
	if cfg.RewindMaxTicks == 0 { cfg.RewindMaxTicks = 5 }

	if cfg.Store == nil || cfg.SaveQ == nil {
		panic("zone: Store and SaveQ required")
	}
	if cfg.SnapshotStore == nil || cfg.SnapshotQ == nil {
		panic("zone: SnapshotStore and SnapshotQ required")
	}
	if cfg.TransferTargetZone == 0 {
		panic("zone: TransferTargetZone required")
	}

	s := &Server{
		cfg: cfg,
		world: NewWorld(),
		grid: spatial.New(cfg.CellSize),
		players: make(map[shared.SessionID]*player),
		transferPending: make(map[shared.SessionID]*pendingTransfer),
		posHist: make(map[shared.EntityID]*posHistory),
		met: &metrics.Counters{},
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	// Step18: attempt snapshot load before serving
	if snap, ok, err := s.cfg.SnapshotStore.LoadSnapshot(ctx, s.cfg.ZoneID); err == nil && ok {
		s.loadSnapshotLocked(snap)
		log.Printf("zone %d loaded snapshot tick=%d ents=%d", s.cfg.ZoneID, snap.ServerTick, len(snap.Entities))
	}

	if s.cfg.HTTPAddr != "" {
		_ = s.met.Serve(s.cfg.HTTPAddr)
		log.Printf("zone metrics http=%s", s.cfg.HTTPAddr)
	}

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
			s.enqueueSnapshotLocked()
			s.mu.Unlock()
			return nil

		case fr, ok := <-inbound:
			if !ok { return errors.New("gateway link closed") }
			s.handleFrame(ctx, fr)

		case <-ticker.C:
			start := time.Now()
			s.step(ctx)
			s.met.ObserveTick(time.Since(start))
		}
	}
}

func (s *Server) handleFrame(ctx context.Context, fr wire.Frame) {
	switch fr.Type {
	case wire.MsgAttachPlayer:
		sid, cid, zid, interest, err := wire.DecodeAttachPlayer(fr.Payload)
		if err != nil || uint32(zid) != s.cfg.ZoneID {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad attach"))
			return
		}
		s.attachFromStore(ctx, sid, cid, interest)
		_ = wire.WriteFrame(s.w, wire.MsgAttachAck, nil)

	case wire.MsgAttachWithState:
		sid, cid, zid, interest, x, y, hp, err := wire.DecodeAttachWithState(fr.Payload)
		if err != nil || uint32(zid) != s.cfg.ZoneID {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadMsg, "bad attach-with-state"))
			return
		}
		s.mu.Lock()
		if _, ok := s.players[sid]; !ok {
			eid := s.world.Spawn(wire.KindPlayer, cid, x, y)
			s.posHist[eid] = newPosHistory(s.cfg.HistoryTicks)
			s.world.HP[eid] = hp
			s.players[sid] = &player{
				SID: sid, CID: cid, EID: eid,
				Interest: interest,
				known: make(map[shared.EntityID]struct{}),
				lastSentPos: make(map[shared.EntityID][2]int16),
				lastSentHP: make(map[shared.EntityID]uint16),
				pendingEvents: []string{"entered zone"},
			}
			// spawn some NPCs around on first attach to show AI/combat
			for _, ne := range s.world.RandomNearbyNPCSpawn(x, y, 3) { s.posHist[ne] = newPosHistory(s.cfg.HistoryTicks) }
		}
		s.mu.Unlock()
		_ = wire.WriteFrame(s.w, wire.MsgAttachAck, nil)

	case wire.MsgDetachPlayer:
		sid, err := wire.DecodeDetachPlayer(fr.Payload)
		if err != nil { return }
		s.mu.Lock()
		s.detachLocked(sid, "detach")
		s.mu.Unlock()

	case wire.MsgPlayerInput:
		sid, tick, mx, my, err := wire.DecodePlayerInput(fr.Payload)
		if err != nil { return }
		s.mu.Lock()
		p := s.players[sid]
		if p == nil { s.mu.Unlock(); return }
		// freeze if transfer pending
		if _, pending := s.transferPending[sid]; pending {
			s.mu.Unlock()
			return
		}
		if tick < p.nextClientTick || tick > p.nextClientTick+64 {
			s.mu.Unlock()
			return
		}
		p.nextClientTick = tick + 1
		eid := p.EID
		s.world.VelX[eid] = mx
		s.world.VelY[eid] = my
		s.mu.Unlock()

	case wire.MsgPlayerAction:
		sid, tick, skill, target, err := wire.DecodePlayerAction(fr.Payload)
		if err != nil { return }
		s.mu.Lock()
		p := s.players[sid]
		if p == nil { s.mu.Unlock(); return }
		// strict anti-cheat: use serverTick for cooldown, ignore client tick besides anti-spam window
		_ = tick
		if skill != 1 {
			s.mu.Unlock()
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadAction, "unknown skill"))
			return
		}
		// Step24: lag compensation - use client-provided action tick as claimed server tick
actionTick := tick
if actionTick == 0 || actionTick > s.serverTick || (s.serverTick-actionTick) > s.cfg.RewindMaxTicks {
	s.mu.Unlock()
	_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadAction, "bad action tick"))
	return
}
ax, ay, okA := s.posAtLocked(p.EID, actionTick)
tx, ty, okT := s.posAtLocked(target, actionTick)
if !okA || !okT {
	s.mu.Unlock()
	_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrBadAction, "no history"))
	return
}
ok, reason := s.world.ResolveSkill1At(p.EID, target, s.serverTick, ax, ay, tx, ty)
if ok {
			p.pendingEvents = append(p.pendingEvents, "hit")
		} else {
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(reason, "action rejected"))
		}
		s.mu.Unlock()

	case wire.MsgTransferCommit:
		sid, err := wire.DecodeTransferCommit(fr.Payload)
		if err != nil { return }
		s.mu.Lock()
		// finalize: remove player/entity
		if pt := s.transferPending[sid]; pt != nil {
			s.detachLocked(sid, "transfer commit")
			delete(s.transferPending, sid)
		}
		s.mu.Unlock()

	case wire.MsgTransferAbort:
		sid, err := wire.DecodeTransferAbort(fr.Payload)
		if err != nil { return }
		s.mu.Lock()
		// unfreeze by clearing pending; keep player alive
		delete(s.transferPending, sid)
		s.mu.Unlock()

	default:
	}
}

func (s *Server) attachFromStore(ctx context.Context, sid shared.SessionID, cid shared.CharacterID, interest wire.InterestMask) {
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

	eid := s.world.Spawn(wire.KindPlayer, cid, base.X, base.Y)
	s.posHist[eid] = newPosHistory(s.cfg.HistoryTicks)
	s.world.HP[eid] = base.HP

	if interest == 0 {
		interest = wire.InterestMove | wire.InterestState | wire.InterestEvent | wire.InterestCombat
	}
	s.players[sid] = &player{
		SID: sid, CID: cid, EID: eid,
		Interest: interest,
		known: make(map[shared.EntityID]struct{}),
		lastSentPos: make(map[shared.EntityID][2]int16),
		lastSentHP: make(map[shared.EntityID]uint16),
		pendingEvents: []string{"welcome"},
	}
	for _, ne := range s.world.RandomNearbyNPCSpawn(base.X, base.Y, 3) { s.posHist[ne] = newPosHistory(s.cfg.HistoryTicks) }
}

func (s *Server) detachLocked(sid shared.SessionID, why string) {
	p := s.players[sid]
	if p == nil { return }
	eid := p.EID
	// persist (enqueue only)
	s.enqueueCharacterLocked(p.CID, eid)
	s.world.Despawn(eid)
	delete(s.posHist, eid)
	delete(s.players, sid)
	delete(s.transferPending, sid)
	_ = why
}

func (s *Server) enqueueCharacterLocked(cid shared.CharacterID, eid shared.EntityID) {
	st := persist.CharacterState{
		CharacterID: cid,
		ZoneID: shared.ZoneID(s.cfg.ZoneID),
		X: s.world.PosX[eid],
		Y: s.world.PosY[eid],
		HP: s.world.HP[eid],
		ServerTick: s.serverTick,
	}
	s.cfg.SaveQ.Enqueue(st)
}

func (s *Server) enqueueDirtyLocked() {
	for _, p := range s.players {
		if s.world.Dirty[p.EID] {
			s.enqueueCharacterLocked(p.CID, p.EID)
			s.world.Dirty[p.EID] = false
		}
	}
}

func (s *Server) enqueueSnapshotLocked() {
	// snapshot world (Step18)
	snap := persist.Snapshot{
		ZoneID: s.cfg.ZoneID,
		ServerTick: s.serverTick,
		Entities: make([]persist.SnapshotEntity, 0, len(s.world.Kind)),
	}
	for eid := range s.world.Kind {
		snap.Entities = append(snap.Entities, persist.SnapshotEntity{
			EID: uint32(eid),
			Kind: uint8(s.world.Kind[eid]),
			Owner: uint64(s.world.Owner[eid]),
			X: s.world.PosX[eid], Y: s.world.PosY[eid],
			VX: s.world.VelX[eid], VY: s.world.VelY[eid],
			HP: s.world.HP[eid],
		})
	}
	s.cfg.SnapshotQ.Enqueue(s.cfg.ZoneID, snap)
}

func (s *Server) loadSnapshotLocked(snap persist.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if snap.ZoneID != s.cfg.ZoneID {
		return
	}
	s.serverTick = snap.ServerTick
	// wipe world (players will reattach later; snapshot is just world state)
	s.world = NewWorld()
	s.world.NextEID = 1
	for _, e := range snap.Entities {
		eid := shared.EntityID(e.EID)
		if eid >= s.world.NextEID {
			s.world.NextEID = eid + 1
		}
		s.world.Kind[eid] = wire.EntityKind(e.Kind)
		s.world.Owner[eid] = shared.CharacterID(e.Owner)
		s.world.PosX[eid] = e.X
		s.world.PosY[eid] = e.Y
		s.world.VelX[eid] = e.VX
		s.world.VelY[eid] = e.VY
		s.world.HP[eid] = e.HP
		s.world.Mask[eid] = wire.InterestMove | wire.InterestState | wire.InterestEvent | wire.InterestCombat
		s.posHist[eid] = newPosHistory(s.cfg.HistoryTicks)
	}
}

func (s *Server) rebuildGridLocked() {
	s.grid.Clear()
	for eid := range s.world.Kind {
		s.grid.Insert(uint32(eid), s.world.PosX[eid], s.world.PosY[eid])
	}
}

type eidDist struct { eid shared.EntityID; d2 int32 }

func (s *Server) step(ctx context.Context) {
	s.mu.Lock()
	s.serverTick++

	// Step17: AI budget + LOD (only NPCs near any player)
	aiBudget := s.cfg.AIBudgetPerTick
	if aiBudget > 0 && len(s.players) > 0 {
		// build a quick list of player positions
		playerPos := make([][2]int16, 0, len(s.players))
		for _, p := range s.players {
			playerPos = append(playerPos, [2]int16{s.world.PosX[p.EID], s.world.PosY[p.EID]})
		}
		for eid, kind := range s.world.Kind {
			if aiBudget <= 0 { break }
			if kind != wire.KindNPC { continue }
			// LOD: only update NPCs within 35 units of any player
			nx, ny := s.world.PosX[eid], s.world.PosY[eid]
			near := false
			for _, pp := range playerPos {
				dx := int32(nx) - int32(pp[0])
				dy := int32(ny) - int32(pp[1])
				if dx*dx+dy*dy <= 35*35 { near = true; break }
			}
			if !near { continue }
			s.world.WanderNPC(eid)
			aiBudget--
		}
	}

	s.world.StepPhysics()
	s.rebuildGridLocked()
// Step24: record position history (after physics)
for eid := range s.world.Kind {
	h := s.posHist[eid]
	if h == nil {
		h = newPosHistory(s.cfg.HistoryTicks)
		s.posHist[eid] = h
	}
	h.add(s.serverTick, s.world.PosX[eid], s.world.PosY[eid])
}

	// Step13: handle transfer timeouts (abort)
	for sid, pt := range s.transferPending {
		if s.serverTick - pt.StartedTick > s.cfg.TransferTimeoutTicks {
			delete(s.transferPending, sid)
			_ = wire.WriteFrame(s.w, wire.MsgError, wire.EncodeError(wire.ErrTransfer, "transfer timeout"))
		}
	}

	sendState := (s.serverTick % uint32(s.cfg.StateEveryTicks)) == 0
	doSave := (s.serverTick % uint32(s.cfg.SaveEveryTicks)) == 0
	doSnap := (s.serverTick % uint32(s.cfg.SnapshotEveryTicks)) == 0

	// detect boundary transfer and emit prepare (Step13)
	for sid, p := range s.players {
		if _, pending := s.transferPending[sid]; pending { continue }
		x := s.world.PosX[p.EID]
		if s.shouldTransfer(x) {
			// freeze movement
			s.world.VelX[p.EID] = 0
			s.world.VelY[p.EID] = 0
			pt := &pendingTransfer{
				TargetZone: shared.ZoneID(s.cfg.TransferTargetZone),
				StartedTick: s.serverTick,
				X: s.world.PosX[p.EID],
				Y: s.world.PosY[p.EID],
				HP: s.world.HP[p.EID],
				Interest: p.Interest,
				CID: p.CID,
				EID: p.EID,
			}
			s.transferPending[sid] = pt
			// enqueue save of character state
			s.enqueueCharacterLocked(p.CID, p.EID)

			st := persist.CharacterState{CharacterID: p.CID, ZoneID: shared.ZoneID(s.cfg.ZoneID), X: pt.X, Y: pt.Y, HP: pt.HP, ServerTick: s.serverTick}
			payload := wire.EncodeTransferPrepare(p.SID, p.CID, pt.TargetZone, pt.Interest, st)
			_ = wire.WriteFrame(s.w, wire.MsgTransferPrepare, payload)
			p.pendingEvents = append(p.pendingEvents, "transfer_prepare")
		}
	}

	// replication per player (AOI + interest filters)
	type perOut struct {
		sid shared.SessionID
		ev []wire.RepEvent
		move []wire.RepEvent
		state []wire.RepEvent
	}
	out := make([]perOut, 0, len(s.players))

	tmp := make([]uint32, 0, 256)
	for sid, p := range s.players {
		// if transfer pending, still allow event channel to show "loading" but stop movement/state
		_, pending := s.transferPending[sid]

		peid := p.EID
		px, py := s.world.PosX[peid], s.world.PosY[peid]

		tmp = tmp[:0]
		cands := s.grid.QueryCircle(px, py, s.cfg.AOIRadius, tmp)

		newSet := make(map[shared.EntityID]struct{}, len(cands))
		dists := make([]eidDist, 0, len(cands))
		for _, eidU := range cands {
			eid := shared.EntityID(eidU)
			newSet[eid] = struct{}{}
			ex, ey, ok := s.grid.GetPos(eidU)
			if !ok { continue }
			dx := int32(ex) - int32(px); dy := int32(ey) - int32(py)
			dists = append(dists, eidDist{eid: eid, d2: dx*dx+dy*dy})
		}
		sort.Slice(dists, func(i,j int) bool {
			if dists[i].d2 == dists[j].d2 { return dists[i].eid < dists[j].eid }
			return dists[i].d2 < dists[j].d2
		})

		ev := make([]wire.RepEvent, 0, 8)
		if p.Interest & wire.InterestEvent != 0 {
			if len(p.pendingEvents) > 0 {
				ev = append(ev, wire.RepEvent{Op: wire.RepEventText, Text: p.pendingEvents[0]})
				p.pendingEvents = p.pendingEvents[1:]
			}
		}

		move := make([]wire.RepEvent, 0, 64)
		state := make([]wire.RepEvent, 0, 16)

		if !pending && (p.Interest & wire.InterestMove != 0) {
			// despawn
			for eid := range p.known {
				if _, ok := newSet[eid]; !ok {
					move = append(move, wire.RepEvent{Op: wire.RepDespawn, EID: eid})
					delete(p.known, eid)
					delete(p.lastSentPos, eid)
					delete(p.lastSentHP, eid)
				}
			}

			// spawn/move with interest filtering
			for _, ed := range dists {
				eid := ed.eid
				mask := s.world.Mask[eid]
				if (mask & p.Interest) == 0 {
					continue
				}
				if _, ok := p.known[eid]; !ok {
					move = append(move, wire.RepEvent{
						Op: wire.RepSpawn, EID: eid, X: s.world.PosX[eid], Y: s.world.PosY[eid],
						Kind: s.world.Kind[eid], Mask: mask,
					})
					p.known[eid] = struct{}{}
					p.lastSentPos[eid] = [2]int16{s.world.PosX[eid], s.world.PosY[eid]}
					p.lastSentHP[eid] = s.world.HP[eid]
				} else {
					prev := p.lastSentPos[eid]
					if prev[0] != s.world.PosX[eid] || prev[1] != s.world.PosY[eid] {
						move = append(move, wire.RepEvent{Op: wire.RepMove, EID: eid, X: s.world.PosX[eid], Y: s.world.PosY[eid]})
						p.lastSentPos[eid] = [2]int16{s.world.PosX[eid], s.world.PosY[eid]}
					}
				}
				if len(move) >= 256 { break }
			}
		}

		if !pending && sendState && (p.Interest & wire.InterestState != 0) {
			for _, ed := range dists {
				eid := ed.eid
				mask := s.world.Mask[eid]
				if (mask & p.Interest) == 0 { continue }
				if _, ok := p.known[eid]; !ok { continue }
				if p.lastSentHP[eid] != s.world.HP[eid] {
					state = append(state, wire.RepEvent{Op: wire.RepStateHP, EID: eid, Val: s.world.HP[eid]})
					p.lastSentHP[eid] = s.world.HP[eid]
				}
				if len(state) >= 64 { break }
			}
		}

		// budget greedy: ev -> move -> state
		b := s.cfg.BudgetBytes
		ev = trimBudget(ev, b); b -= estSize(ev)
		move = trimBudget(move, b); b -= estSize(move)
		state = trimBudget(state, b)

		if len(ev)+len(move)+len(state) > 0 {
			out = append(out, perOut{sid: p.SID, ev: ev, move: move, state: state})
		}
	}

	if doSave { s.enqueueDirtyLocked() }
	if doSnap { s.enqueueSnapshotLocked() }

	// update metrics
	s.met.Entities.Store(int64(len(s.world.Kind)))
	s.met.Players.Store(int64(len(s.players)))

	s.mu.Unlock()

	for _, m := range out {
		if len(m.ev) > 0 {
			p := wire.EncodeReplicate(m.sid, s.serverTick, wire.ChanEvent, m.ev)
			s.met.AddRepBytes(len(p))
			_ = wire.WriteFrame(s.w, wire.MsgReplicate, p)
		}
		if len(m.move) > 0 {
			p := wire.EncodeReplicate(m.sid, s.serverTick, wire.ChanMove, m.move)
			s.met.AddRepBytes(len(p))
			_ = wire.WriteFrame(s.w, wire.MsgReplicate, p)
		}
		if len(m.state) > 0 {
			p := wire.EncodeReplicate(m.sid, s.serverTick, wire.ChanState, m.state)
			s.met.AddRepBytes(len(p))
			_ = wire.WriteFrame(s.w, wire.MsgReplicate, p)
		}
	}

	_ = ctx
}

func (s *Server) shouldTransfer(x int16) bool {
	b := s.cfg.TransferBoundaryX
	if b > 0 { return x > b }
	if b < 0 { return x < b }
	return false
}

// budget estimation for wire.RepEvent list (coarse upper bound)
func estSize(evs []wire.RepEvent) int {
	sz := 23
	for _, e := range evs {
		switch e.Op {
		case wire.RepSpawn:
			sz += 1+4+1+4+4
		case wire.RepMove:
			sz += 1+4+4
		case wire.RepDespawn:
			sz += 1+4
		case wire.RepStateHP:
			sz += 1+4+2
		case wire.RepEventText:
			l := len(e.Text); if l > 65535 { l = 65535 }
			sz += 1+2+l
		}
	}
	return sz
}
func trimBudget(evs []wire.RepEvent, budget int) []wire.RepEvent {
	if budget <= 23 { return evs[:0] }
	out := evs[:0]
	for _, e := range evs {
		out = append(out, e)
		if estSize(out) > budget {
			out = out[:len(out)-1]
			break
		}
	}
	return out
}
