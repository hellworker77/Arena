package runtime

import (
	"game-server/internal/ecs/ecs_signatures/static"
)

type InterestState struct {
	Visible map[static.EntityID]struct{}
}
