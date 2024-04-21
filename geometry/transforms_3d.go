package geometry

import "math"

func Scale2DMatrix(t float64) Matrix3D {
	if t == 0 {
		panic("Cannot scale by 0.0")
	}
	return Matrix3D{
		t,
		0,
		0,

		0,
		t,
		0,

		0,
		0,
		t,
	}
}

func RotatePitch3D(t float64) Matrix3D {
	return Matrix3D{
		1,
		0,
		0,

		0,
		math.Cos(t),
		math.Sin(t),

		0,
		-math.Sin(t),
		math.Cos(t),
	}
}

func RotateYaw3D(t float64) Matrix3D {
	return Matrix3D{
		math.Cos(t),
		0,
		math.Sin(t),

		0,
		1,
		0,

		-math.Sin(t),
		0,
		math.Cos(t),
	}
}

func RotateRoll3D(t float64) Matrix3D {
	return Matrix3D{
		math.Cos(t),
		-math.Sin(t),
		0,

		math.Sin(t),
		math.Cos(t),
		0,

		0,
		0,
		1,
	}
}
