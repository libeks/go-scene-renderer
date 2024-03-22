package geometry

import "math"

func TranslationMatrix(v Vector3D) HomogeneusMatrix {
	return HomogeneusMatrix{
		1.0,
		0.0,
		0.0,
		v.X,

		0.0,
		1.0,
		0.0,
		v.Y,

		0.0,
		0.0,
		1.0,
		v.Z,

		0.0,
		0.0,
		0.0,
		1.0,
	}
}

func ScaleMatrix(t float64) HomogeneusMatrix {
	if t == 0.0 {
		panic("Cannot scale by 0.0")
	}
	return HomogeneusMatrix{
		1.0,
		0.0,
		0.0,
		0.0,

		0.0,
		1.0,
		0.0,
		0.0,

		0.0,
		0.0,
		1.0,
		0.0,

		0.0,
		0.0,
		0.0,
		1.0 / t,
	}
}

func RotateMatrixX(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		1,
		0,
		0,
		0,

		0,
		math.Cos(t),
		-math.Sin(t),
		0,

		0,
		math.Sin(t),
		math.Cos(t),
		0,

		0,
		0,
		0,
		1,
	}
}

func RotateMatrixY(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		math.Cos(t),
		0,
		math.Sin(t),
		0,

		0,
		1,
		0,
		0,

		-math.Sin(t),
		0,
		math.Cos(t),
		0,

		0,
		0,
		0,
		1,
	}
}

func RotateMatrixZ(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		math.Cos(t),
		-math.Sin(t),
		0,
		0,

		math.Sin(t),
		math.Cos(t),
		0,
		0,

		0,
		0,
		1,
		0,

		0,
		0,
		0,
		1,
	}
}
