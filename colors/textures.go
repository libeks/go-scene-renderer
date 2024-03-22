package colors

import (
	"math/rand"
)

type HorizGradient struct {
	Gradient
}

func (d HorizGradient) GetTextureColor(x, y float64) Color {
	valZeroOne := x/2 + 0.5
	return d.Gradient.Interpolate(valZeroOne)
}

type Uniform struct {
	Color
}

func (d Uniform) GetTextureColor(x, y float64) Color {
	return d.Color
}

type Random struct{}

func (d Random) GetTextureColor(x, y float64) Color {
	if rand.Float32() > 0.5 {
		return Black
	}
	return White
}
