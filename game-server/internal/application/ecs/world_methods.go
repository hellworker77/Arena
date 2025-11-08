package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/ecs/ecs_signatures/tag"
)

func (w *World) WritePosition(eID static.EntityID, v runtime.Position) {
	rec := w.entities[eID]
	writeComponent(w, eID, CPos, &rec.Archetype.Positions, v)
}

func (w *World) GetPosition(eID static.EntityID) *runtime.Position {
	rec := w.entities[eID]
	return getComponent(w, eID, CPos, &rec.Archetype.Positions)
}

func (w *World) RemovePosition(eID static.EntityID) {
	removeComponent(w, eID, CPos)
}

func (w *World) WriteVelocity(eID static.EntityID, v runtime.Velocity) {
	rec := w.entities[eID]
	writeComponent(w, eID, CVel, &rec.Archetype.Velocities, v)
}

func (w *World) GetVelocity(eID static.EntityID) *runtime.Velocity {
	rec := w.entities[eID]
	return getComponent(w, eID, CVel, &rec.Archetype.Velocities)
}
func (w *World) RemoveVelocity(eID static.EntityID) {
	removeComponent(w, eID, CVel)
}

func (w *World) WriteHealth(eID static.EntityID, v runtime.Health) {
	rec := w.entities[eID]
	writeComponent(w, eID, CHealth, &rec.Archetype.Healths, v)
}

func (w *World) GetHealth(eID static.EntityID) *runtime.Health {
	rec := w.entities[eID]
	return getComponent(w, eID, CHealth, &rec.Archetype.Healths)
}
func (w *World) RemoveHealth(eID static.EntityID) {
	removeComponent(w, eID, CHealth)
}

func (w *World) WriteExperience(eID static.EntityID, v runtime.Experience) {
	rec := w.entities[eID]
	writeComponent(w, eID, CExperience, &rec.Archetype.Experiences, v)
}

func (w *World) GetExperience(eID static.EntityID) *runtime.Experience {
	rec := w.entities[eID]
	return getComponent(w, eID, CExperience, &rec.Archetype.Experiences)
}
func (w *World) RemoveExperience(eID static.EntityID) {
	removeComponent(w, eID, CExperience)
}

func (w *World) WriteAttackCooldown(eID static.EntityID, v runtime.AttackCooldown) {
	rec := w.entities[eID]
	writeComponent(w, eID, CAttackCooldown, &rec.Archetype.AttackCooldowns, v)
}

func (w *World) GetAttackCooldown(eID static.EntityID) *runtime.AttackCooldown {
	rec := w.entities[eID]
	return getComponent(w, eID, CAttackCooldown, &rec.Archetype.AttackCooldowns)
}
func (w *World) RemoveAttackCooldown(eID static.EntityID) {
	removeComponent(w, eID, CAttackCooldown)
}

func (w *World) WriteTarget(eID static.EntityID, v runtime.Target) {
	rec := w.entities[eID]
	writeComponent(w, eID, CTarget, &rec.Archetype.Targets, v)
}

func (w *World) GetTarget(eID static.EntityID) *runtime.Target {
	rec := w.entities[eID]
	return getComponent(w, eID, CTarget, &rec.Archetype.Targets)
}
func (w *World) RemoveTarget(eID static.EntityID) {
	removeComponent(w, eID, CTarget)
}

func (w *World) WriteLifespan(eID static.EntityID, v runtime.Lifespan) {
	rec := w.entities[eID]
	writeComponent(w, eID, CLifespan, &rec.Archetype.Lifespans, v)
}

func (w *World) GetLifespan(eID static.EntityID) *runtime.Lifespan {
	rec := w.entities[eID]
	return getComponent(w, eID, CLifespan, &rec.Archetype.Lifespans)
}
func (w *World) RemoveLifespan(eID static.EntityID) {
	removeComponent(w, eID, CLifespan)
}

func (w *World) WriteProjectileState(eID static.EntityID, v runtime.ProjectileState) {
	rec := w.entities[eID]
	writeComponent(w, eID, CProjectileState, &rec.Archetype.ProjectileStates, v)
}

func (w *World) GetProjectileState(eID static.EntityID) *runtime.ProjectileState {
	rec := w.entities[eID]
	return getComponent(w, eID, CProjectileState, &rec.Archetype.ProjectileStates)
}
func (w *World) RemoveProjectileState(eID static.EntityID) {
	removeComponent(w, eID, CProjectileState)
}

func (w *World) WriteCollider(eID static.EntityID, v static.Collider) {
	rec := w.entities[eID]
	writeComponent(w, eID, CCollider, &rec.Archetype.Colliders, v)
}

func (w *World) GetCollider(eID static.EntityID) *static.Collider {
	rec := w.entities[eID]
	return getComponent(w, eID, CCollider, &rec.Archetype.Colliders)
}
func (w *World) RemoveCollider(eID static.EntityID) {
	removeComponent(w, eID, CCollider)
}

func (w *World) WriteCombatAttributes(eID static.EntityID, v static.CombatAttributes) {
	rec := w.entities[eID]
	writeComponent(w, eID, CCombatAttrs, &rec.Archetype.CombatAttrs, v)
}

