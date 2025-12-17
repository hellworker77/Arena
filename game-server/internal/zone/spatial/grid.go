package spatial

import "math"

type CellKey struct{ X, Y int32 }

type Grid struct {
	CellSize int16
	cells    map[CellKey][]uint32
	xs       map[uint32]int16
	ys       map[uint32]int16
}

func New(cellSize int16) *Grid {
	if cellSize <= 0 {
		cellSize = 8
	}
	return &Grid{
		CellSize: cellSize,
		cells:    make(map[CellKey][]uint32),
		xs:       make(map[uint32]int16),
		ys:       make(map[uint32]int16),
	}
}

func (g *Grid) Clear() {
	for k := range g.cells {
		delete(g.cells, k)
	}
	for k := range g.xs {
		delete(g.xs, k)
	}
	for k := range g.ys {
		delete(g.ys, k)
	}
}

func (g *Grid) Insert(eid uint32, x, y int16) {
	g.xs[eid] = x
	g.ys[eid] = y
	ck := g.cellOf(x, y)
	g.cells[ck] = append(g.cells[ck], eid)
}

func (g *Grid) cellOf(x, y int16) CellKey {
	cs := float64(g.CellSize)
	return CellKey{
		X: int32(math.Floor(float64(x) / cs)),
		Y: int32(math.Floor(float64(y) / cs)),
	}
}

func (g *Grid) QueryCircle(cx, cy, r int16, out []uint32) []uint32 {
	if r <= 0 {
		return out
	}
	minX := int32(math.Floor(float64(cx-r) / float64(g.CellSize)))
	maxX := int32(math.Floor(float64(cx+r) / float64(g.CellSize)))
	minY := int32(math.Floor(float64(cy-r) / float64(g.CellSize)))
	maxY := int32(math.Floor(float64(cy+r) / float64(g.CellSize)))

	rr := int32(r) * int32(r)
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			ids := g.cells[CellKey{X: x, Y: y}]
			for _, eid := range ids {
				ex := int32(g.xs[eid])
				ey := int32(g.ys[eid])
				dx := ex - int32(cx)
				dy := ey - int32(cy)
				if dx*dx+dy*dy <= rr {
					out = append(out, eid)
				}
			}
		}
	}
	return out
}

func (g *Grid) GetPos(eid uint32) (x, y int16, ok bool) {
	x, ok = g.xs[eid]
	if !ok {
		return 0, 0, false
	}
	y = g.ys[eid]
	return x, y, true
}
