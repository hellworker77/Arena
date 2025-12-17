package zone

import (
	"math"
	"math/rand"
	"time"

	"game-server/internal/shared"
	"game-server/internal/shared/wire"
)

type World struct {
	NextEID shared.EntityID

	// Component stores (ECS-ish)
	Kind   map[shared.EntityID]wire.EntityKind
	Owner  map[shared.EntityID]shared.CharacterID
	PosX   map[shared.EntityID]int16
	PosY   map[shared.EntityID]int16
	VelX   map[shared.EntityID]int16
	VelY   map[shared.EntityID]int16
	HP     map[shared.EntityID]uint16
	Mask   map[shared.EntityID]wire.InterestMask

	Dirty  map[shared.EntityID]bool

	// Cooldowns: eid -> next serverTick allowed for skill1
	Skill1CD map[shared.EntityID]uint32
}

func NewWorld() *World {
	rand.Seed(time.Now().UnixNano())
	return &World{
		NextEID: 1,
		Kind: make(map[shared.EntityID]wire.EntityKind),
		Owner: make(map[shared.EntityID]shared.CharacterID),
		PosX: make(map[shared.EntityID]int16),
		PosY: make(map[shared.EntityID]int16),
		VelX: make(map[shared.EntityID]int16),
		VelY: make(map[shared.EntityID]int16),
		HP: make(map[shared.EntityID]uint16),
		Mask: make(map[shared.EntityID]wire.InterestMask),
		Dirty: make(map[shared.EntityID]bool),
		Skill1CD: make(map[shared.EntityID]uint32),
	}
}

func (w *World) Spawn(kind wire.EntityKind, owner shared.CharacterID, x, y int16) shared.EntityID {
	eid := w.NextEID
	w.NextEID++
	w.Kind[eid] = kind
	w.Owner[eid] = owner
	w.PosX[eid] = x
	w.PosY[eid] = y
	w.VelX[eid] = 0
	w.VelY[eid] = 0
	if kind == wire.KindNPC {
		w.HP[eid] = 50
	} else {
		w.HP[eid] = 100
	}
	w.Mask[eid] = wire.InterestMove | wire.InterestState | wire.InterestEvent | wire.InterestCombat
	w.Dirty[eid] = true
	return eid
}

func (w *World) Despawn(eid shared.EntityID) {
	delete(w.Kind, eid)
	delete(w.Owner, eid)
	delete(w.PosX, eid)
	delete(w.PosY, eid)
	delete(w.VelX, eid)
	delete(w.VelY, eid)
	delete(w.HP, eid)
	delete(w.Mask, eid)
	delete(w.Dirty, eid)
	delete(w.Skill1CD, eid)
}

func (w *World) StepPhysics() {
	for eid := range w.Kind {
		vx := w.VelX[eid]
		vy := w.VelY[eid]
		if vx != 0 || vy != 0 {
			w.PosX[eid] += vx
			w.PosY[eid] += vy
			w.Dirty[eid] = true
		}
	}
}

func dist2(ax, ay, bx, by int16) int32 {
	dx := int32(ax) - int32(bx)
	dy := int32(ay) - int32(by)
	return dx*dx + dy*dy
}

func within(ax, ay, bx, by int16, r int16) bool {
	return dist2(ax, ay, bx, by) <= int32(r)*int32(r)
}

func clampInt16(v int16, lo, hi int16) int16 {
	if v < lo { return lo }
	if v > hi { return hi }
	return v
}

func (w *World) WanderNPC(eid shared.EntityID) {
	// tiny wander: random velocity in [-1,1]
	w.VelX[eid] = int16(rand.Intn(3) - 1)
	w.VelY[eid] = int16(rand.Intn(3) - 1)
}

// ResolveSkill1: toy melee hit, server-authoritative (Step15)
func (w *World) ResolveSkill1(attacker, target shared.EntityID, serverTick uint32) (ok bool, reason wire.ErrCode) {
	if w.Kind[attacker] != wire.KindPlayer {
		return false, wire.ErrBadAction
	}
	if _, exists := w.Kind[target]; !exists {
		return false, wire.ErrBadAction
	}
	if serverTick < w.Skill1CD[attacker] {
		return false, wire.ErrCooldown
	}
	ax, ay := w.PosX[attacker], w.PosY[attacker]
	tx, ty := w.PosX[target], w.PosY[target]
	if !within(ax, ay, tx, ty, 4) { // melee range
		return false, wire.ErrOutOfRange
	}

	hp := w.HP[target]
	if hp == 0 {
		return false, wire.ErrBadAction
	}
	dmg := uint16(5)
	if hp <= dmg {
		w.HP[target] = 0
	} else {
		w.HP[target] = hp - dmg
	}
	w.Dirty[target] = true
	// 10 ticks cooldown
	w.Skill1CD[attacker] = serverTick + 10
	return true, 0
}

func (w *World) RandomNearbyNPCSpawn(centerX, centerY int16, n int) []shared.EntityID {
	out := make([]shared.EntityID, 0, n)
	for i := 0; i < n; i++ {
		ang := rand.Float64() * 2 * math.Pi
		r := rand.Float64() * 10
		x := centerX + int16(math.Round(math.Cos(ang)*r))
		y := centerY + int16(math.Round(math.Sin(ang)*r))
		out = append(out, w.Spawn(wire.KindNPC, 0, x, y))
	}
	return out
}
