package runtime

import "game-server/internal/application/ecs/ecs_signatures/static"

type ProjectileState struct {
	OwnerID  static.EntityID
	TargetID static.EntityID
	SpawnPos Position
	Elapsed  float32
}
