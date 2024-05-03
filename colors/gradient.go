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

// // Subsample the original gradient to produce a new one
// func (g SimpleGradient) Subsample(a, b float64) Gradient {
// 	return SimpleGradient{
// 		g.Interpolate(a),
// 		g.Interpolate(b),
// 	}
// }

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

// func (g LinearGradient) Subsample(a, b float64) Gradient {

// }

// type ColorDot struct {
// 	t float64
// 	c Color
// }

// type ParameterizedGradient struct {
// 	// assume that colors[0] with have t=0 and colors[-1] will have t=1
// 	// colorA Color      // startpoint
// 	colors []ColorDot // midpoints with parameterized t, sorted
// 	// colorB Color      // endpoint
// }

// func (g ParameterizedGradient) Interpolate(v float64) Color {
// 	if len(g.colors) < 2 {
// 		panic(fmt.Errorf("parameterized gradient %s should have at least two colors", g.colors))
// 	}
// 	if v <= 0 {
// 		return g.colors[0].c
// 	}
// 	if v >= 1 {
// 		return g.colors[len(g.colors)-1].c
// 	}
// 	index := g.getIndex(v)
// 	// interpolate between colors[index] to colors[index+1]
// 	t1, t2 := c.colors[index].t, c.colors[index+1].t
// 	tNew := (v - t1) / (t2 - t1)
// 	if tNew <= 0 || tNew > 1 {
// 		panic(fmt.Errorf("parameterized gradient value at v=%.3f results in t=%0.3f, %s", v, tNew, g.colors))
// 	}
// 	return SimpleGradient{c.colors[index].c, c.colors[index+1].c}.Interpolate(tNew)
// }

// func (g ParameterizedGradient) getIndex(v float64) int {
// 	for i, c := range g.colors {
// 		if v >= c.t {
// 			return i
// 		}
// 	}
// }

// func (g ParameterizedGradient) interpolateT(v float64, index int) float64 {
// 	// interpolate between colors[index] to colors[index+1]
// 	t1, t2 := c.colors[index].t, c.colors[index+1].t
// 	return (v - t1) / (t2 - t1)
// }

// func (g ParameterizedGradient) Subsample(a, b float64) Gradient {
// 	var colors []ColorDot
// 	start := g.getIndex(a)
// 	end := g.getIndex(b)
// 	if start == end {
// 		return SimpleGradient{
// 			Start: g.Interpolate(a),
// 			End:   g.Interpolate(b),
// 		}
// 	}
// 	colors = []ColorDot{
// 		ColorDot{
// 			c: g.Interpolate(a),
// 			t: 0,
// 		},
// 	}
// 	width := b - a
// 	for i := start; i < end; i++ {
// 		t1, t2 := c.colors[index].t, c.colors[index+1].t

// 		colors = append(colors, ColorDot{
// 			c: g.colors[i+1].c,
// 			t: (g.colors[i+1].t / (b - a)),
// 		})
// 	}
// 	colors = append(colors,
// 		ColorDot{
// 			c: g.Interpolate(b),
// 			t: 1,
// 		},
// 	)
// }

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
