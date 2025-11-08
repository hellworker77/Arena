package ecs_data

import "game-server/internal/application/ecs/ecs_types"

var EnemyPresets = map[string]ecs_types.EnemyCfg{
	"engine.core.fallen": {
		Name:        "Fallen",
		Level:       1,
		BaseExpGain: 10,
		DefaultsStats: ecs_types.Stats{
			Health:          30,
			MaxHealth:       30,
			AttackSpeed:     1.5,
			AttackDamage:    3,
			MovementSpeed:   1.0,
			Defense:         1,
			IsMelee:         true,
			AttackRange:     1.0,
			AttackRating:    5,
			ProjectileSpeed: 0,
		},
		ProjectileCfg: nil,
	},
	"engine.core.fallen_shaman": {
		Name:        "Fallen Shaman",
		Level:       2,
		BaseExpGain: 20,
		DefaultsStats: ecs_types.Stats{
			Health:          30,
			MaxHealth:       30,
			AttackSpeed:     1.5,
			AttackDamage:    11,
			MovementSpeed:   0.7,
			Defense:         1,
			IsMelee:         false,
			AttackRange:     50.0,
			AttackRating:    20,
			ProjectileSpeed: 0,
		},
		ProjectileCfg: &ecs_types.ProjectileCfg{
			ProjectileID: "engine.core.fallen_fire_ball",
			Trajectory:   ecs_types.TrajectoryLinear,
			Params: map[string]float32{
				"radius":       2.0,
				"angularSpeed": 0.0,
				"explosion":    0.0,
			},
		},
	},
}
