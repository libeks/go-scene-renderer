package colors

import (
	"github.com/libeks/go-scene-renderer/sampler"
)

func RotatingLine(on, off Color, thickness float64) BinaryDynamicSamplerWithColors {
	return BinaryDynamicSamplerWithColors{
		DynamicSampler: sampler.RotatingLine(thickness),
		On:             on,
		Off:            off,
	}
}

func RotatingCross(on, off Color, thickness float64) AnimatedTexture {
	return BinaryDynamicSamplerWithColors{
		DynamicSampler: sampler.RotatingCross(thickness),
		On:             on,
		Off:            off,
	}
}

func PulsingSquare(on, off Color) AnimatedTexture {
	return BinaryAnimatedSamplerWithColors{
		Sampler: sampler.PulsingSquare{},
		On:      on,
		Off:     off,
	}
}

func VerticalLine(on, off Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.VerticalLine(thickness),
		On:            on,
		Off:           off,
	}
}

func HorizontalLine(on, off Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.HorizontalLine(thickness),
		On:            on,
		Off:           off,
	}
}

func Cross(on, off Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Cross(thickness),
		On:            on,
		Off:           off,
	}
}

func Square(on, off Color, thickness float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Square{
			Thickness: thickness,
		},
		On:  on,
		Off: off,
	}
}

func RoundedSquare(on, off Color, halfwidth, radius float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.RoundedSquare{
			HalfWidth: halfwidth,
			Radius:    radius,
		},
		On:  on,
		Off: off,
	}
}

func Circle(on, off Color, radius float64) Texture {
	return BinarySamplerWithColors{
		StaticSampler: sampler.Circle{
			Thickness: radius,
		},
		On:  on,
		Off: off,
	}
}
