package ecs

type Signature uint64

const (
	CPos             Signature = 1 << iota // runtime.Position
	CVel                                   // runtime.Velocity
	CHealth                                // runtime.Health
	CExperience                            // runtime.Experience
	CAttackCooldown                        // runtime.AttackCooldown
	CTarget                                // runtime.Target
	CLifespan                              // runtime.Lifespan
	CProjectileState                       // runtime.ProjectileState
	CInventory                             // runtime.Inventory
	CWorldItem                             // runtime.WorldItem
	CInterestState                         // runtime.InterestState
)

const (
	CCollider         Signature = 1 << (iota + 16) // static.Collider
	CCombatAttrs                                   // static.CombatAttributes
	CMovementAttrs                                 // static.MovementAttributes
	CEnemyPreset                                   // static.EnemyPreset
	CProjectilePreset                              // static.ProjectilePreset
)

const (
	CPlayerTag Signature = 1 << (iota + 32)
	CEnemyTag
	CNpcTag
	CProjectileTag
	CItemTag
)
