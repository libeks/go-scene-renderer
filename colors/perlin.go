package colors

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/libeks/go-scene-renderer/maths"
)

var (
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = int32(10)
	perlinSeed  = int64(109)
)

func NewPerlinNoise(gradient Gradient) perlinNoise {
	p := perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed))
	return perlinNoise{
		noise:    p,
		gradient: gradient,
	}
}

type perlinNoise struct {
	noise    *perlin.Perlin
	gradient Gradient
}

// x,y range from -1 to 1
func (p perlinNoise) GetTextureColor(x, y float64) Color {
	valZeroOne := maths.Sigmoid(p.noise.Noise2D(x, y) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p perlinNoise) GetFrameColor(x, y, t float64) Color {
	valZeroOne := maths.Sigmoid(p.noise.Noise3D(x, y, t) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p perlinNoise) GetFrameValue(x, y, t float64) float64 {
	valZeroOne := maths.Sigmoid(p.noise.Noise3D(x, y, t) * 10)
	return valZeroOne
}
