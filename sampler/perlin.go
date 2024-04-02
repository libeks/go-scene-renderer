package sampler

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
	randConst   = float64(10000)
)

type PerlinNoise struct {
	noise   *perlin.Perlin
	offsetX float64
	offsetY float64
}

func NewPerlinNoise() PerlinNoise {
	return PerlinNoise{noise: perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed))}
}

func NewRandomPerlinNoise() PerlinNoise {
	return PerlinNoise{
		noise:   perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed)),
		offsetX: rand.Float64() * randConst,
		offsetY: rand.Float64() * randConst,
	}
}

func (p PerlinNoise) GetFrameValue(x, y, t float64) float64 {
	valZeroOne := maths.Sigmoid(p.noise.Noise3D(x+p.offsetX, y+p.offsetY, t) * 5)
	return valZeroOne
}
