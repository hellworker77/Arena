package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	"game-server/internal/ecs/ecs_signatures/static"
)

func GetFaction(w *ecs2.World, eID static.EntityID) string {
	if w.HasComponent(eID, ecs2.CPlayerTag) {
		return "player"
	}
	if w.HasComponent(eID, ecs2.CPlayerTag) {
		return "enemy"
	}
	return "neutral"
}
