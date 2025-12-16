package game

import (
	"errors"
	"game-server/internal/ecs"
	runtimeSig "game-server/internal/ecs/ecs_signatures/runtime"
	staticSig "game-server/internal/ecs/ecs_signatures/static"
	ecsSystems "game-server/internal/ecs/ecs_systems"
	"game-server/pkg/protocol"
)

const (
	// maxBufferedInputsPerPlayer bounds memory and worst-case per-tick work.
	maxBufferedInputsPerPlayer = 64

	// maxInputLeadTicks rejects client ticks that are too far in the future
	// relative to what the server expects for this player. This prevents a client
	// from spamming huge tick numbers and forcing unbounded buffering.
	maxInputLeadTicks = 128

	// maxInputLateTicks rejects client ticks that are too far in the past.
	// This avoids churn from duplicates / replays.
	maxInputLateTicks = 8
)

type playerState struct {
	entity staticSig.EntityID

	// haveBase becomes true on first input. Until then, the server keeps the player still.
	haveBase bool

	// nextTick is the next client tick the server will try to consume.
	nextTick uint32

	// lastApplied is the last input we applied (or synthesized if missing).
	lastApplied protocol.Input

	// buffer holds inputs keyed by client tick.
	buffer map[uint32]protocol.Input
}

// Engine owns the authoritative simulation state.
// It is deliberately single-threaded: call all methods from the simulation loop.
type Engine struct {
	world   *ecs.World
	systems []ecsSystems.ECSSystem

	// players maps an external client key (currently: addr string) to per-player state.
	players map[string]*playerState

	// per-client replication state for delta snapshots.
	snap map[string]*clientSnapState
}

type clientSnapState struct {
	seq          uint32
	lastFullTick uint32
	lastSent     map[uint32]protocol.EntityState // entity id -> last state
}

func NewEngine() *Engine {
	return &Engine{
		world:   ecs.NewWorld(),
		systems: []ecsSystems.ECSSystem{ecsSystems.MovementSystem{}},
		players: map[string]*playerState{},
		snap:    map[string]*clientSnapState{},
	}
}

func (e *Engine) World() *ecs.World { return e.world }

func (e *Engine) PlayerEntity(clientKey string) (staticSig.EntityID, bool) {
	ps, ok := e.players[clientKey]
	if !ok {
		return 0, false
	}
	return ps.entity, true
}

// AddPlayer creates a new player entity and returns its id.
func (e *Engine) AddPlayer(clientKey string) staticSig.EntityID {
	if ps, ok := e.players[clientKey]; ok {
		return ps.entity
	}

	ent := e.world.CreateEntity(ecs.CPlayerTag | ecs.CPos | ecs.CVel | ecs.CHealth)
	ecs.Set(e.world, ent, ecs.Position, runtimeSig.Position{})
	ecs.Set(e.world, ent, ecs.Velocity, runtimeSig.Velocity{})
	ecs.Set(e.world, ent, ecs.Health, runtimeSig.Health{Current: 100})

	e.players[clientKey] = &playerState{
		entity: ent,
		buffer: make(map[uint32]protocol.Input, 16),
	}
	// Initialize replication baseline for this client.
	if _, ok := e.snap[clientKey]; !ok {
		e.snap[clientKey] = &clientSnapState{lastSent: make(map[uint32]protocol.EntityState, 64)}
	}
	return ent
}

func (e *Engine) RemovePlayer(clientKey string) {
	// World doesn't currently expose entity destruction; keep mapping removal only.
	delete(e.players, clientKey)
	delete(e.snap, clientKey)
}

// QueueInput stores an input for later consumption by Step().
// It enforces ordering constraints based on clientTick.
func (e *Engine) QueueInput(clientKey string, in protocol.Input) error {
	ps, ok := e.players[clientKey]
	if !ok {
		return errors.New("unknown player")
	}

	// Initialize base tick on first input.
	if !ps.haveBase {
		ps.haveBase = true
		ps.nextTick = in.ClientTick
		ps.lastApplied = protocol.Input{ClientTick: in.ClientTick, Speed: in.Speed}
	}

	// Reject too-old inputs (duplicates / replays).
	if in.ClientTick+maxInputLateTicks < ps.nextTick {
		return errors.New("stale input")
	}

	// Reject unreasonably-future inputs (buffer bloat / abuse).
	if in.ClientTick > ps.nextTick+maxInputLeadTicks {
		return errors.New("input too far ahead")
	}

	// Store / overwrite (duplicates collapse naturally).
	ps.buffer[in.ClientTick] = in

	// Bound buffer size: if it grows too large, drop the farthest-future ticks first.
	for len(ps.buffer) > maxBufferedInputsPerPlayer {
		var maxTick uint32
		first := true
		for t := range ps.buffer {
			if first || t > maxTick {
				maxTick = t
				first = false
			}
		}
		delete(ps.buffer, maxTick)
	}

	return nil
}

