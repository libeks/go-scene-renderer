package colors

import (
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/sampler"
)

func NewPerlinNoiseTexture(gradient Gradient) perlinNoiseTexture {
	p := sampler.NewPerlinNoise()
	return perlinNoiseTexture{
		noise:    p,
		gradient: gradient,
	}
}

func NewRandomPerlinNoiseTexture(gradient Gradient) perlinNoiseTexture {
	p := sampler.NewRandomPerlinNoise()
	return perlinNoiseTexture{
		noise:    p,
		gradient: gradient,
	}
}

type perlinNoiseTexture struct {
	noise    sampler.PerlinNoise
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
