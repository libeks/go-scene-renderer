package textures

import "fmt"

func GetCoord(x, y, n int) int {
	return x*n + y
}

func IndexToCoord(i, n int) (int, int) {
	y := i % n
	x := i / n
	return x, y
}

func NewGrid(n int) grid {
	return grid{
		vals: make([]float64, n*n),
		N:    n,
	}
}

type grid struct {
	vals []float64
	N    int
}

func (g grid) Get(x, y int) float64 {
	return g.vals[g.getCoord(x, y)]
}

func (g grid) getCoord(x, y int) int {
	return GetCoord(x, y, g.N)
}

func (g grid) Set(x, y int, val float64) {
	g.vals[g.getCoord(x, y)] = val
}

func (g grid) String() string {
	return fmt.Sprintf("Grid: %v", g.vals)
}
