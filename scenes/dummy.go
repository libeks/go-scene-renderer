package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
)

type Dummy struct{}

func (d Dummy) GetPixel(x, y int, t float64) color.Color {
	return color.GrayscaleColor(t)
}

func (d Dummy) GetColor(x, y float64, t float64) color.Color {
	return color.GrayscaleColor(t)
}

type HorizGradient struct {
	Gradient color.Gradient
}

func (d HorizGradient) GetColor(x, y float64, t float64) color.Color {
	valZeroOne := x/2 + 0.5
	return d.Gradient.Interpolate(valZeroOne)
}

func (d HorizGradient) GetColorPalette(t float64) []color.Color {
	return color.GetGradientColorPalette(d.Gradient)
}

type Uniform struct {
	Color color.Color
}

func (d Uniform) GetColor(x, y float64, t float64) color.Color {
	return d.Color
}

func (d Uniform) GetColorPalette(t float64) []color.Color {
	return []color.Color{d.Color}
}
