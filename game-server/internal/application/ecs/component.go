package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/ecs/ecs_signatures/tag"
)

type Component[T any] struct {
	Mask  Signature
	Slice func(a *Archetype) *[]T
}

var (
	Position = Component[runtime.Position]{
		Mask:  CPos,
		Slice: func(a *Archetype) *[]runtime.Position { return &a.Positions },
	}
	Velocity = Component[runtime.Velocity]{
		Mask:  CVel,
		Slice: func(a *Archetype) *[]runtime.Velocity { return &a.Velocities },
	}
	Health = Component[runtime.Health]{
		Mask:  CHealth,
		Slice: func(a *Archetype) *[]runtime.Health { return &a.Healths },
	}
	Experience = Component[runtime.Experience]{
		Mask:  CExperience,
		Slice: func(a *Archetype) *[]runtime.Experience { return &a.Experiences },
	}
	AttackCooldown = Component[runtime.AttackCooldown]{
		Mask:  CAttackCooldown,
		Slice: func(a *Archetype) *[]runtime.AttackCooldown { return &a.AttackCooldowns },
	}
	Target = Component[runtime.Target]{
		Mask:  CTarget,
		Slice: func(a *Archetype) *[]runtime.Target { return &a.Targets },
	}
	Lifespan = Component[runtime.Lifespan]{
		Mask:  CLifespan,
		Slice: func(a *Archetype) *[]runtime.Lifespan { return &a.Lifespans },
	}
	ProjectileState = Component[runtime.ProjectileState]{
		Mask:  CProjectileState,
		Slice: func(a *Archetype) *[]runtime.ProjectileState { return &a.ProjectileStates },
	}
	Inventory = Component[runtime.Inventory]{
		Mask:  CInventory,
		Slice: func(a *Archetype) *[]runtime.Inventory { return &a.Inventories },
	}
	WorldItem = Component[runtime.WorldItem]{
		Mask:  CWorldItem,
		Slice: func(a *Archetype) *[]runtime.WorldItem { return &a.WorldItems },
	}
	InterestState = Component[runtime.InterestState]{
		Mask:  CInterestState,
		Slice: func(a *Archetype) *[]runtime.InterestState { return &a.InterestStates },
	}
	Collider = Component[static.Collider]{
		Mask:  CCollider,
		Slice: func(a *Archetype) *[]static.Collider { return &a.Colliders },
	}
	CombatAttributes = Component[static.CombatAttributes]{
		Mask:  CCombatAttrs,
		Slice: func(a *Archetype) *[]static.CombatAttributes { return &a.CombatAttrs },
	}
	MovementAttributes = Component[static.MovementAttributes]{
		Mask:  CMovementAttrs,
		Slice: func(a *Archetype) *[]static.MovementAttributes { return &a.MovementAttrs },
	}
	EnemyPreset = Component[static.EnemyPreset]{
		Mask:  CEnemyPreset,
		Slice: func(a *Archetype) *[]static.EnemyPreset { return &a.EnemyPresets },
	}
	ProjectilePreset = Component[static.ProjectilePreset]{
		Mask:  CProjectilePreset,
		Slice: func(a *Archetype) *[]static.ProjectilePreset { return &a.ProjectilePresets },
	}
	PlayerTag = Component[tag.PlayerTag]{
		Mask:  CPlayerTag,
		Slice: func(a *Archetype) *[]tag.PlayerTag { return &a.PlayerTags },
	}
	EnemyTag = Component[tag.EnemyTag]{
		Mask:  CEnemyTag,
		Slice: func(a *Archetype) *[]tag.EnemyTag { return &a.EnemyTags },
	}
	NpcTag = Component[tag.NpcTag]{
		Mask:  CNpcTag,
		Slice: func(a *Archetype) *[]tag.NpcTag { return &a.NpcTags },
	}
	ProjectileTag = Component[tag.ProjectileTag]{
		Mask:  CProjectileTag,
		Slice: func(a *Archetype) *[]tag.ProjectileTag { return &a.ProjectileTags },
	}
	ItemTag = Component[tag.ItemTag]{
		Mask:  CItemTag,
		Slice: func(a *Archetype) *[]tag.ItemTag { return &a.ItemTags },
	}
)
