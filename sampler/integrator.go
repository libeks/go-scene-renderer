package sampler

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/grid"
)

func Integrate(s DynamicSampler, steps int, nBlocks int, intConstant float64) grid.DynamicGrid {
	fmt.Printf("Generating scene integral...")
	invStep := 1 / float64(steps)
	d := 1 / float64(nBlocks)
	g := grid.NewGrid(nBlocks)
	grids := grid.NewDynamicGrid()
	for i := range steps {
		t := float64(i) / float64(steps-1)
		sampler := s.GetFrame(t)
		newGrid := grid.NewGrid(nBlocks)
		for xIdx := range nBlocks {
			x := (float64(xIdx) + 0.5) * d
			for yIdx := range nBlocks {
				y := (float64(yIdx) + 0.5) * d
				val := sampler.GetValue(x, y) * invStep * intConstant
				newGrid.Set(xIdx, yIdx, val+g.Get(xIdx, yIdx))
			}
		}
		grids.AddFrame(t, newGrid)
		g = newGrid
	}
	fmt.Printf(" Done!\n")
	return grids
}
