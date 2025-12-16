package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
)

type ItemPickupSystem struct {
	ItemPickupRange float32
}

func (s ItemPickupSystem) Run(w *ecs2.World, dt float32) {
	it := w.Query(ecs2.CPlayerTag | ecs2.CPos | ecs2.CInventory).Iter()

	for it.Next() {
		playerPos := *it.Position()

		nearItems := w.Query(ecs2.CItemTag|ecs2.CPos|ecs2.CWorldItem).RangeIter(playerPos, s.ItemPickupRange)

		for nearItems.Next() {
			itemID := nearItems.EntityID()

			// AddToInventory(inv, item.BaseID, item.Count)

			w.RemoveEntity(itemID)
		}
	}
}

func (ItemPickupSystem) Reads() ecs2.Signature {
	return 0
}

func (ItemPickupSystem) Writes() ecs2.Signature {
	return 0
}
