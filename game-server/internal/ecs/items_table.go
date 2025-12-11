package ecs

import "game-server/internal/ecs/ecs_signatures/runtime"

var ItemDB = map[int]runtime.Item{}

var nextItemID = 1

func GenerateItem(item runtime.Item) int {
	item.ID = nextItemID
	nextItemID++

	ItemDB[item.ID] = item
	return item.ID
}
