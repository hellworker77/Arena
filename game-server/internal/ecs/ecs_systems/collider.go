package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
	ecs_utils2 "game-server/internal/ecs/ecs_utils"
)

type ColliderSystem struct{}

func (c ColliderSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CCollider | ecs2.CPos).Iter()

	for it.Next() {
		posA := it.Position()
		colliderA := it.Collider()
		eID := it.EntityID()

		nearIDs := w.Grid.Query(1, *posA, 10.0)

		for _, otherID := range nearIDs {
			if otherID == eID {
				continue
			}
			rec, ok := w.GetEntity(otherID)
			if !ok || rec.Archetype.Signature&(ecs2.CPos|ecs2.CCollider) != (ecs2.CPos|ecs2.CCollider) {
				continue
			}

			posB := rec.Archetype.Positions[rec.Index]
			colB := rec.Archetype.Colliders[rec.Index]

			if ecs_utils2.CheckCollision(*posA, posB, *colliderA, colB) {
				ecs_utils2.HandleCollision(w, eID, otherID)
			}
		}

	}
}

func (ColliderSystem) Reads() ecs2.Signature {
	return ecs2.CCollider | ecs2.CPos
}

func (ColliderSystem) Writes() ecs2.Signature {
	return ecs2.CPos | ecs2.CHealth | ecs2.CVel
}

/// TODO: Spatial partitioning: to reduce the number of collision checks
/// TODO: Layer / Mask: For example, players don't collide with other players
/// TODO: Interest management: handle only nearby players
