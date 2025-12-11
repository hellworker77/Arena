package udp

import (
	ecs2 "game-server/internal/ecs"
	runtime2 "game-server/internal/ecs/ecs_signatures/runtime"
	"game-server/internal/ecs/ecs_signatures/static"
)

func (s *Server) createPlayerEntity(addrStr string) static.EntityID {
	e := s.world.CreateEntity(ecs2.CPlayerTag | ecs2.CPos | ecs2.CVel | ecs2.CHealth)
	ecs2.Set(s.world, e, ecs2.Position, runtime2.Position{})
	ecs2.Set(s.world, e, ecs2.Velocity, runtime2.Velocity{})
	ecs2.Set(s.world, e, ecs2.Health, runtime2.Health{Current: 100})
	return e
}
