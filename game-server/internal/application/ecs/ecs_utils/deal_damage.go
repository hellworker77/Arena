package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

func DealDamage(w *ecs.World, attackerId, targetId static.EntityID, combatAttrs static.CombatAttributes) {
	health := ecs.Get(w, targetId, ecs.Health)

	if health == nil {
		return
	}

	health.Current -= combatAttrs.AttackDamage
	if health.Current < 0 {
		health.Current = 0
	}
}
