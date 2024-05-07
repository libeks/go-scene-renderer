package colors

import (
	"fmt"
	"image/color"
	"math"

	"github.com/crazy3lf/colorconv"
)

const (
	maxUInt32   = 0xffff
	gammaFactor = 1.8
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
	return uint32(maxUInt32 * gamma(c.R)),
		uint32(maxUInt32 * gamma(c.G)),
		uint32(maxUInt32 * gamma(c.B)),
		maxUInt32
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", floatToUInt32(c.R), floatToUInt32(c.G), floatToUInt32(c.B))
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

// Inverse of gamma fuction, so that gamma(inverseGamma) = input
func inverseGamma(v float64) float64 {
	return math.Pow(v, float64(1)/1.8)
}

func gamma(v float64) float64 {
	return math.Pow(v, 1.8)
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
		R: inverseGamma(uInt32ToFloat(r)),
		G: inverseGamma(uInt32ToFloat(g)),
		B: inverseGamma(uInt32ToFloat(b)),
	}
	return c
}

// h,s,v each range from 0 to 1
func HSL(h, s, l float64) Color {
	r, g, b, err := colorconv.HSLToRGB(h*360, s, l)
	if err != nil {
		fmt.Printf("hsl %.3f, %.3f, %.3f -> rgb: %d %d %d\n", h, s, l, r, g, b)
		panic(err)
	}
	c := Color{
		R: inverseGamma(uInt8ToFloat(r)),
		G: inverseGamma(uInt8ToFloat(g)),
		B: inverseGamma(uInt8ToFloat(b)),
	}
	return c
}

// h,s,v each range from 0 to 1
func HSV(h, s, v float64) Color {
	r, g, b, err := colorconv.HSVToRGB(h*360, s, v)
	if err != nil {
		fmt.Printf("hsv %.3f, %.3f, %.3f -> rgb: %d %d %d\n", h, s, v, r, g, b)
		panic(err)
	}
	c := Color{
		R: inverseGamma(uInt8ToFloat(r)),
		G: inverseGamma(uInt8ToFloat(g)),
		B: inverseGamma(uInt8ToFloat(b)),
	}
	return c
}

func uInt8ToFloat(r uint8) float64 {
	return float64(r) / float64(0xff)
}

func uInt32ToFloat(r uint32) float64 {
	return float64(r) / float64(0xff)
}

func floatToUInt32(r float64) uint32 {
	return uint32(r * float64(0xff))
}

func ToInterfaceSlice(colors []Color) []color.Color {
	out := make([]color.Color, len(colors))
	for i := range colors {
		out[i] = colors[i]
	}
	return out
}
