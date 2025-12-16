package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
	ecs_utils2 "game-server/internal/ecs/ecs_utils"
	"time"
)

type CombatSystem struct{}

func (CombatSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CTarget | ecs2.CCombatAttrs | ecs2.CAttackCooldown | ecs2.CProjectilePreset).Iter()
	now := time.Now().Unix()

	for it.Next() {
		attackerID := it.EntityID()
		combatAttributes := it.CombatAttributes()
		projectilePreset := it.ProjectilePreset()
		cd := it.AttackCooldown()
		targetID := it.Target().ID
		attackerPos := it.Position()
		targetPos := ecs2.Get(w, targetID, ecs2.Position)

		if targetPos == nil || attackerPos == nil || combatAttributes == nil || projectilePreset == nil {
			continue
		}

		attackDelay := 1.0 / combatAttributes.AttackSpeed
		if float32(now-cd.LastAttackUnix) < attackDelay {
			continue
		}

		if combatAttributes.AttackRange < attackerPos.DistanceTo(*targetPos) {
			ecs_utils2.DealDamage(w, attackerID, targetID, *combatAttributes)
		} else {
			ecs_utils2.CreateProjectile(w, attackerID, targetID, *combatAttributes, *projectilePreset)
		}

		cd.LastAttackUnix = now
	}
}

func (CombatSystem) Reads() ecs2.Signature {
	return ecs2.CTarget | ecs2.CCombatAttrs | ecs2.CAttackCooldown | ecs2.CProjectilePreset
}

func (CombatSystem) Writes() ecs2.Signature {
	return ecs2.CHealth | ecs2.CAttackCooldown | ecs2.CProjectileState
}
