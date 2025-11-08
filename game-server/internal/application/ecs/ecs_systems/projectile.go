package ecs_systems

import (
	"game-server/internal/application/ecs"
	"time"
)

type ProjectileSystem struct{}

func (ProjectileSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CProjectileTag | ecs.CPos | ecs.CVel | ecs.CProjectilePreset | ecs.CProjectileState).Iter()
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

func (ProjectileSystem) Reads() ecs.Signature {
	return ecs.CProjectileTag | ecs.CPos | ecs.CVel | ecs.CProjectilePreset | ecs.CProjectileState
}

func (ProjectileSystem) Writes() ecs.Signature {
	return ecs.CPos | ecs.CVel
}
