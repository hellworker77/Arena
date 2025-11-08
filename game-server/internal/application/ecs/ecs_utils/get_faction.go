package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

func GetFaction(w *ecs.World, eID static.EntityID) string {
	if w.HasComponent(eID, ecs.CPlayerTag) {
		return "player"
	}
	if w.HasComponent(eID, ecs.CPlayerTag) {
		return "enemy"
	}
	return "neutral"
}
