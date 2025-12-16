package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
	"game-server/internal/ecs/ecs_signatures/runtime"
	ecs_utils2 "game-server/internal/ecs/ecs_utils"
)

type EnemyBehaviorSystem struct{}

func (EnemyBehaviorSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CEnemyTag | ecs2.CPos | ecs2.CVel | ecs2.CEnemyPreset | ecs2.CMovementAttrs).Iter()

	for it.Next() {
		eID := it.EntityID()
		pos := it.Position()
		preset := it.EnemyPreset()
		movementAttrs := it.MovementAttributes()

		_, playerPos, ok := ecs_utils2.FindClosestPlayer(w, *pos, preset.AggroRange)
		if !ok {
			continue
		}

		targetDir := playerPos.Sub(*pos).Normalized()
		avoid := ecs_utils2.ComputeSeparationForce(w, eID, *pos, preset.KeepDistance)

		move := targetDir.Add(avoid).Normalized().Mul(movementAttrs.MoveSpeed)

		newPos := pos.Add(move.Mul(dt))

		ecs2.Set(w, eID, ecs2.Position, runtime.Position{X: newPos.X, Y: newPos.Y})
	}
}

func (EnemyBehaviorSystem) Reads() ecs2.Signature {
	return ecs2.CEnemyTag | ecs2.CPos | ecs2.CVel | ecs2.CEnemyPreset | ecs2.CMovementAttrs
}

func (EnemyBehaviorSystem) Writes() ecs2.Signature {
	return ecs2.CVel | ecs2.CPos
}
