package ecs_utils

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

type CollisionHandler func(w *ecs.World, aID, bID static.EntityID)

var collisionHandlers = map[static.CollisionGroup]map[static.CollisionGroup]CollisionHandler{
	static.LayerPlayer: {
		static.LayerWorld:      handlePlayerWorldCollision,
		static.LayerEnemy:      handlePlayerEnemyCollision,
		static.LayerItem:       handlePlayerItemCollision,
		static.LayerProjectile: handlePlayerProjectileCollision,
	},
	static.LayerEnemy: {
		static.LayerEnemy:      handleEnemyEnemyCollision,
		static.LayerProjectile: handleEnemyProjectileCollision,
	},
	static.LayerProjectile: {
		static.LayerEnemy:  handleEnemyProjectileCollision,
		static.LayerPlayer: handlePlayerProjectileCollision,
	},
}

func handlePlayerWorldCollision(w *ecs.World, playerID, wallID static.EntityID) {
	vel := ecs.Get(w, playerID, ecs.Velocity)
	if vel == nil {
		return
	}
	vel.X, vel.Y = 0, 0
}

func handlePlayerEnemyCollision(w *ecs.World, playerID, enemyID static.EntityID) {
	// no reaction for now
}

func handlePlayerItemCollision(w *ecs.World, playerID, itemID static.EntityID) {
	// no reaction for now
}

func handlePlayerProjectileCollision(w *ecs.World, playerID, projectileID static.EntityID) {
	projState := ecs.Get(w, projectileID, ecs.ProjectileState)
	projPreset := ecs.Get(w, projectileID, ecs.ProjectilePreset)

	if projState == nil || projPreset == nil {
		return
	}

	ownerFaction := GetFaction(w, projState.OwnerID)
	targetFaction := GetFaction(w, playerID)
	combatAttrs := ecs.Get(w, projState.OwnerID, ecs.CombatAttributes)

	if ownerFaction == targetFaction {
		return
	}

	health := ecs.Get(w, playerID, ecs.Health)

	if health != nil && combatAttrs != nil {
		health.Current -= combatAttrs.AttackDamage
	}

	w.RemoveEntity(projectileID)
}

func handleEnemyEnemyCollision(w *ecs.World, enemyAID, enemyBID static.EntityID) {
	posA := ecs.Get(w, enemyAID, ecs.Position)
	posB := ecs.Get(w, enemyBID, ecs.Position)

	if posA == nil || posB == nil {
		return
	}

	dir := posA.Sub(*posB).Normalized()
	shift := dir.Mul(0.1)

	posA.Add(shift)
	posB.Add(shift.Neg())
}

func handleEnemyProjectileCollision(w *ecs.World, enemyID, projectileID static.EntityID) {
	projState := ecs.Get(w, projectileID, ecs.ProjectileState)
	projPreset := ecs.Get(w, projectileID, ecs.ProjectilePreset)

	if projState == nil || projPreset == nil {
		return
	}

	ownerFaction := GetFaction(w, projState.OwnerID)
	targetFaction := GetFaction(w, enemyID)
	combatAttrs := ecs.Get(w, projState.OwnerID, ecs.CombatAttributes)

	if ownerFaction == targetFaction {
		return
	}

	health := ecs.Get(w, enemyID, ecs.Health)

	if health != nil && combatAttrs != nil {
		health.Current -= combatAttrs.AttackDamage
	}

	w.RemoveEntity(projectileID)
}

func HandleCollision(w *ecs.World, aID, bID static.EntityID) {
	a := ecs.Get(w, aID, ecs.Collider)
	b := ecs.Get(w, bID, ecs.Collider)

	if a == nil || b == nil {
		return
	}

	if (a.Mask & static.CollisionGroup(b.Layer)) == 0 {
		return
	}
	if (b.Mask & static.CollisionGroup(a.Layer)) == 0 {
		return
	}

	if handlers, ok := collisionHandlers[static.CollisionGroup(a.Layer)]; ok {
		if handler, ok := handlers[static.CollisionGroup(b.Layer)]; ok {
			handler(w, aID, bID)
			return
		}
	}
}
