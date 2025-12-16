package ecs

import (
	runtime2 "game-server/internal/ecs/ecs_signatures/runtime"
	static2 "game-server/internal/ecs/ecs_signatures/static"
	tag2 "game-server/internal/ecs/ecs_signatures/tag"
)

type Archetype struct {
	Signature Signature
	Count     int

	// static components
	Colliders         []static2.Collider
	CombatAttrs       []static2.CombatAttributes
	EnemyPresets      []static2.EnemyPreset
	EntityIDs         []static2.EntityID
	MovementAttrs     []static2.MovementAttributes
	ProjectilePresets []static2.ProjectilePreset
	// runtime components
	AttackCooldowns  []runtime2.AttackCooldown
	Experiences      []runtime2.Experience
	Healths          []runtime2.Health
	Lifespans        []runtime2.Lifespan
	Positions        []runtime2.Position
	ProjectileStates []runtime2.ProjectileState
	Targets          []runtime2.Target
	Velocities       []runtime2.Velocity
	Inventories      []runtime2.Inventory
	WorldItems       []runtime2.WorldItem
	InterestStates   []runtime2.InterestState
	// tag components
	PlayerTags     []tag2.PlayerTag
	EnemyTags      []tag2.EnemyTag
	NpcTags        []tag2.NpcTag
	ProjectileTags []tag2.ProjectileTag
	ItemTags       []tag2.ItemTag
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

type componentMeta struct {
	mask  Signature
	slice interface{}
}

func (a *Archetype) getComponentMetas() []componentMeta {
	return []componentMeta{
		{mask: CPos, slice: &a.Positions},
		{mask: CVel, slice: &a.Velocities},
		{mask: CHealth, slice: &a.Healths},
		{mask: CAttackCooldown, slice: &a.AttackCooldowns},
		{mask: CLifespan, slice: &a.Lifespans},
		{mask: CTarget, slice: &a.Targets},
		{mask: CExperience, slice: &a.Experiences},
		{mask: CProjectileState, slice: &a.ProjectileStates},
		{mask: CInventory, slice: &a.Inventories},
		{mask: CWorldItem, slice: &a.WorldItems},
		{mask: CInterestState, slice: &a.InterestStates},
		{mask: CCollider, slice: &a.Colliders},
		{mask: CCombatAttrs, slice: &a.CombatAttrs},
		{mask: CMovementAttrs, slice: &a.MovementAttrs},
		{mask: CEnemyPreset, slice: &a.EnemyPresets},
		{mask: CProjectilePreset, slice: &a.ProjectilePresets},
		{mask: CPlayerTag, slice: &a.PlayerTags},
		{mask: CEnemyTag, slice: &a.EnemyTags},
		{mask: CNpcTag, slice: &a.NpcTags},
		{mask: CProjectileTag, slice: &a.ProjectileTags},
		{mask: CItemTag, slice: &a.ItemTags},
	}
}

func growSlice(slicePtr interface{}, count int) {
	switch s := slicePtr.(type) {
	case *[]runtime2.Position:
		*s = grow(*s, count)
	case *[]runtime2.Velocity:
		*s = grow(*s, count)
	case *[]runtime2.Health:
		*s = grow(*s, count)
	case *[]runtime2.AttackCooldown:
		*s = grow(*s, count)
	case *[]runtime2.Lifespan:
		*s = grow(*s, count)
	case *[]runtime2.Target:
		*s = grow(*s, count)
	case *[]runtime2.Experience:
		*s = grow(*s, count)
	case *[]runtime2.ProjectileState:
		*s = grow(*s, count)
	case *[]runtime2.Inventory:
		*s = grow(*s, count)
	case *[]runtime2.WorldItem:
		*s = grow(*s, count)
	case *[]runtime2.InterestState:
		*s = grow(*s, count)
	case *[]static2.Collider:
		*s = grow(*s, count)
	case *[]static2.CombatAttributes:
		*s = grow(*s, count)
	case *[]static2.MovementAttributes:
		*s = grow(*s, count)
	case *[]static2.EnemyPreset:
		*s = grow(*s, count)
	case *[]static2.ProjectilePreset:
		*s = grow(*s, count)
	case *[]tag2.PlayerTag:
		*s = grow(*s, count)
	case *[]tag2.EnemyTag:
		*s = grow(*s, count)
	case *[]tag2.NpcTag:
		*s = grow(*s, count)
	case *[]tag2.ProjectileTag:
		*s = grow(*s, count)
	case *[]tag2.ItemTag:
		*s = grow(*s, count)
	default:
		panic("unknown slice type")
	}
}

func swapSlice(slicePtr any, i, j int) {
	switch s := slicePtr.(type) {
	case *[]runtime2.Position:
		swap(*s, i, j)
	case *[]runtime2.Velocity:
		swap(*s, i, j)
	case *[]runtime2.Health:
		swap(*s, i, j)
	case *[]runtime2.AttackCooldown:
		swap(*s, i, j)
	case *[]runtime2.Lifespan:
		swap(*s, i, j)
	case *[]runtime2.Target:
		swap(*s, i, j)
	case *[]runtime2.Experience:
		swap(*s, i, j)
	case *[]runtime2.ProjectileState:
		swap(*s, i, j)
	case *[]runtime2.Inventory:
		swap(*s, i, j)
	case *[]runtime2.WorldItem:
		swap(*s, i, j)
	case *[]runtime2.InterestState:
		swap(*s, i, j)
	case *[]static2.Collider:
		swap(*s, i, j)
	case *[]static2.CombatAttributes:
		swap(*s, i, j)
	case *[]static2.MovementAttributes:
		swap(*s, i, j)
	case *[]static2.EnemyPreset:
		swap(*s, i, j)
	case *[]static2.ProjectilePreset:
		swap(*s, i, j)
	case *[]tag2.PlayerTag:
		swap(*s, i, j)
	case *[]tag2.EnemyTag:
		swap(*s, i, j)
	case *[]tag2.NpcTag:
		swap(*s, i, j)
	case *[]tag2.ProjectileTag:
		swap(*s, i, j)
	case *[]tag2.ItemTag:
		swap(*s, i, j)
	default:
		panic("unknown slice type")
	}
}

func (a *Archetype) InsertEmpty(eID static2.EntityID) int {
	i := a.Count
	a.Count++

	a.EntityIDs = grow(a.EntityIDs, a.Count)
	a.EntityIDs[i] = eID

	for _, m := range a.getComponentMetas() {
		if a.Signature&m.mask != 0 {
			growSlice(m.slice, a.Count)
		}
	}

	return i
}

func (a *Archetype) Remove(index int) static2.EntityID {
	last := a.Count - 1
	movedID := a.EntityIDs[last]

	swap(a.EntityIDs, index, last)

	for _, m := range a.getComponentMetas() {
		if a.Signature&m.mask != 0 {
			swapSlice(m.slice, index, last)
		}
	}

	a.Count--
	return movedID
}
