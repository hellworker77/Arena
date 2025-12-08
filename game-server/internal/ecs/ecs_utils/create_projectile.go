package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	runtime2 "game-server/internal/ecs/ecs_signatures/runtime"
	static2 "game-server/internal/ecs/ecs_signatures/static"
	"time"
)

func CreateProjectile(w *ecs2.World, attackerId, targetId static2.EntityID, combatAttrs static2.CombatAttributes, prjPreset static2.ProjectilePreset) {
	aRec, aOk := w.GetEntity(attackerId)
	if !aOk {
		return
	}
	aPos := &aRec.Archetype.Positions[aRec.Index]

	tRec, tOk := w.GetEntity(targetId)
	if !tOk {
		return
	}

	tPos := &tRec.Archetype.Positions[tRec.Index]

	dir := tPos.Sub(*aPos).Normalized()
	vel := dir.Mul(prjPreset.Speed)

	eID := w.CreateEntity(ecs2.CPos | ecs2.CVel | ecs2.CProjectileTag | ecs2.CLifespan)

	ecs2.Set(w, eID, ecs2.Position, *aPos)
	ecs2.Set(w, eID, ecs2.Velocity, vel)
	ecs2.Set(w, eID, ecs2.Lifespan, runtime2.Lifespan{CreatedAtUnix: time.Now().Unix(), DurationSecs: 5})
	ecs2.Set(w, eID, ecs2.ProjectileState, runtime2.ProjectileState{
		OwnerID:  attackerId,
		TargetID: targetId,
		SpawnPos: *aPos,
		Elapsed:  0,
	})
}
