package colors

import (
	"fmt"
	"image/color"
	"math"
)

const (
	maxUInt32 = 0xffff
)

var (
	Black   = GrayscaleColor(0)
	White   = GrayscaleColor(1)
	Gray    = GrayscaleColor(0.5)
	Red     = Color{R: 1, G: 0, B: 0}
	Blue    = Color{R: 0, G: 0, B: 1}
	Green   = Color{R: 0, G: 1, B: 0}
	Cyan    = Color{R: 0, G: 1, B: 1}
	Magenta = Color{R: 1, G: 0, B: 1}
	Yellow  = Color{R: 1, G: 1, B: 0}

	Grayscale = SimpleGradient{
		Start: Black,
		End:   White,
	}
)

type Color struct {
	// represented in the range 0-1
	R float64
	G float64
	B float64
}

func (c Color) RGBA() (r, g, b, a uint32) {
	// only apply gamma correction when rendered values are requested,
	// keep raw values otherwise
	return uint32(maxUInt32 * math.Pow(c.R, 1.8)),
		uint32(maxUInt32 * math.Pow(c.G, 1.8)),
		uint32(maxUInt32 * math.Pow(c.B, 1.8)),
		maxUInt32
}

// Average the colors in the slice
func Average(colors []Color) Color {
	if len(colors) == 1 {
		return colors[0]
	}
	retCol := Color{}
	n := float64(len(colors))
	for _, c := range colors {
		retCol.R += c.R
		retCol.G += c.G
		retCol.B += c.B
	}
	retCol.R = retCol.R / n
	retCol.G = retCol.G / n
	retCol.B = retCol.B / n
	return retCol
}

func (c Color) Add(d Color) Color {
	// add the color components of the two colors, maxing out at 255
	return Color{
		min(1, c.R+d.R),
		min(1, c.G+d.G),
		min(1, c.B+d.B),
	}
}

func GrayscaleColor(v float64) Color {
	return Color{
		R: v,
		G: v,
		B: v,
	}
}

// Parses Hex color value into Color
// adapted from https://stackoverflow.com/a/54200713
func Hex(s string) Color {
	var r, g, b uint32
	switch len(s) {
	case 7:
		_, _ = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	case 4:
		_, _ = fmt.Sscanf(s, "#%1x%1x%1x", &r, &g, &b)
		// Double the hex digits:
		r *= 17
		g *= 17
		b *= 17
	}
	c := Color{
		R: uInt32ToFloat(r),
		G: uInt32ToFloat(g),
		B: uInt32ToFloat(b),
	}
	return c
}

func uInt32ToFloat(r uint32) float64 {
	return float64(r) / float64(0xff)
}

func ToInterfaceSlice(colors []Color) []color.Color {
	out := make([]color.Color, len(colors))
	for i := range colors {
		out[i] = colors[i]
	}
	return out
}
