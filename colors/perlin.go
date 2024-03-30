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
	randConst   = float64(10000)
)

func NewPerlinNoiseTexture(gradient Gradient) perlinNoiseTexture {
	p := NewPerlinNoise()
	return perlinNoiseTexture{
		noise:    p,
		gradient: gradient,
	}
}

func NewRandomPerlinNoiseTexture(gradient Gradient) perlinNoiseTexture {
	p := NewRandomPerlinNoise()
	return perlinNoiseTexture{
		noise:    p,
		gradient: gradient,
	}
}

type perlinNoiseTexture struct {
	noise    PerlinNoise
	gradient Gradient
}

// x,y range from -1 to 1
func (p perlinNoiseTexture) GetTextureColor(x, y float64) Color {
	valZeroOne := maths.Sigmoid(p.noise.GetFrameValue(x, y, 0) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p perlinNoiseTexture) GetFrameColor(x, y, t float64) Color {
	valZeroOne := maths.Sigmoid(p.noise.GetFrameValue(x, y, t) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

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
