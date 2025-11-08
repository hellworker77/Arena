package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/ecs/ecs_signatures/tag"
)

type Archetype struct {
	Signature Signature
	Count     int

	// static components
	Colliders         []static.Collider
	CombatAttrs       []static.CombatAttributes
	EnemyPresets      []static.EnemyPreset
	EntityIDs         []static.EntityID
	MovementAttrs     []static.MovementAttributes
	ProjectilePresets []static.ProjectilePreset
	// runtime components
	AttackCooldowns  []runtime.AttackCooldown
	Experiences      []runtime.Experience
	Healths          []runtime.Health
	Lifespans        []runtime.Lifespan
	Positions        []runtime.Position
	ProjectileStates []runtime.ProjectileState
	Targets          []runtime.Target
	Velocities       []runtime.Velocity
	Inventories      []runtime.Inventory
	WorldItems       []runtime.WorldItem
	InterestStates   []runtime.InterestState
	// tag components
	PlayerTags     []tag.PlayerTag
	EnemyTags      []tag.EnemyTag
	NpcTags        []tag.NpcTag
	ProjectileTags []tag.ProjectileTag
	ItemTags       []tag.ItemTag
}

func grow[T any](slice []T, count int) []T {
	if len(slice) >= count {
		return slice
	}
	needed := count - len(slice)
	ext := make([]T, needed)
	return append(slice, ext...)
}

func swap[T any](slice []T, i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (a *Archetype) InsertEmpty(eID static.EntityID) int {
	i := a.Count
	a.Count++

	a.EntityIDs = grow(a.EntityIDs, a.Count)
	a.EntityIDs[i] = eID

	// runtime
	if a.Signature&CPos != 0 {
		a.Positions = grow(a.Positions, a.Count)
	}
	if a.Signature&CVel != 0 {
		a.Velocities = grow(a.Velocities, a.Count)
	}
	if a.Signature&CHealth != 0 {
		a.Healths = grow(a.Healths, a.Count)
	}
	if a.Signature&CAttackCooldown != 0 {
		a.AttackCooldowns = grow(a.AttackCooldowns, a.Count)
	}
	if a.Signature&CLifespan != 0 {
		a.Lifespans = grow(a.Lifespans, a.Count)
	}
	if a.Signature&CTarget != 0 {
		a.Targets = grow(a.Targets, a.Count)
	}
	if a.Signature&CExperience != 0 {
		a.Experiences = grow(a.Experiences, a.Count)
	}
	if a.Signature&CProjectileState != 0 {
		a.ProjectileStates = grow(a.ProjectileStates, a.Count)
	}
	if a.Signature&CInventory != 0 {
		a.Inventories = grow(a.Inventories, a.Count)
	}
	if a.Signature&CWorldItem != 0 {
		a.WorldItems = grow(a.WorldItems, a.Count)
	}
	if a.Signature&CInterestState != 0 {
		a.InterestStates = grow(a.InterestStates, a.Count)
	}
	// static
	if a.Signature&CCollider != 0 {
		a.Colliders = grow(a.Colliders, a.Count)
	}
	if a.Signature&CCombatAttrs != 0 {
		a.CombatAttrs = grow(a.CombatAttrs, a.Count)
	}
	if a.Signature&CMovementAttrs != 0 {
		a.MovementAttrs = grow(a.MovementAttrs, a.Count)
	}
	if a.Signature&CEnemyPreset != 0 {
		a.EnemyPresets = grow(a.EnemyPresets, a.Count)
	}
	if a.Signature&CProjectilePreset != 0 {
		a.ProjectilePresets = grow(a.ProjectilePresets, a.Count)
	}

	// tags
	if a.Signature&CPlayerTag != 0 {
		a.PlayerTags = grow(a.PlayerTags, a.Count)
	}
	if a.Signature&CEnemyTag != 0 {
		a.EnemyTags = grow(a.EnemyTags, a.Count)
	}
	if a.Signature&CNpcTag != 0 {
		a.NpcTags = grow(a.NpcTags, a.Count)
	}
	if a.Signature&CProjectileTag != 0 {
		a.ProjectileTags = grow(a.ProjectileTags, a.Count)
	}
	if a.Signature&CItemTag != 0 {
		a.ItemTags = grow(a.ItemTags, a.Count)
	}

	return i
}

func (a *Archetype) Remove(index int) static.EntityID {
	last := a.Count - 1
	movedID := a.EntityIDs[last]

	swap(a.EntityIDs, index, last)
	// runtime
	if a.Signature&CPos != 0 {
		swap(a.Positions, index, last)
	}
	if a.Signature&CVel != 0 {
		swap(a.Velocities, index, last)
	}
	if a.Signature&CHealth != 0 {
		swap(a.Healths, index, last)
	}
	if a.Signature&CAttackCooldown != 0 {
		swap(a.AttackCooldowns, index, last)
	}
	if a.Signature&CLifespan != 0 {
		swap(a.Lifespans, index, last)
	}
	if a.Signature&CTarget != 0 {
		swap(a.Targets, index, last)
	}
	if a.Signature&CExperience != 0 {
		swap(a.Experiences, index, last)
	}
	if a.Signature&CProjectileState != 0 {
		swap(a.ProjectileStates, index, last)
	}
	if a.Signature&CInventory != 0 {
		swap(a.Inventories, index, last)
	}
	if a.Signature&CWorldItem != 0 {
		swap(a.WorldItems, index, last)
	}
	if a.Signature&CInterestState != 0 {
		swap(a.InterestStates, index, last)
	}
	// static
	if a.Signature&CCollider != 0 {
		swap(a.Colliders, index, last)
	}
	if a.Signature&CCombatAttrs != 0 {
		swap(a.CombatAttrs, index, last)
	}
	if a.Signature&CMovementAttrs != 0 {
		swap(a.MovementAttrs, index, last)
	}
	if a.Signature&CEnemyPreset != 0 {
		swap(a.EnemyPresets, index, last)
	}
	if a.Signature&CProjectilePreset != 0 {
		swap(a.ProjectilePresets, index, last)
	}
	// tags
	if a.Signature&CPlayerTag != 0 {
		swap(a.PlayerTags, index, last)
	}
	if a.Signature&CEnemyTag != 0 {
		swap(a.EnemyTags, index, last)
	}
	if a.Signature&CNpcTag != 0 {
		swap(a.NpcTags, index, last)
	}
	if a.Signature&CProjectileTag != 0 {
		swap(a.ProjectileTags, index, last)
	}
	if a.Signature&CItemTag != 0 {
		swap(a.ItemTags, index, last)
	}

	a.Count--
	return movedID
}
