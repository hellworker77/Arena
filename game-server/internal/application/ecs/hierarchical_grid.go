package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

type cellKey struct {
	X, Y int32
}

type GridLevel struct {
	CellSize float32
	Cells    map[cellKey][]static.EntityID
}

type HierarchicalGrid struct {
	Levels []GridLevel
}

func NewHierarchicalGrid(levelSizes ...float32) *HierarchicalGrid {
	levels := make([]GridLevel, len(levelSizes))

	for i, size := range levelSizes {
		levels[i] = GridLevel{
			CellSize: size,
			Cells:    make(map[cellKey][]static.EntityID),
		}
	}

	return &HierarchicalGrid{Levels: levels}
}

func (hl *HierarchicalGrid) cellKey(pos runtime.Position, cellSIze float32) cellKey {
	return cellKey{
		X: int32(pos.X / cellSIze),
		Y: int32(pos.Y / cellSIze),
	}
}

func (hl *HierarchicalGrid) Insert(entityID static.EntityID, pos runtime.Position) {
	for i := range hl.Levels {
		key := hl.cellKey(pos, hl.Levels[i].CellSize)
		hl.Levels[i].Cells[key] = append(hl.Levels[i].Cells[key], entityID)
	}
}

func (hl *HierarchicalGrid) Remove(entityID static.EntityID, pos runtime.Position) {
	for i := range hl.Levels {
		key := hl.cellKey(pos, hl.Levels[i].CellSize)
		list := hl.Levels[i].Cells[key]

		for j := range list {
			if list[j] == entityID {
				list[j] = list[len(list)-1]
				hl.Levels[i].Cells[key] = list[:len(list)-1]
				break
			}
		}
	}
}

func (hl *HierarchicalGrid) Move(eID static.EntityID, oldPos, newPos runtime.Position) {
	for i := range hl.Levels {
		oldKey := hl.cellKey(oldPos, hl.Levels[i].CellSize)
		newKey := hl.cellKey(newPos, hl.Levels[i].CellSize)

		if oldKey != newKey {
			hl.Remove(eID, oldPos)
			hl.Insert(eID, newPos)
		}
	}
}

func (hl *HierarchicalGrid) Query(level int, pos runtime.Position, rangeDist float32) []static.EntityID {
	if level < 0 || level >= len(hl.Levels) {
		return nil
	}

	gl := &hl.Levels[level]
	cellRange := int32(rangeDist / gl.CellSize)
	centerKey := hl.cellKey(pos, gl.CellSize)

	result := make([]static.EntityID, 0, 64)

	for dx := -cellRange; dx <= cellRange; dx++ {
		for dy := -cellRange; dy <= cellRange; dy++ {
			key := cellKey{X: centerKey.X + dx, Y: centerKey.Y + dy}

			if entities, exists := gl.Cells[key]; exists {
				result = append(result, entities...)
			}
		}
	}

	return result
}
