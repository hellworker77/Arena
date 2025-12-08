package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
	"time"
)

type ProjectileSystem struct{}

func (ProjectileSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CProjectileTag | ecs2.CPos | ecs2.CVel | ecs2.CProjectilePreset | ecs2.CProjectileState).Iter()
	now := time.Now().Unix()

	for it.Next() {
		eID := it.EntityID()
		pos := it.Position()
		vel := it.Velocity()
		preset := it.ProjectilePreset()
		state := it.ProjectileState()
		life := it.Lifespan()

		if now-life.CreatedAtUnix > life.DurationSecs {
			w.RemoveEntity(eID)
			continue
		}

		if preset.IsHoming && state.TargetID != 0 {
			tRec, ok := w.GetEntity(state.TargetID)
			if ok {
				tPos := tRec.Archetype.Positions[tRec.Index]
				dir := tPos.Sub(*pos).Normalized()
				vel.X = dir.X * vel.Magnitude()
				vel.Y = dir.Y * vel.Magnitude()
			}
		}

		pos.X += vel.X * dt
		pos.Y += vel.Y * dt
	}
}

func (ProjectileSystem) Reads() ecs2.Signature {
	return ecs2.CProjectileTag | ecs2.CPos | ecs2.CVel | ecs2.CProjectilePreset | ecs2.CProjectileState
}

func (ProjectileSystem) Writes() ecs2.Signature {
	return ecs2.CPos | ecs2.CVel
}
