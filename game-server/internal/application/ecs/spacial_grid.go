package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

type gridCellKey struct {
	X, Y float32
}

type SpatialGrid struct {
	CellSize float32
	Cells    map[gridCellKey][]static.EntityID
}

func NewSpatialGrid(cellSize float32) *SpatialGrid {
	return &SpatialGrid{
		CellSize: cellSize,
		Cells:    make(map[gridCellKey][]static.EntityID),
	}
}

func (g *SpatialGrid) cell(pos runtime.Position) gridCellKey {
	return gridCellKey{
		X: pos.X / g.CellSize,
		Y: pos.Y / g.CellSize,
	}
}

func (g *SpatialGrid) Add(eID static.EntityID, pos runtime.Position) {
	key := g.cell(pos)
	g.Cells[key] = append(g.Cells[key], eID)
}

func (g *SpatialGrid) Remove(eID static.EntityID, pos runtime.Position) {
	key := g.cell(pos)
	entities := g.Cells[key]

	for i := range entities {
		if entities[i] == eID {
			entities[i] = entities[len(entities)-1]
			g.Cells[key] = entities[:len(entities)-1]
			return
		}
	}
}

func (g *SpatialGrid) Move(eID static.EntityID, oldPos, newPos runtime.Position) {
	oldKey := g.cell(oldPos)
	newKey := g.cell(newPos)
	if oldKey != newKey {
		g.Remove(eID, oldPos)
		g.Add(eID, newPos)
	}
}

func (g *SpatialGrid) QueryNearby(pos runtime.Position) []static.EntityID {
	key := g.cell(pos)

	res := make([]static.EntityID, 0, 16)

	for dx := int32(-1); dx <= 1; dx++ {
		for dy := int32(-1); dy <= 1; dy++ {
			k := gridCellKey{X: key.X + float32(dx), Y: key.Y + float32(dy)}
			if list, ok := g.Cells[k]; ok {
				res = append(res, list...)
			}
		}
	}

	return res
}

func (g *SpatialGrid) QueryInRange(pos runtime.Position, rangeDist float32) []static.EntityID {
	key := g.cell(pos)
	cellsRange := int32(rangeDist / g.CellSize)

	res := make([]static.EntityID, 0, 16)

	for dx := -cellsRange; dx <= cellsRange; dx++ {
		for dy := -cellsRange; dy <= cellsRange; dy++ {
			k := gridCellKey{X: key.X + float32(dx), Y: key.Y + float32(dy)}
			if list, ok := g.Cells[k]; ok {
				res = append(res, list...)
			}
		}
	}

	return res
}
