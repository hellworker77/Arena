package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

func HandleDeath(w *ecs.World, eID static.EntityID) {
	rec, ok := w.GetEntity(eID)
	if !ok {
		return
	}

	newSig := rec.Archetype.Signature &^ ecs.CStats
	w.MoveEntity(eID, newSig)
}
