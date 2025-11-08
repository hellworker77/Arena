package ecs_systems

import "game-server/internal/application/ecs"

type ItemPickupSystem struct {
	ItemPickupRange float32
}

func (s ItemPickupSystem) Run(w *ecs.World, dt float32) {
	it := w.Query(ecs.CPlayerTag | ecs.CPos | ecs.CInventory).Iter()

	for it.Next() {
		playerID := it.EntityID()
		playerPos := *it.Position()
		inv := it.Inventory()

		nearItems := w.Query(ecs.CItemTag|ecs.CPos|ecs.CWorldItem).RangeIter(playerPos, s.ItemPickupRange)

		for nearItems.Next() {
			itemID := nearItems.EntityID()
			item := nearItems.WorldItem()

			// AddToInventory(inv, item.BaseID, item.Count)

			w.RemoveEntity(itemID)
		}
	}
}

func (ItemPickupSystem) Reads() ecs.Signature {
	return 0
}

func (ItemPickupSystem) Writes() ecs.Signature {
	return 0
}
