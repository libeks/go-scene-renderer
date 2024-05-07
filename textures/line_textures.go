package textures

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/sampler"
)

func RotatingLine(on, off colors.Color, thickness float64) BinaryDynamicSamplerWithColors {
	return BinaryDynamicSamplerWithColors{
		DynamicSampler: sampler.RotatingLine(thickness),
		On:             on,
		Off:            off,
	}
}

func RotatingCross(on, off colors.Color, thickness float64) BinaryDynamicSamplerWithColors {
	return BinaryDynamicSamplerWithColors{
		DynamicSampler: sampler.RotatingCross(thickness),
		On:             on,
		Off:            off,
	}
}

func PulsingSquare(on, off colors.Color) AnimatedTexture {
	return BinaryAnimatedSamplerWithColors{
		Sampler: sampler.PulsingSquare{},
		On:      on,
		Off:     off,
	}
}

func VerticalLine(on, off colors.Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.VerticalLine(thickness),
		On:            on,
		Off:           off,
	}
}

func HorizontalLine(on, off colors.Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.HorizontalLine(thickness),
		On:            on,
		Off:           off,
	}
}

func Cross(on, off colors.Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Cross(thickness),
		On:            on,
		Off:           off,
	}
}

func Square(on, off colors.Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Square{
			Thickness: thickness,
		},
		On:  on,
		Off: off,
	}
}

func RoundedSquare(on, off colors.Color, halfwidth, radius float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.RoundedSquare{
			HalfWidth: halfwidth,
			Radius:    radius,
		},
		On:  on,
		Off: off,
	}
}

func Circle(on, off colors.Color, radius float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Circle{
			Thickness: radius,
		},
		On:  on,
		Off: off,
	}
}

func GetSpecialMapper(on, off colors.Color, thickness float64) StaticMapper {
	return StaticMapper{
		Mapping: []TextureValueMapping{
			{0.9, Square(on, off, 1.0)},
			{0.8, Square(on, off, max(0.7, 2*thickness))},
			{0.7, Cross(on, off, thickness)},
			{0.5, HorizontalLine(on, off, thickness)},
			{0.4, VerticalLine(on, off, thickness)},
			{0.1, Circle(on, off, thickness)},
			{0.0, Uniform(off)},
		},
	}
}

func Rainbow() AnimatedTexture {
	return rainbow{}
}

type rainbow struct{}

func (r rainbow) GetFrameColor(x, y, t float64) colors.Color {
	return colors.HSL(math.Mod(t, 1), 0.75, 0.5)
}