// Step runs one simulation tick.
func (e *Engine) Step(dt float32) {
	// Consume one tick worth of input per player.
	for _, ps := range e.players {
		if !ps.haveBase {
			// Keep still until we see first input.
			vel := ecs.Get(e.world, ps.entity, ecs.Velocity)
			vel.X = 0
			vel.Y = 0
			continue
		}

		if in, ok := ps.buffer[ps.nextTick]; ok {
			ps.lastApplied = in
			delete(ps.buffer, ps.nextTick)
		} else {
			// Missing tick: synthesize "no movement" to avoid stuck-walking forever.
			tmp := ps.lastApplied
			tmp.ClientTick = ps.nextTick
			tmp.MoveX = 0
			tmp.MoveY = 0
			ps.lastApplied = tmp
		}
		ps.nextTick++

		vel := ecs.Get(e.world, ps.entity, ecs.Velocity)
		vel.X = ps.lastApplied.MoveX * ps.lastApplied.Speed
		vel.Y = ps.lastApplied.MoveY * ps.lastApplied.Speed
	}

	for _, sys := range e.systems {
		sys.Run(e.world, dt)
	}
}

// BuildSnapshotForClient builds either a full or delta snapshot payload for a given client.
// This method must be called only from the simulation thread.
//
// interestRadius controls which entities are relevant (simple distance culling).
// fullEveryTicks forces a full snapshot at that cadence to bound drift.
func (e *Engine) BuildSnapshotForClient(clientKey string, serverTick uint32, interestRadius float32, fullEveryTicks uint32) ([]byte, bool) {
	ps, ok := e.players[clientKey]
	if !ok {
		return nil, false
	}
	cs, ok := e.snap[clientKey]
	if !ok {
		cs = &clientSnapState{lastSent: make(map[uint32]protocol.EntityState, 64)}
		e.snap[clientKey] = cs
	}
	cs.seq++

	// Determine client's position.
	center := ecs.Get(e.world, ps.entity, ecs.Position)
	r2 := interestRadius * interestRadius

	// Collect relevant entities (currently: all players, culled by radius).
	states := make(map[uint32]protocol.EntityState, 64)
	for _, other := range e.players {
		pos := ecs.Get(e.world, other.entity, ecs.Position)
		dx := pos.X - center.X
		dy := pos.Y - center.Y
		if dx*dx+dy*dy > r2 {
			continue
		}
		vel := ecs.Get(e.world, other.entity, ecs.Velocity)
		id := uint32(other.entity)
		states[id] = protocol.EntityState{ID: id, X: pos.X, Y: pos.Y, VX: vel.X, VY: vel.Y}
	}

	forceFull := fullEveryTicks > 0 && (serverTick == 0 || serverTick-cs.lastFullTick >= fullEveryTicks)
	if forceFull || len(cs.lastSent) == 0 {
		out := make([]protocol.EntityState, 0, len(states))
		for _, st := range states {
			out = append(out, st)
		}
		cs.lastSent = states
		cs.lastFullTick = serverTick
		return protocol.MarshalSnapshotFull(serverTick, cs.seq, out), true
	}

	// Delta: compute upserts + removes.
	upserts := make([]protocol.EntityState, 0, 32)
	removes := make([]uint32, 0, 16)

	for id, cur := range states {
		prev, ok := cs.lastSent[id]
		if !ok || prev.X != cur.X || prev.Y != cur.Y || prev.VX != cur.VX || prev.VY != cur.VY {
			upserts = append(upserts, cur)
		}
	}
	for id := range cs.lastSent {
		if _, ok := states[id]; !ok {
			removes = append(removes, id)
		}
	}

	// Update baseline.
	cs.lastSent = states
	return protocol.MarshalSnapshotDelta(serverTick, cs.seq, upserts, removes), true
}
