package maths

import "math"

// The best a 1D bezier will give is a slowly growing curve
//  see https://www.desmos.com/calculator/y0z57e8ndr
// Using sigmoid with a high k is preferred, it is much more abrupt

func BezierZeroOne(t float64) float64 {
	return -2*math.Pow(t, 3) + 3*math.Pow(t, 2)
}

type Vec3 struct {
	X, Y, Z float64
}

func (t Vec3) Add(a Vec3) Vec3 {
	return Vec3{
		t.X + a.X,
		t.Y + a.Y,
		t.Z + a.Z,
	}
}

func (t Vec3) ScalarMultiply(a float64) Vec3 {
	return Vec3{
		t.X * a,
		t.Y * a,
		t.Z * a,
	}
}

type Linear float64

func (l Linear) Add(b Linear) Linear {
	return Linear(l + b)
}

func (l Linear) ScalarMultiply(a float64) Linear {
	return Linear(a * float64(l))
}

// TODO generalize this to 2d, 3d points, colors
// taken from https://stackoverflow.com/a/73851453
type Pointlike interface {
	Vec3 | Linear
}

type point[T any] interface {
	Pointlike
	Add(T) T
	ScalarMultiply(float64) T
}

type Bezier[T point[T]] struct {
	points []T
}

// t ranges from 0.0 to 1.0
func (b Bezier[T]) At(t float64) T {
	if len(b.points) == 1 {
		return b.points[0]
	}
	return Bezier[T]{
		b.points[0 : len(b.points)-1],
	}.At(t).ScalarMultiply(1 - t).Add(Bezier[T]{
		b.points[1:len(b.points)],
	}.At(t).ScalarMultiply(t))
}
