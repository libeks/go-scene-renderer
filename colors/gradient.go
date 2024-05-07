package colors

import (
	"fmt"
	"math"
)

type Gradient interface {
	// v=0.0 returns start color
	// v=1.0 return end color
	Interpolate(v float64) Color
}

type SimpleGradient struct {
	Start Color
	End   Color
}

func (g SimpleGradient) String() string {
	return fmt.Sprintf("SimpleGradient [%s %s]", g.Start, g.End)
}

func (g SimpleGradient) Interpolate(v float64) Color {
	if v < 0 {
		return Red
	}
	if v > 1 {
		return Blue
	}
	return Color{
		R: interpolate(v, g.Start.R, g.End.R),
		G: interpolate(v, g.Start.G, g.End.G),
		B: interpolate(v, g.Start.B, g.End.B),
	}
}

type LinearGradient struct {
	Points []Color
}

func (g LinearGradient) Interpolate(v float64) Color {
	if v <= 0 {
		return g.Points[0]
	}
	if v >= 1 {
		return g.Points[len(g.Points)-1]
	}
	nPoints := len(g.Points)
	segmentLength := 1.0 / float64(nPoints-1)
	segment := int(math.Floor(v / segmentLength))
	remainder := math.Mod(v, segmentLength) / segmentLength
	if segment == nPoints-1 {
		return g.Points[nPoints-1]
	}
	return SimpleGradient{g.Points[segment], g.Points[segment+1]}.Interpolate(remainder)
}

func (g LinearGradient) String() string {
	return fmt.Sprintf("LinearGradient %s", g.Points)
}

func flipGradient(g Gradient) Gradient {
	return flippedGradient{g}
}

type flippedGradient struct {
	Gradient
}

func (g flippedGradient) Interpolate(t float64) Color {
	return g.Gradient.Interpolate(1 - t)
}

type sampledGradient struct {
	Gradient
	a float64
	b float64
}

func (g sampledGradient) Interpolate(t float64) Color {
	newT := (t * (g.b - g.a)) + g.a
	return g.Gradient.Interpolate(newT)
}

func Subsample(g Gradient, a, b float64) Gradient {
	if b < a {
		return sampledGradient{
			Gradient: flipGradient(g),
			a:        1 - b,
			b:        1 - a,
		}
	}
	return sampledGradient{
		g,
		a,
		b,
	}
}

func interpolate(t, a, b float64) float64 {
	return (b*t + a*(1-t))
}
