package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/physics"
)

func FindClosestPlayer(w *ecs.World, pos runtime.Position, rangeDist float32) (static.EntityID, physics.Vector2D, bool) {
	it := w.Query(ecs.CPlayerTag | ecs.CPos).Iter()

	var bestId static.EntityID
	var bestPos runtime.Position
	bestDist := rangeDist

	for it.Next() {
		pPos := it.Pos()
		dist := pos.DistanceTo(*pPos)
		if dist < bestDist {
			bestDist = dist
			bestId = it.EntityID()
			bestPos = *pPos
		}
	}

	return bestId, bestPos, bestDist < rangeDist
}
