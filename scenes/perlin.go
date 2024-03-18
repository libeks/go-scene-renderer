package scenes

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/objects"
)

var (
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = int32(5)
	perlinSeed  = int64(103)
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

func (p PerlinNoise) GetFrameColor(x, y float64, t float64) color.Color {
	valZeroOne := sigmoid(p.noise.Noise2D(x, y) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p PerlinNoise) GetFrame(t float64) Frame {
	return PerlinNoise{
		noise:    perlin.NewPerlinRandSource(1.0+t*2, perlinBeta+t*2, perlinN, rand.NewSource(perlinSeed)),
		gradient: p.gradient,
	}
}

func (p PerlinNoise) GetColor(x, y float64) color.Color {
	valZeroOne := sigmoid(p.noise.Noise2D(x, y) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p PerlinNoise) GetObjects() []objects.Object {
	return nil
}
