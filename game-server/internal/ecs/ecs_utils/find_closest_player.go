package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	"game-server/internal/ecs/ecs_signatures/runtime"
	"game-server/internal/ecs/ecs_signatures/static"
	"game-server/pkg/physics"
)

func FindClosestPlayer(w *ecs2.World, pos runtime.Position, rangeDist float32) (static.EntityID, physics.Vector2D, bool) {
	it := w.Query(ecs2.CPlayerTag | ecs2.CPos).Iter()

	var bestId static.EntityID
	var bestPos runtime.Position
	bestDist := rangeDist

	for it.Next() {
		pPos := it.Position()
		dist := pos.DistanceTo(*pPos)
		if dist < bestDist {
			bestDist = dist
			bestId = it.EntityID()
			bestPos = *pPos
		}
	}

	return bestId, bestPos, bestDist < rangeDist
}
