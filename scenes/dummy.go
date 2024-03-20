package scenes

import (
	"math/rand"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/objects"
)

type Dummy struct{}

func (d Dummy) GetFrameColor(x, y int, t float64) color.Color {
	return color.GrayscaleColor(t)
}

type HorizGradient struct {
	Gradient color.Gradient
}

func (d HorizGradient) GetFrameColor(x, y float64, t float64) color.Color {
	valZeroOne := x/2 + 0.5
	return d.Gradient.Interpolate(valZeroOne)
}

func (d HorizGradient) GetColorPalette(t float64) []color.Color {
	return color.GetGradientColorPalette(d.Gradient)
}

type Uniform struct {
	Color color.Color
}

func (d Uniform) GetColor(x, y float64) color.Color {
	return d.Color
}

func (d Uniform) GetFrame(t float64) Frame {
	return d
}

func (d Uniform) Flatten() []*objects.Triangle {
	return nil
}

type Random struct{}

func (d Random) GetColor(x, y float64) color.Color {
	if rand.Float32() > 0.5 {
		return color.Black
	}
	return color.White
}

func (d Random) GetFrame(t float64) Frame {
	return d
}

func (d Random) Flatten() []*objects.Triangle {
	return nil
}
