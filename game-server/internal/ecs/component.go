package ecs

import (
	runtime2 "game-server/internal/ecs/ecs_signatures/runtime"
	static2 "game-server/internal/ecs/ecs_signatures/static"
	tag2 "game-server/internal/ecs/ecs_signatures/tag"
)

type Component[T any] struct {
	Mask  Signature
	Slice func(a *Archetype) *[]T
}

var (
	Position = Component[runtime2.Position]{
		Mask:  CPos,
		Slice: func(a *Archetype) *[]runtime2.Position { return &a.Positions },
	}
	Velocity = Component[runtime2.Velocity]{
		Mask:  CVel,
		Slice: func(a *Archetype) *[]runtime2.Velocity { return &a.Velocities },
	}
	Health = Component[runtime2.Health]{
		Mask:  CHealth,
		Slice: func(a *Archetype) *[]runtime2.Health { return &a.Healths },
	}
	Experience = Component[runtime2.Experience]{
		Mask:  CExperience,
		Slice: func(a *Archetype) *[]runtime2.Experience { return &a.Experiences },
	}
	AttackCooldown = Component[runtime2.AttackCooldown]{
		Mask:  CAttackCooldown,
		Slice: func(a *Archetype) *[]runtime2.AttackCooldown { return &a.AttackCooldowns },
	}
	Target = Component[runtime2.Target]{
		Mask:  CTarget,
		Slice: func(a *Archetype) *[]runtime2.Target { return &a.Targets },
	}
	Lifespan = Component[runtime2.Lifespan]{
		Mask:  CLifespan,
		Slice: func(a *Archetype) *[]runtime2.Lifespan { return &a.Lifespans },
	}
	ProjectileState = Component[runtime2.ProjectileState]{
		Mask:  CProjectileState,
		Slice: func(a *Archetype) *[]runtime2.ProjectileState { return &a.ProjectileStates },
	}
	Inventory = Component[runtime2.Inventory]{
		Mask:  CInventory,
		Slice: func(a *Archetype) *[]runtime2.Inventory { return &a.Inventories },
	}
	WorldItem = Component[runtime2.WorldItem]{
		Mask:  CWorldItem,
		Slice: func(a *Archetype) *[]runtime2.WorldItem { return &a.WorldItems },
	}
	InterestState = Component[runtime2.InterestState]{
		Mask:  CInterestState,
		Slice: func(a *Archetype) *[]runtime2.InterestState { return &a.InterestStates },
	}
	Collider = Component[static2.Collider]{
		Mask:  CCollider,
		Slice: func(a *Archetype) *[]static2.Collider { return &a.Colliders },
	}
	CombatAttributes = Component[static2.CombatAttributes]{
		Mask:  CCombatAttrs,
		Slice: func(a *Archetype) *[]static2.CombatAttributes { return &a.CombatAttrs },
	}
	MovementAttributes = Component[static2.MovementAttributes]{
		Mask:  CMovementAttrs,
		Slice: func(a *Archetype) *[]static2.MovementAttributes { return &a.MovementAttrs },
	}
	EnemyPreset = Component[static2.EnemyPreset]{
		Mask:  CEnemyPreset,
		Slice: func(a *Archetype) *[]static2.EnemyPreset { return &a.EnemyPresets },
	}
	ProjectilePreset = Component[static2.ProjectilePreset]{
		Mask:  CProjectilePreset,
		Slice: func(a *Archetype) *[]static2.ProjectilePreset { return &a.ProjectilePresets },
	}
	PlayerTag = Component[tag2.PlayerTag]{
		Mask:  CPlayerTag,
		Slice: func(a *Archetype) *[]tag2.PlayerTag { return &a.PlayerTags },
	}
	EnemyTag = Component[tag2.EnemyTag]{
		Mask:  CEnemyTag,
		Slice: func(a *Archetype) *[]tag2.EnemyTag { return &a.EnemyTags },
	}
	NpcTag = Component[tag2.NpcTag]{
		Mask:  CNpcTag,
		Slice: func(a *Archetype) *[]tag2.NpcTag { return &a.NpcTags },
	}
	ProjectileTag = Component[tag2.ProjectileTag]{
		Mask:  CProjectileTag,
		Slice: func(a *Archetype) *[]tag2.ProjectileTag { return &a.ProjectileTags },
	}
	ItemTag = Component[tag2.ItemTag]{
		Mask:  CItemTag,
		Slice: func(a *Archetype) *[]tag2.ItemTag { return &a.ItemTags },
	}
)
