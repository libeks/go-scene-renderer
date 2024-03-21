package scenes

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
)

var (
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = int32(5)
	perlinSeed  = int64(109)
)

func NewPerlinNoise(gradient color.Gradient) PerlinNoise {
	p := perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed))
	return PerlinNoise{
		noise:    p,
		gradient: gradient,
	}
}

type PerlinNoise struct {
	noise    *perlin.Perlin
	gradient color.Gradient
}

// x,y ranges from -1 to 1
func (p PerlinNoise) GetFrameColor(x, y float64, t float64) color.Color {
	valZeroOne := maths.Sigmoid(p.noise.Noise2D(x, y) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p PerlinNoise) GetFrame(t float64) Frame {
	return p
	// return PerlinNoise{
	// 	noise:    perlin.NewPerlinRandSource(1.0+t*2, perlinBeta+t*2, perlinN, rand.NewSource(perlinSeed)),
	// 	gradient: p.gradient,
	// }
}

func (p PerlinNoise) GetColor(x, y float64) color.Color {
	valZeroOne := maths.Sigmoid(p.noise.Noise2D(x, y) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

// x,y range from -1 to 1
func (p PerlinNoise) GetTextureColor(x, y float64) color.Color {
	// fmt.Printf("%.3f, %.3f\n", x, y)
	return p.GetColor(x*2-1, y*2-1)
}

func (d PerlinNoise) Flatten() []*objects.Triangle {
	return nil
}