func (w *World) GetCombatAttributes(eID static.EntityID) *static.CombatAttributes {
	rec := w.entities[eID]
	return getComponent(w, eID, CCombatAttrs, &rec.Archetype.CombatAttrs)
}
func (w *World) RemoveCombatAttributes(eID static.EntityID) {
	removeComponent(w, eID, CCombatAttrs)
}

func (w *World) WriteMovementAttributes(eID static.EntityID, v static.MovementAttributes) {
	rec := w.entities[eID]
	writeComponent(w, eID, CMovementAttrs, &rec.Archetype.MovementAttrs, v)
}

func (w *World) GetMovementAttributes(eID static.EntityID) *static.MovementAttributes {
	rec := w.entities[eID]
	return getComponent(w, eID, CMovementAttrs, &rec.Archetype.MovementAttrs)
}
func (w *World) RemoveMovementAttributes(eID static.EntityID) {
	removeComponent(w, eID, CMovementAttrs)
}

func (w *World) WriteEnemyPreset(eID static.EntityID, v static.EnemyPreset) {
	rec := w.entities[eID]
	writeComponent(w, eID, CEnemyPreset, &rec.Archetype.EnemyPresets, v)
}

func (w *World) GetEnemyPreset(eID static.EntityID) *static.EnemyPreset {
	rec := w.entities[eID]
	return getComponent(w, eID, CEnemyPreset, &rec.Archetype.EnemyPresets)
}
func (w *World) RemoveEnemyPreset(eID static.EntityID) {
	removeComponent(w, eID, CEnemyPreset)
}

func (w *World) WriteProjectilePreset(eID static.EntityID, v static.ProjectilePreset) {
	rec := w.entities[eID]
	writeComponent(w, eID, CProjectilePreset, &rec.Archetype.ProjectilePresets, v)
}

func (w *World) GetProjectilePreset(eID static.EntityID) *static.ProjectilePreset {
	rec := w.entities[eID]
	return getComponent(w, eID, CProjectilePreset, &rec.Archetype.ProjectilePresets)
}
func (w *World) RemoveProjectilePreset(eID static.EntityID) {
	removeComponent(w, eID, CProjectilePreset)
}

func (w *World) WritePlayerTag(eID static.EntityID, v tag.PlayerTag) {
	rec := w.entities[eID]
	writeComponent(w, eID, CPlayerTag, &rec.Archetype.PlayerTags, v)
}

func (w *World) GetPlayerTag(eID static.EntityID) *tag.PlayerTag {
	rec := w.entities[eID]
	return getComponent(w, eID, CPlayerTag, &rec.Archetype.PlayerTags)
}
func (w *World) RemovePlayerTag(eID static.EntityID) {
	removeComponent(w, eID, CPlayerTag)
}

func (w *World) WriteEnemyTag(eID static.EntityID, v tag.EnemyTag) {
	rec := w.entities[eID]
	writeComponent(w, eID, CEnemyTag, &rec.Archetype.EnemyTags, v)
}

func (w *World) GetEnemyTag(eID static.EntityID) *tag.EnemyTag {
	rec := w.entities[eID]
	return getComponent(w, eID, CEnemyTag, &rec.Archetype.EnemyTags)
}
func (w *World) RemoveEnemyTag(eID static.EntityID) {
	removeComponent(w, eID, CEnemyTag)
}

func (w *World) WriteNpcTag(eID static.EntityID, v tag.NpcTag) {
	rec := w.entities[eID]
	writeComponent(w, eID, CNpcTag, &rec.Archetype.NpcTags, v)
}

func (w *World) GetNpcTag(eID static.EntityID) *tag.NpcTag {
	rec := w.entities[eID]
	return getComponent(w, eID, CNpcTag, &rec.Archetype.NpcTags)
}
func (w *World) RemoveNpcTag(eID static.EntityID) {
	removeComponent(w, eID, CNpcTag)
}

func (w *World) WriteProjectileTag(eID static.EntityID, v tag.ProjectileTag) {
	rec := w.entities[eID]
	writeComponent(w, eID, CProjectileTag, &rec.Archetype.ProjectileTags, v)
}

func (w *World) GetProjectileTag(eID static.EntityID) *tag.ProjectileTag {
	rec := w.entities[eID]
	return getComponent(w, eID, CProjectileTag, &rec.Archetype.ProjectileTags)
}
func (w *World) RemoveProjectileTag(eID static.EntityID) {
	removeComponent(w, eID, CProjectileTag)
}

func (w *World) WriteItemTag(eID static.EntityID, v tag.ItemTag) {
	rec := w.entities[eID]
	writeComponent(w, eID, CItemTag, &rec.Archetype.ItemTags, v)
}

func (w *World) GetItemTag(eID static.EntityID) *tag.ItemTag {
	rec := w.entities[eID]
	return getComponent(w, eID, CItemTag, &rec.Archetype.ItemTags)
}
func (w *World) RemoveItemTag(eID static.EntityID) {
	removeComponent(w, eID, CItemTag)
}
