package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
)

type MovementSystem struct{}

func (MovementSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CPos | ecs2.CVel).Iter()

	for it.Next() {
		pos := it.Position()
		vel := it.Velocity()

		pos.Add(vel.Mul(dt))
	}
}

func (MovementSystem) Reads() ecs2.Signature {
	return ecs2.CPos | ecs2.CVel
}

func (MovementSystem) Writes() ecs2.Signature {
	return ecs2.CPos
}
