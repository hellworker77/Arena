package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/physics"
)

func ComputeSeparationForce(w *ecs.World, selfID static.EntityID, pos physics.Vector2D, minDist float32) physics.Vector2D {
	it := w.Query(ecs.CEnemyTag | ecs.CPos).Iter()

	force := physics.NewVector2D(0, 0)

	for it.Next() {
		otherID := it.EntityID()
		if otherID == selfID {
			continue
		}

		otherPos := it.Position()
		d := pos.DistanceTo(*otherPos)

		if d < minDist && d > 0 {
			diff := pos.Sub(*otherPos).Normalized()
			scale := (minDist - d) / minDist
			force = force.Add(diff.Mul(scale))
		}
	}

	return force
}
