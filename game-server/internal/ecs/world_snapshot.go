package ecs

import (
	"game-server/internal/ecs/ecs_signatures/runtime"
	"game-server/internal/ecs/ecs_signatures/static"
)

type WorldSnapshot struct {
	Entities []EntitySnapshot
}

type EntitySnapshot struct {
	ID  static.EntityID
	Pos runtime.Position
}

func (w *World) TakeSnapshot() WorldSnapshot {
	snapshot := WorldSnapshot{
		Entities: make([]EntitySnapshot, 0, len(w.entities)),
	}

	for id, rec := range w.entities {
		a := rec.Archetype
		i := rec.Index

		var e EntitySnapshot
		e.ID = id

		if a.Signature&CPos != 0 {
			e.Pos = a.Positions[i]
		}

		snapshot.Entities = append(snapshot.Entities, e)
	}

	return snapshot
}
