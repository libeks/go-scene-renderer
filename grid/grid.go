package grid

import (
	"fmt"
	"slices"
)

func GetCoord(x, y, n int) int {
	return x*n + y
}

func IndexToCoord(i, n int) (int, int) {
	y := i % n
	x := i / n
	return x, y
}

func NewGrid(n int) Grid {
	return Grid{
		vals: make([]float64, n*n),
		N:    n,
	}
}

type Grid struct {
	vals []float64
	N    int
}

func (g Grid) Get(x, y int) float64 {
	return g.vals[g.getCoord(x, y)]
}

func (g Grid) getCoord(x, y int) int {
	return GetCoord(x, y, g.N)
}

func (g Grid) Set(x, y int, val float64) {
	g.vals[g.getCoord(x, y)] = val
}

func (g Grid) String() string {
	return fmt.Sprintf("Grid: %v", g.vals)
}

func NewDynamicGrid() DynamicGrid {
	return DynamicGrid{
		index: []float64{},
		grids: []Grid{},
	}
}

type DynamicGrid struct {
	index []float64 // maps from index in grids to frame float. First should be 0.0, last should be 1.0. This must be sorted.
	grids []Grid    // same cardinality as index, which indexes the frame t to this grid
}

func (g *DynamicGrid) getIndex(t float64) int {
	idx, _ := slices.BinarySearch(g.index, t) // most of the time there won't be an exact match, ignore second parameter
	return idx
}

func (g *DynamicGrid) AddFrame(t float64, newGrid Grid) {
	idx := g.getIndex(t)
	g.index = slices.Insert(g.index, idx, t)
	g.grids = slices.Insert(g.grids, idx, newGrid)
}

func (g *DynamicGrid) GetFrame(t float64) Grid {
	idx := g.getIndex(t)
	return g.grids[idx]
}
