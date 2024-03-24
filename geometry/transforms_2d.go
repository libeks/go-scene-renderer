package geometry

import "math"

func ScaleMatrix2D(t float64) Matrix2D {
	if t == 0 {
		panic("Cannot scale by 0.0")
	}
	return Matrix2D{
		t,
		0,

		0,
		t,
	}
}

func RotateMatrix2D(t float64) Matrix2D {
	return Matrix2D{
		math.Cos(t),
		-math.Sin(t),

		math.Sin(t),
		math.Cos(t),
	}
}
