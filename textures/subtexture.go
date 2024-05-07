package textures

import (
	"fmt"
	"math/rand"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/grid"
	"github.com/libeks/go-scene-renderer/sampler"
)

func DynamicSubtexturer(s AnimatedTexture, n int, sampler sampler.Sampler) *dynamicSubtexturer {
	return &dynamicSubtexturer{
		subtexture:   s,
		N:            n,
		PointSampler: sampler,
	}
}

type dynamicSubtexturer struct {
	subtexture   AnimatedTexture
	N            int // number of squares to tile
	PointSampler sampler.Sampler
}

func (s *dynamicSubtexturer) getCellValue(xMeta, yMeta, t float64) float64 {
	val := s.PointSampler.GetFrameValue(xMeta, yMeta, t)

	return val
}

func (s dynamicSubtexturer) GetFrame(t float64) Texture {
	d := 1 / float64(s.N)
	// calculate the grid sampler values to use in this frame
	grid := grid.NewGrid(s.N)
	for x := range s.N {
		for y := range s.N {
			grid.Set(x, y, s.getCellValue(float64(x)*d, float64(y)*d, t))
		}
	}
	fmt.Printf("grid %s\n", grid)
	return staticSubtexture{
		AnimatedTexture: s.subtexture,
		N:               s.N,
		Grid:            grid,
	}
}

func DynamicGridSubtexturer(s AnimatedTexture, N int, g grid.DynamicGrid) *dynamicGridSubtexturer {
	return &dynamicGridSubtexturer{
		DynamicGrid: g,
		subtexture:  s,
		N:           N,
	}
}

type dynamicGridSubtexturer struct {
	grid.DynamicGrid
	subtexture AnimatedTexture
	N          int
}

func (s dynamicGridSubtexturer) GetFrame(t float64) Texture {
	return staticSubtexture{
		AnimatedTexture: s.subtexture,
		N:               s.N,
		Grid:            s.DynamicGrid.GetFrame(t),
	}
}

type staticSubtexture struct {
	AnimatedTexture // renders the visual of each subcell, indexed by t, which doesn't have to be by time
	N               int
	grid.Grid       // grid contains the sampler values for this frame
}

func (s staticSubtexture) GetTextureColor(b, c float64) colors.Color {
	d := 1 / float64(s.N)
	xMeta, xValue := bucketRemainder(b, d)
	yMeta, yValue := bucketRemainder(c, d)

	tHere := s.Grid.Get(int(xMeta*float64(s.N)), int(yMeta*float64(s.N)))
	return s.AnimatedTexture.GetFrameColor(xValue, yValue, tHere)
}

func GetRandomCellRemapper(d DynamicTexture, n int, threshold float64) DynamicTexture {
	cellMapping := make([]int, n*n)
	for i := range n * n {
		cellMapping[i] = i
	}
	rand.Shuffle(n*n, func(i, j int) {
		if rand.Float64() > threshold {
			cellMapping[i], cellMapping[j] = cellMapping[j], cellMapping[i]
		}
	})
	return cellRemapperDynamic{
		DynamicTexture: d,
		cellMapping:    cellMapping,
		N:              n,
	}
}

// implements DynamicTexture
type cellRemapperDynamic struct {
	DynamicTexture
	cellMapping []int // maps from 2d int,int to 2d int,int; one-to-one and onto
	N           int
}

func (s cellRemapperDynamic) GetFrame(t float64) Texture {
	return cellRemapper{
		Texture:     s.DynamicTexture.GetFrame(t),
		cellMapping: s.cellMapping,
		N:           s.N,
	}
}

type cellRemapper struct {
	Texture
	cellMapping []int
	N           int
}

func (s cellRemapper) GetTextureColor(x, y float64) colors.Color {
	d := 1 / float64(s.N)
	xMeta, xValue := bucketRemainder(x, d)
	yMeta, yValue := bucketRemainder(y, d)
	coord := grid.GetCoord(int(xMeta*float64(s.N)), int(yMeta*float64(s.N)), s.N)
	newCoord := s.cellMapping[coord]
	xInt, yInt := grid.IndexToCoord(newCoord, s.N)
	newX, newY := float64(xInt)*d+xValue*d, float64(yInt)*d+yValue*d
	return s.Texture.GetTextureColor(newX, newY)
}

func QuadriMapper(n int, A, B, C, D DynamicTexture) DynamicTexture {
	return quadriMapperDynamic{
		N: n,
		A: A,
		B: B,
		C: C,
		D: D,
	}
}

type quadriMapperDynamic struct {
	N int
	A DynamicTexture
	B DynamicTexture
	C DynamicTexture
	D DynamicTexture
}

func (s quadriMapperDynamic) GetFrame(t float64) Texture {
	return quadriMapper{
		N: s.N,
		A: s.A.GetFrame(t),
		B: s.B.GetFrame(t),
		C: s.C.GetFrame(t),
		D: s.D.GetFrame(t),
	}
}

type quadriMapper struct {
	N int
	A Texture
	B Texture
	C Texture
	D Texture
}

func (s quadriMapper) GetTextureColor(x, y float64) colors.Color {
	d := 1 / float64(s.N)
	xMeta, _ := bucketRemainder(x, d)
	yMeta, _ := bucketRemainder(y, d)
	var text Texture
	quadrant := int(xMeta/d)%2 + 2*(int(yMeta/d)%2)
	// fmt.Printf("x, y: %.3f, %.3f; %d %d Quadrant: %d\n", xMeta, yMeta, int(xMeta*float64(s.N)), int(yMeta*float64(s.N)), quadrant)
	switch quadrant {
	case 0:
		text = s.A
	case 1:
		text = s.B
	case 2:
		text = s.C
	case 3:
		text = s.D
	}
	return text.GetTextureColor(x, y)
}
