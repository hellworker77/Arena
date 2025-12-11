package ecs_utils

import (
	ecs2 "game-server/internal/ecs"
	static2 "game-server/internal/ecs/ecs_signatures/static"
)

type CollisionHandler func(w *ecs2.World, aID, bID static2.EntityID)

var collisionHandlers = map[static2.CollisionGroup]map[static2.CollisionGroup]CollisionHandler{
	static2.LayerPlayer: {
		static2.LayerWorld:      handlePlayerWorldCollision,
		static2.LayerEnemy:      handlePlayerEnemyCollision,
		static2.LayerItem:       handlePlayerItemCollision,
		static2.LayerProjectile: handlePlayerProjectileCollision,
	},
	static2.LayerEnemy: {
		static2.LayerEnemy:      handleEnemyEnemyCollision,
		static2.LayerProjectile: handleEnemyProjectileCollision,
	},
	static2.LayerProjectile: {
		static2.LayerEnemy:  handleEnemyProjectileCollision,
		static2.LayerPlayer: handlePlayerProjectileCollision,
	},
}

func handlePlayerWorldCollision(w *ecs2.World, playerID, wallID static2.EntityID) {
	vel := ecs2.Get(w, playerID, ecs2.Velocity)
	if vel == nil {
		return
	}
	vel.X, vel.Y = 0, 0
}

func handlePlayerEnemyCollision(w *ecs2.World, playerID, enemyID static2.EntityID) {
	// no reaction for now
}

func handlePlayerItemCollision(w *ecs2.World, playerID, itemID static2.EntityID) {
	// no reaction for now
}

func handlePlayerProjectileCollision(w *ecs2.World, playerID, projectileID static2.EntityID) {
	projState := ecs2.Get(w, projectileID, ecs2.ProjectileState)
	projPreset := ecs2.Get(w, projectileID, ecs2.ProjectilePreset)

	if projState == nil || projPreset == nil {
		return
	}

	ownerFaction := GetFaction(w, projState.OwnerID)
	targetFaction := GetFaction(w, playerID)
	combatAttrs := ecs2.Get(w, projState.OwnerID, ecs2.CombatAttributes)

	if ownerFaction == targetFaction {
		return
	}

	health := ecs2.Get(w, playerID, ecs2.Health)

	if health != nil && combatAttrs != nil {
		health.Current -= combatAttrs.AttackDamage
	}

	w.RemoveEntity(projectileID)
}

func handleEnemyEnemyCollision(w *ecs2.World, enemyAID, enemyBID static2.EntityID) {
	posA := ecs2.Get(w, enemyAID, ecs2.Position)
	posB := ecs2.Get(w, enemyBID, ecs2.Position)

	if posA == nil || posB == nil {
		return
	}

	dir := posA.Sub(*posB).Normalized()
	shift := dir.Mul(0.1)

	posA.Add(shift)
	posB.Add(shift.Neg())
}

func handleEnemyProjectileCollision(w *ecs2.World, enemyID, projectileID static2.EntityID) {
	projState := ecs2.Get(w, projectileID, ecs2.ProjectileState)
	projPreset := ecs2.Get(w, projectileID, ecs2.ProjectilePreset)

	if projState == nil || projPreset == nil {
		return
	}

	ownerFaction := GetFaction(w, projState.OwnerID)
	targetFaction := GetFaction(w, enemyID)
	combatAttrs := ecs2.Get(w, projState.OwnerID, ecs2.CombatAttributes)

	if ownerFaction == targetFaction {
		return
	}

	health := ecs2.Get(w, enemyID, ecs2.Health)

	if health != nil && combatAttrs != nil {
		health.Current -= combatAttrs.AttackDamage
	}

	w.RemoveEntity(projectileID)
}

func HandleCollision(w *ecs2.World, aID, bID static2.EntityID) {
	a := ecs2.Get(w, aID, ecs2.Collider)
	b := ecs2.Get(w, bID, ecs2.Collider)

	if a == nil || b == nil {
		return
	}

	if (a.Mask & static2.CollisionGroup(b.Layer)) == 0 {
		return
	}
	if (b.Mask & static2.CollisionGroup(a.Layer)) == 0 {
		return
	}

	if handlers, ok := collisionHandlers[static2.CollisionGroup(a.Layer)]; ok {
		if handler, ok := handlers[static2.CollisionGroup(b.Layer)]; ok {
			handler(w, aID, bID)
			return
		}
	}
}
