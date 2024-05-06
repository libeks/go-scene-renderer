package textures

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/sampler"
)

func NewPerlinNoiseTexture(gradient colors.Gradient) perlinNoiseTexture {
	p := sampler.NewPerlinNoise()
	return perlinNoiseTexture{
		noise:    p,
		gradient: gradient,
	}
}

type perlinNoiseTexture struct {
	noise    sampler.PerlinNoise
	gradient colors.Gradient
}

// x,y range from -1 to 1
func (p perlinNoiseTexture) GetTextureColor(x, y float64) colors.Color {
	valZeroOne := maths.Sigmoid(p.noise.GetFrameValue(x, y, 0) * 10)
	return p.gradient.Interpolate(valZeroOne)
}

func (p perlinNoiseTexture) GetFrameColor(x, y, t float64) colors.Color {
	valZeroOne := maths.Sigmoid(p.noise.GetFrameValue(x, y, t) * 10)
	return p.gradient.Interpolate(valZeroOne)
}
