package sampler

import (
	"math"

	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
)

func RotatingLine(thickness float64) DynamicSampler {
	return RotatingStatic(
		HorizontalLine(thickness),
		1,
	)
}

func RotatedLine(thickness float64, angle float64) StaticSampler {
	v1 := geometry.Vector2D{X: 0, Y: 1} // orthogonal vector is rotated 90 degrees from vertical
	v2 := geometry.Vector2D{X: 1, Y: 0} // orthogonal vector is rotated 90 degrees from vertical
	m := geometry.RotateMatrix2D(angle)
	v1, v2 = m.MultVect(v1), m.MultVect(v2)
	return rotatedLine{
		Thickness: thickness,
		v1:        v1,
		v2:        v2,
	}
}

type rotatedLine struct {
	Thickness float64
	v1        geometry.Vector2D
	v2        geometry.Vector2D
}

func (t rotatedLine) GetValue(x, y float64) float64 {
	v := geometry.Vector2D{X: x*2 - 1, Y: y*2 - 1} // do math in the square -1 to 1
	xdistance := math.Abs(v.DotProduct(t.v1))
	ydistance := math.Abs(v.DotProduct(t.v2))
	if xdistance < t.Thickness && ydistance < 1.0 {
		return 1.0
	}
	return 0.0
}

func RotatedCross(thickness, angle float64) StaticSampler {
	v1 := geometry.Vector2D{X: 0, Y: 1} // orthogonal vector is rotated 90 degrees from vertical
	v2 := geometry.Vector2D{X: 1, Y: 0} // orthogonal vector is rotated 90 degrees from vertical
	m := geometry.RotateMatrix2D(angle)
	v1, v2 = m.MultVect(v1), m.MultVect(v2)
	return rotatedCross{
		Thickness: thickness,
		v1:        v1,
		v2:        v2,
	}
}

type rotatedCross struct {
	Thickness float64
	v1        geometry.Vector2D
	v2        geometry.Vector2D
}

func (t rotatedCross) GetValue(x, y float64) float64 {
	v := geometry.Vector2D{X: x*2 - 1, Y: y*2 - 1} // do math in the square -1 to 1
	d1, d2 := math.Abs(v.DotProduct(t.v1)), math.Abs(v.DotProduct(t.v2))
	if (d1 < t.Thickness && d2 < 1.0) || (d2 < t.Thickness && d1 < 1.0) {
		return 1
	}
	if v.Mag() < 2*t.Thickness {
		return 1
	}
	return 0
}

func Cross(thickness float64) StaticSampler {
	return MaxStaticCombiner(
		HorizontalLine(thickness),
		VerticalLine(thickness),
	)

}

func RotatingCross(thickness float64) DynamicSampler {
	return RotatingStatic(Cross(thickness), 1)
}

type PulsingSquare struct{}

func (t PulsingSquare) GetFrameValue(x, y, f float64) float64 {
	x, y = x*2-1, y*2-1 // do math in the square -1 to 1
	if math.Abs(x) < f && math.Abs(y) < f {
		return 1
	}
	return 0
}

func VerticalLine(thickness float64) StaticSampler {
	return verticalLine{thickness}
}

type verticalLine struct {
	Thickness float64
}

func (s verticalLine) GetValue(x, y float64) float64 {
	x, y = x*2-1, y*2-1 // do math in the square -1 to 1
	xdistance := math.Abs(x)
	ydistance := math.Abs(y)
	if xdistance < s.Thickness && ydistance < 1.0 {
		return 1.0
	}
	return 0.0
}

func HorizontalLine(thickness float64) StaticSampler {
	return horizontalLine{thickness}
}

type horizontalLine struct {
	Thickness float64
}

func (s horizontalLine) GetValue(x, y float64) float64 {
	x, y = x*2-1, y*2-1 // do math in the square -1 to 1
	xdistance := math.Abs(x)
	ydistance := math.Abs(y)
	if ydistance < s.Thickness && xdistance < 1.0 {
		return 1.0
	}
	return 0.0
}

type Square struct {
	Thickness float64
}

func (s Square) GetValue(x, y float64) float64 {
	return PulsingSquare{}.GetFrameValue(x, y, s.Thickness)
}

type RoundedSquare struct {
	HalfWidth float64
	Radius    float64 // takes away from Thickness
}

func (s RoundedSquare) GetValue(x, y float64) float64 {
	x, y = x*2-1, y*2-1
	if math.Abs(x) > s.HalfWidth || math.Abs(y) > s.HalfWidth {
		// definitely outside
		return 0
	}
	if math.Abs(x) <= s.HalfWidth-s.Radius && math.Abs(y) <= s.HalfWidth-s.Radius {
		// inside the inner square
		return 1
	}
	if (math.Abs(x) >= s.HalfWidth-s.Radius) != (math.Abs(y) >= s.HalfWidth-s.Radius) {
		// one of the side rectangles
		return 1
	}
	if maths.Radius(s.HalfWidth-s.Radius-math.Abs(x), s.HalfWidth-s.Radius-math.Abs(y)) <= s.Radius {
		return 1
	}
	return 0
}

type Circle struct {
	Thickness float64 // in texture coordinates
}

func (s Circle) GetValue(x, y float64) float64 {
	v := geometry.Vector2D{X: x*2 - 1, Y: y*2 - 1} // do math in the square -1 to 1
	if v.Mag() < 2*s.Thickness {
		return 1
	}
	return 0
}
