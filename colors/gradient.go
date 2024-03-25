package colors

import (
	"math"
)

type Gradient interface {
	// v=0.0 returns start color
	// v=1.0 return end color
	Interpolate(v float64) Color
}

func GetGradientColorPalette(g Gradient) []Color {
	out := make([]Color, 256)
	for i := range 256 {
		out[i] = g.Interpolate(float64(i) / 256.0)
	}
	return out
}

type SimpleGradient struct {
	Start Color
	End   Color
}

func (g SimpleGradient) Interpolate(v float64) Color {
	return Color{
		R: interpolate(v, g.Start.R, g.End.R),
		G: interpolate(v, g.Start.G, g.End.G),
		B: interpolate(v, g.Start.B, g.End.B),
	}
}

// Subsample the original gradient to produce a new one
func (g SimpleGradient) Subsample(a, b float64) Gradient {
	return SimpleGradient{
		g.Interpolate(a),
		g.Interpolate(b),
	}
}

type LinearGradient struct {
	Points []Color
}

func (g LinearGradient) Interpolate(v float64) Color {
	nPoints := len(g.Points)
	segmentLength := 1.0 / float64(nPoints-1)
	segment := int(math.Floor(v / segmentLength))
	remainder := math.Mod(v, segmentLength) / segmentLength
	if segment == nPoints-1 {
		return g.Points[nPoints-1]
	}
	return SimpleGradient{g.Points[segment], g.Points[segment+1]}.Interpolate(remainder)
}

func interpolate(t, a, b float64) float64 {
	return (b*t + a*(1-t))
}
