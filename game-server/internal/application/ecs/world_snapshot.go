package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
	"game-server/internal/application/ecs/ecs_types"
)

type WorldSnapshot struct {
	Entities []EntitySnapshot
}

type EntitySnapshot struct {
	ID    static.EntityID
	Pos   runtime.Position
	Stats ecs_types.Stats
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

		if a.Signature&CStats != 0 {
			e.Stats = a.Stats[i]
		}

		snapshot.Entities = append(snapshot.Entities, e)
	}

	return snapshot
}
