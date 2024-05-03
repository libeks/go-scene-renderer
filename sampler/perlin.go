package sampler

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
)

var (
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = int32(10)
	perlinSeed  = int64(109)
)

type PerlinNoise struct {
	noise   *perlin.Perlin
	offsetX float64
	offsetY float64
}

func NewPerlinNoise() PerlinNoise {
	return PerlinNoise{noise: perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed))}
}

// returns a value from -1 to 1, based on Perlin Noise
func (p PerlinNoise) GetFrameValue(x, y, t float64) float64 {
	val := p.noise.Noise3D(x+p.offsetX, y+p.offsetY, t)
	return val
}
