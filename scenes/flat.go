package scenes

import (
	"github.com/libeks/go-scene-renderer/sampler"
	"github.com/libeks/go-scene-renderer/textures"
)

func PerlinColors() DynamicScene {
	perlin := sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 20}
	offset := 0.05
	redOffset := offset
	greenOffset := 2 * offset
	blueOffset := 0.0
	// redPerlin := sampler.Shifted(perlin, redOffset, redOffset)
	// greenPerlin := sampler.Shifted(perlin, greenOffset, greenOffset)
	// bluePerlin := sampler.Shifted(perlin, blueOffset, blueOffset)
	redPerlin := sampler.TimeShifted(perlin, redOffset)
	greenPerlin := sampler.TimeShifted(perlin, greenOffset)
	bluePerlin := sampler.TimeShifted(perlin, blueOffset)
	background := BackgroundFromTexture(textures.RBGSamplerTexture(redPerlin, greenPerlin, bluePerlin))
	return BackgroundScene(background)
}

func ColorRotation() DynamicScene {
	texture := sampler.RotatingCross(0.1)
	offset := 0.005
	redOffset := offset
	greenOffset := 2 * offset
	blueOffset := 0.0
	// redPerlin := sampler.Shifted(perlin, redOffset, redOffset)
	// greenPerlin := sampler.Shifted(perlin, greenOffset, greenOffset)
	// bluePerlin := sampler.Shifted(perlin, blueOffset, blueOffset)
	redTexture := sampler.TimeShiftedDynamic(texture, redOffset)
	greenTexture := sampler.TimeShiftedDynamic(texture, greenOffset)
	blueTexture := sampler.TimeShiftedDynamic(texture, blueOffset)
	background := BackgroundFromTexture(textures.RBGSamplerDynamicTexture(redTexture, greenTexture, blueTexture))
	return BackgroundScene(background)
}
