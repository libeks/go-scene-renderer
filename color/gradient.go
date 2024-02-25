package color

type Gradient struct {
	Start Color
	End   Color
}

// v=0.0 returns start color
// v=1.0 return end color
func (g Gradient) Interpolate(v float64) Color {
	return Color{
		R: interpolate(v, g.Start.R, g.End.R),
		G: interpolate(v, g.Start.G, g.End.G),
		B: interpolate(v, g.Start.B, g.End.B),
	}
}

func interpolate(t, a, b float64) float64 {
	return (b*t + a*(1-t))
}
