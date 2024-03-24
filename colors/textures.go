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

type Checkerboard struct {
	Squares int
}

func (c Checkerboard) GetTextureColor(x, y float64) Color {
	r := 1 / float64(c.Squares)
	xV, yV := int(x/r), int(y/r)
	if (xV+yV)%2 == 0 {
		return Black
	}
	return White
}
