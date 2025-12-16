package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	"game-server/internal/ecs/ecs_signatures/static"
)

func HandleDeath(w *ecs2.World, eID static.EntityID) {
	rec, ok := w.GetEntity(eID)
	if !ok {
		return
	}

	newSig := rec.Archetype.Signature &^ ecs2.CHealth
	w.MoveEntity(eID, newSig)
}
