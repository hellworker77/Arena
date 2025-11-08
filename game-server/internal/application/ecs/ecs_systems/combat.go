package ecs_systems

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_utils"
	"time"
)

type CombatSystem struct{}

func (CombatSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CTarget | ecs.CCombatAttrs | ecs.CAttackCooldown | ecs.CProjectilePreset).Iter()
	now := time.Now().Unix()

	for it.Next() {
		attackerID := it.EntityID()
		combatAttributes := it.CombatAttributes()
		projectilePreset := it.ProjectilePreset()
		cd := it.AttackCooldown()
		targetID := it.Target().ID
		attackerPos := it.Position()
		targetPos := ecs.Get(w, targetID, ecs.Position)

		if targetPos == nil || attackerPos == nil || combatAttributes == nil || projectilePreset == nil {
			continue
		}

		attackDelay := 1.0 / combatAttributes.AttackSpeed
		if float32(now-cd.LastAttackUnix) < attackDelay {
			continue
		}

		if combatAttributes.AttackRange < attackerPos.DistanceTo(*targetPos) {
			ecs_utils.DealDamage(w, attackerID, targetID, *combatAttributes)
		} else {
			ecs_utils.CreateProjectile(w, attackerID, targetID, *combatAttributes, *projectilePreset)
		}

		cd.LastAttackUnix = now
	}
}

func (CombatSystem) Reads() ecs.Signature {
	return ecs.CTarget | ecs.CCombatAttrs | ecs.CAttackCooldown | ecs.CProjectilePreset
}

func (CombatSystem) Writes() ecs.Signature {
	return ecs.CHealth | ecs.CAttackCooldown | ecs.CProjectileState
}
