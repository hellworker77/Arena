package ecs_systems

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_utils"
)

type ColliderSystem struct{}

func (c ColliderSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CCollider | ecs.CPos).Iter()

	for it.Next() {
		posA := it.Position()
		colliderA := it.Collider()
		eID := it.EntityID()

		nearIDs := w.Grid.QueryNearby(*posA)

		for _, otherID := range nearIDs {
			if otherID == eID {
				continue
			}
			rec, ok := w.GetEntity(otherID)
			if !ok || rec.Archetype.Signature&(ecs.CPos|ecs.CCollider) != (ecs.CPos|ecs.CCollider) {
				continue
			}

			posB := rec.Archetype.Positions[rec.Index]
			colB := rec.Archetype.Colliders[rec.Index]

			if ecs_utils.CheckCollision(*posA, posB, *colliderA, colB) {
				ecs_utils.HandleCollision(w, eID, otherID)
			}
		}

	}
}

func (ColliderSystem) Reads() ecs.Signature {
	return ecs.CCollider | ecs.CPos
}

func (ColliderSystem) Writes() ecs.Signature {
	return ecs.CPos | ecs.CHealth | ecs.CVel
}

/// TODO: Spatial partitioning: to reduce the number of collision checks
/// TODO: Layer / Mask: For example, players don't collide with other players
/// TODO: Interest management: handle only nearby players
