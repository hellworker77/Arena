package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"time"
)

func CreateProjectile(w *ecs.World, attackerId, targetId static.EntityID, combatAttrs static.CombatAttributes, prjPreset static.ProjectilePreset) {
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

	eID := w.CreateEntity(ecs.CPos | ecs.CVel | ecs.CProjectileTag | ecs.CLifespan)

	ecs.Set(w, eID, ecs.Position, *aPos)
	ecs.Set(w, eID, ecs.Velocity, vel)
	ecs.Set(w, eID, ecs.Lifespan, runtime.Lifespan{CreatedAtUnix: time.Now().Unix(), DurationSecs: 5})
	ecs.Set(w, eID, ecs.ProjectileState, runtime.ProjectileState{
		OwnerID:  attackerId,
		TargetID: targetId,
		SpawnPos: *aPos,
		Elapsed:  0,
	})
}
