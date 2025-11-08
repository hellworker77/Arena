package ecs_systems

import (
	"game-server/internal/application/ecs"
)

type MovementSystem struct{}

func (MovementSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CPos | ecs.CVel).Iter()

	for it.Next() {
		pos := it.Position()
		vel := it.Velocity()

		pos.Add(vel.Mul(dt))
	}
}

func (MovementSystem) Reads() ecs.Signature {
	return ecs.CPos | ecs.CVel
}

func (MovementSystem) Writes() ecs.Signature {
	return ecs.CPos
}
