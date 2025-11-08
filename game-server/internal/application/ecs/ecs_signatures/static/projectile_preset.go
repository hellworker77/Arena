package static

import "game-server/internal/application/ecs/ecs_types"

type ProjectilePreset struct {
	ProjectileID string // Unique identifier for display in client UI
	Speed        float32
	Trajectory   ecs_types.TrajectoryType
	Params       map[string]float32
	IsHoming     bool
	Pierce       int // -1 for infinite pierce
}
