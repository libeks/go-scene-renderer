package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/sampler"
	"github.com/libeks/go-scene-renderer/textures"
)

func PerlinColors() DynamicScene {
	perlin := sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 20}
	offset := 0.05
	redOffset := offset
	greenOffset := 2 * offset
	blueOffset := 0.0
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
	redTexture := sampler.TimeShiftedDynamic(texture, redOffset)
	greenTexture := sampler.TimeShiftedDynamic(texture, greenOffset)
	blueTexture := sampler.TimeShiftedDynamic(texture, blueOffset)
	background := BackgroundFromTexture(textures.RBGSamplerDynamicTexture(redTexture, greenTexture, blueTexture))
	return BackgroundScene(background)
}

func ShuffledColorRotation() DynamicScene {
	texture := sampler.RotatingCross(0.1)
	offset := 0.005
	redOffset := offset
	greenOffset := 2 * offset
	blueOffset := 0.0
	redTexture := sampler.TimeShiftedDynamic(texture, redOffset)
	greenTexture := sampler.TimeShiftedDynamic(texture, greenOffset)
	blueTexture := sampler.TimeShiftedDynamic(texture, blueOffset)
	background := BackgroundFromTexture(
		textures.GetRandomCellRemapper(
			textures.RBGSamplerDynamicTexture(redTexture, greenTexture, blueTexture),
			100,
			0.4, // number of cells being shuffled
		),
	)
	return BackgroundScene(background)
}

func FourColorSquares() DynamicScene {
	orange := colors.Hex("#FFA500")
	background := BackgroundFromTexture(
		textures.QuadriMapper(
			67,
			textures.StaticTexture(textures.SquareGradientTexture(colors.Red, orange, colors.Yellow, colors.Green)),
			textures.StaticTexture(textures.SquareGradientTexture(orange, colors.Yellow, colors.Green, colors.Red)),
			textures.StaticTexture(textures.SquareGradientTexture(colors.Yellow, colors.Green, colors.Red, orange)),
			textures.StaticTexture(textures.SquareGradientTexture(colors.Green, colors.Red, orange, colors.Blue)),
		),
	)
	return BackgroundScene(background)
}

func VerticalWiggler() DynamicScene {
	off := colors.Hex("#0D2B52")
	on := colors.Hex("#FFF6CC")
	nLines := 11
	angle := 0.015

	background := BackgroundFromTexture(
		textures.BinaryDynamicSamplerWithColors{
			DynamicSampler: sampler.Repeat(sampler.VerticalWiggler(nLines, angle), 5),
			On:             on,
			Off:            off,
		},
	)
	return BackgroundScene(background)
}

func VerticalLineConcentricCircles() DynamicScene {
	blue := colors.Hex("#75DDDE")
	red := colors.Hex("#E80E02")
	nLines := 13
	angle := 0.015
	lineWidth := 1 / float64(nLines)
	s := sampler.VerticalLines(nLines)

	for i := range 8 {
		s = sampler.InvertAndRotateCircle(
			s,
			float64(i+3)*lineWidth,
			angle,
		)
	}

	background := BackgroundFromTexture(
		textures.StaticTexture(
			textures.BinarySamplerWithColors{
				StaticSampler: s,
				On:            blue,
				Off:           red,
			},
		),
	)
	return BackgroundScene(background)
}

func ShuffledConcentricCircles() DynamicScene {
	on := colors.Red
	off := colors.White
	nLines := 20
	angle := 1.00
	background := BackgroundFromTexture(
		textures.GetRandomCellRemapper(
			textures.BinaryDynamicSamplerWithColors{
				DynamicSampler: sampler.Repeat(sampler.VerticalWiggler(nLines, angle), 5),
				On:             on,
				Off:            off,
			},
			20,
			0.8, // number of cells being shuffled
		),
	)
	return BackgroundScene(background)
}

func ConcentricCircles() DynamicScene {
	on := colors.Red
	off := colors.White
	background := BackgroundFromTexture(
		textures.StaticTexture(
			textures.BinarySamplerWithColors{
				StaticSampler: sampler.ConcentricCircles(0.1),
				On:            on,
				Off:           off,
			},
		),
	)
	return BackgroundScene(background)

}
