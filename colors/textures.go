package colors

import (
	"math/rand"
)

type HorizontalGradient struct {
	Gradient
}

func (d HorizontalGradient) GetTextureColor(x, y float64) Color {
	return d.Gradient.Interpolate(x)
}

type VerticalGradient struct {
	Gradient
}

func (d VerticalGradient) GetTextureColor(x, y float64) Color {
	return d.Gradient.Interpolate(y)
}

type Fuzzy struct {
	Texture Texture
	StdDev  float64
}

func (g Fuzzy) GetTextureColor(x, y float64) Color {
	dx := rand.NormFloat64() * g.StdDev
	dy := rand.NormFloat64() * g.StdDev
	x = x + dx
	if x < 0 {
		x = 0
	}
	if x > 1 {
		x = 1
	}
	y = y + dy
	if y < 0 {
		y = 0
	}
	if y > 1 {
		y = 1
	}
	return g.Texture.GetTextureColor(x, y)
}

type FuzzyDynamic struct {
	Texture DynamicTexture
	StdDev  float64
}

func (g FuzzyDynamic) GetFrame(t float64) Texture {
	return Fuzzy{
		Texture: g.Texture.GetFrame(t),
		StdDev:  g.StdDev,
	}
}

func Uniform(c Color) Texture {
	return uniform{c}
}

type uniform struct {
	Color
}

func (d uniform) GetTextureColor(x, y float64) Color {
	return d.Color
}

func Random() Texture {
	return random{}
}

type random struct{}

func (d random) GetTextureColor(x, y float64) Color {
	if rand.Float32() > 0.5 {
		return Black
	}
	return White
}

type MiddleBand struct {
	Min float64
	Max float64
}

func (c MiddleBand) GetAlpha(x, y float64) bool {
	return y > c.Min && y < c.Max
}

type Checkerboard struct {
	Squares int
}

func (c Checkerboard) GetTextureColor(x, y float64) Color {
	// don't render outside of texture boundaries
	if x < 0 || x > 1 || y < 0 || y > 1 {
		return Red
	}
	r := 1 / float64(c.Squares)
	xV, yV := int(x/r), int(y/r)
	if (xV+yV)%2 == 0 {
		return Black
	}
	return White
}

func (c Checkerboard) GetAlpha(x, y float64) bool {
	r := 1 / float64(c.Squares)
	xV, yV := int(x/r), int(y/r)
	return (xV+yV)%2 == 0
}
