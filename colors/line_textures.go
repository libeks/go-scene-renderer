package colors

import (
	"math"

	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
)

type RotatingLine struct {
	On        Color
	Off       Color
	Thickness float64 // in texture coordinates
}

func (t RotatingLine) GetFrameColor(x, y, f float64) Color {
	v := geometry.Vector2D{x*2 - 1, y*2 - 1} // do math in the square -1 to 1
	orth := geometry.Vector2D{1, 0}          // orthogonal vector is rotated 90 degrees from vertical
	direct := geometry.Vector2D{0, 1}
	m := geometry.RotateMatrix2D(f * maths.Rotation)
	orth, direct = m.MultVect(orth), m.MultVect(direct)
	xdistance := math.Abs(v.DotProduct(orth))
	ydistance := math.Abs(v.DotProduct(direct))
	if xdistance < t.Thickness && ydistance < 1.0 {
		return t.On
	}
	return t.Off
}

type RotatingCross struct {
	On        Color
	Off       Color
	Thickness float64 // in texture coordinates
}

func (t RotatingCross) GetFrameColor(x, y, f float64) Color {
	v := geometry.Vector2D{x*2 - 1, y*2 - 1} // do math in the square -1 to 1
	v1 := geometry.Vector2D{0, 1}            // orthogonal vector is rotated 90 degrees from vertical
	v2 := geometry.Vector2D{1, 0}            // orthogonal vector is rotated 90 degrees from vertical
	m := geometry.RotateMatrix2D(f * maths.Rotation)
	v1, v2 = m.MultVect(v1), m.MultVect(v2)
	d1, d2 := math.Abs(v.DotProduct(v1)), math.Abs(v.DotProduct(v2))
	// distance := min(math.Abs(v.DotProduct(v1)), math.Abs(v.DotProduct(v2)))
	if (d1 < t.Thickness && d2 < 1.0) || (d2 < t.Thickness && d1 < 1.0) {
		return t.On
	}
	if v.Mag() < 2*t.Thickness {
		return t.On
	}
	return t.Off
}

type PulsingSquare struct {
	On  Color
	Off Color
}

func (t PulsingSquare) GetFrameColor(x, y, f float64) Color {
	x, y = x*2-1, y*2-1 // do math in the square -1 to 1
	if math.Abs(x) < f && math.Abs(y) < f {
		return t.On
	}
	return t.Off
}

type VerticalLine struct {
	On        Color
	Off       Color
	Thickness float64
}

func (t VerticalLine) GetTextureColor(x, y float64) Color {
	return RotatingLine(t).GetFrameColor(x, y, 0)
}

type HorizontalLine struct {
	On        Color
	Off       Color
	Thickness float64
}

func (t HorizontalLine) GetTextureColor(x, y float64) Color {
	return RotatingLine(t).GetFrameColor(x, y, .25)
}

type Cross struct {
	On        Color
	Off       Color
	Thickness float64
}

func (t Cross) GetTextureColor(x, y float64) Color {
	return RotatingCross(t).GetFrameColor(x, y, 0)
}

type Square struct {
	On        Color
	Off       Color
	Thickness float64
}

func (t Square) GetTextureColor(x, y float64) Color {
	return PulsingSquare{t.On, t.Off}.GetFrameColor(x, y, t.Thickness)
}

type Circle struct {
	On        Color
	Off       Color
	Thickness float64 // in texture coordinates
}

func (t Circle) GetTextureColor(x, y float64) Color {
	v := geometry.Vector2D{x*2 - 1, y*2 - 1} // do math in the square -1 to 1
	if v.Mag() < 2*t.Thickness {
		return t.On
	}
	return t.Off
}
