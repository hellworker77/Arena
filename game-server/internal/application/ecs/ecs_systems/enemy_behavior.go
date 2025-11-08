package ecs_systems

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_utils"
)

type EnemyBehaviorSystem struct{}

func (EnemyBehaviorSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CEnemyTag | ecs.CPos | ecs.CVel | ecs.CEnemyPreset | ecs.CMovementAttrs).Iter()

	for it.Next() {
		eID := it.EntityID()
		pos := it.Position()
		preset := it.EnemyPreset()
		movementAttrs := it.MovementAttributes()

		_, playerPos, ok := ecs_utils.FindClosestPlayer(w, *pos, preset.AggroRange)
		if !ok {
			continue
		}

		targetDir := playerPos.Sub(*pos).Normalized()
		avoid := ecs_utils.ComputeSeparationForce(w, eID, *pos, preset.KeepDistance)

		move := targetDir.Add(avoid).Normalized().Mul(movementAttrs.MoveSpeed)

		newPos := pos.Add(move.Mul(dt))

		ecs.Set(w, eID, ecs.Position, runtime.Position{X: newPos.X, Y: newPos.Y})
	}
}

func (EnemyBehaviorSystem) Reads() ecs.Signature {
	return ecs.CEnemyTag | ecs.CPos | ecs.CVel | ecs.CEnemyPreset | ecs.CMovementAttrs
}

func (EnemyBehaviorSystem) Writes() ecs.Signature {
	return ecs.CVel | ecs.CPos
}
