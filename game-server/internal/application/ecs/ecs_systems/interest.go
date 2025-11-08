package ecs_systems

import (
	"game-server/internal/application/ecs"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

type InterestSystem struct {
	VisibilityRange float32
}

func (s InterestSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CPlayerTag | ecs.CPos | ecs.CInterestState).Iter()

	for it.Next() {
		playerID := it.EntityID()
		playerPos := *it.Position()
		interestState := it.InterestState()

		nearby := w.Grid.Query(1, playerPos, s.VisibilityRange)

		newVisible := make(map[static.EntityID]struct{}, len(nearby))
		for _, eID := range nearby {
			if eID != playerID {
				newVisible[eID] = struct{}{}
			}
		}

		appeared := make([]static.EntityID, 0)
		for eid := range newVisible {
			if _, ok := interestState.Visible[eid]; !ok {
				appeared = append(appeared, eid)
			}
		}

		disappeared := make([]static.EntityID, 0)
		for eid := range interestState.Visible {
			if _, ok := newVisible[eid]; !ok {
				disappeared = append(disappeared, eid)
			}
		}

		interestState.Visible = newVisible

		if len(appeared) == 0 && len(disappeared) == 0 {
			continue
		}

		//net.SendVisibleDelta(playerID, appeared, disappeared)
	}
}

func (InterestSystem) Reads() ecs.Signature {
	return 0
}

func (InterestSystem) Writes() ecs.Signature {
	return 0
}
