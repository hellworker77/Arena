package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	static2 "game-server/internal/ecs/ecs_signatures/static"
)

func DealDamage(w *ecs2.World, attackerId, targetId static2.EntityID, combatAttrs static2.CombatAttributes) {
	health := ecs2.Get(w, targetId, ecs2.Health)

	if health == nil {
		return
	}

	health.Current -= combatAttrs.AttackDamage
	if health.Current < 0 {
		health.Current = 0
	}
}
