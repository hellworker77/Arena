package spatial

import (
	"math"
)

// Grid is a simple spatial hash grid.
//
// It is optimized for frequent rebuilds (once per tick) and fast local queries.
// Entities are inserted into exactly one cell based on their position.
//
// This is not a full physics broadphase: it is used for interest management
// (relevance) to avoid O(N^2) per-client scans.
type Grid struct {
	cellSize float32
	cells    map[int64][]uint32 // packed cell coords -> entity ids
}

func NewGrid(cellSize float32) *Grid {
	if cellSize <= 0 {
		cellSize = 8
	}
	return &Grid{cellSize: cellSize, cells: make(map[int64][]uint32, 256)}
}

func (g *Grid) SetCellSize(cellSize float32) {
	if cellSize <= 0 {
		return
	}
	g.cellSize = cellSize
	g.Clear()
}

func (g *Grid) Clear() {
	// Allocate a new map to drop old buckets quickly.
	g.cells = make(map[int64][]uint32, 256)
}

func (g *Grid) Insert(id uint32, x, y float32) {
	cx, cy := g.cellCoords(x, y)
	key := packCell(cx, cy)
	g.cells[key] = append(g.cells[key], id)
}

// QueryCircle returns candidate entity ids in cells overlapped by the circle.
// Callers should still do an exact distance check.
func (g *Grid) QueryCircle(cx, cy, r float32) []uint32 {
	if r <= 0 {
		return nil
	}
	minX := int32(math.Floor(float64((cx - r) / g.cellSize)))
	maxX := int32(math.Floor(float64((cx + r) / g.cellSize)))
	minY := int32(math.Floor(float64((cy - r) / g.cellSize)))
	maxY := int32(math.Floor(float64((cy + r) / g.cellSize)))

	// Rough capacity guess: (cells in bbox) * ~4
	capGuess := int((maxX-minX+1)*(maxY-minY+1)) * 4
	if capGuess < 0 {
		capGuess = 0
	}
	out := make([]uint32, 0, capGuess)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			key := packCell(x, y)
			if ids, ok := g.cells[key]; ok {
				out = append(out, ids...)
			}
		}
	}
	return out
}

func (g *Grid) cellCoords(x, y float32) (int32, int32) {
	return int32(math.Floor(float64(x / g.cellSize))), int32(math.Floor(float64(y / g.cellSize)))
}

func packCell(x, y int32) int64 {
	return (int64(x) << 32) | (int64(uint32(y)))
}
