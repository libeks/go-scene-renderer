package geometry

import (
	"math"
)

func TranslationMatrix(v Vector3D) HomogeneusMatrix {
	return HomogeneusMatrix{
		1,
		0,
		0,
		v.X,

		0,
		1,
		0,
		v.Y,

		0,
		0,
		1,
		v.Z,

		0,
		0,
		0,
		1,
	}
}

func ScaleMatrix(t float64) HomogeneusMatrix {
	if t == 0 {
		panic("Cannot scale by 0.0")
	}
	return HomogeneusMatrix{
		t,
		0,
		0,
		0,

		0,
		t,
		0,
		0,

		0,
		0,
		t,
		0,

		0,
		0,
		0,
		1,
	}
}

func RotateYaw(t float64) HomogeneusMatrix {
	return RotateMatrixY(t)
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

func RotatePitch(t float64) HomogeneusMatrix {
	return RotateMatrixX(t)
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

func RotateRoll(t float64) HomogeneusMatrix {
	return RotateMatrixZ(t)
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

func PointTowards(p Point) HomogeneusMatrix {
	// make it so the -z vector points towards point p
	// assume the object is at 0,0,0
	forward := p.Vector().Unit()
	right := forward.CrossProduct(V3(0, 1, 0)).Unit()
	up := right.CrossProduct(forward)
	m := Matrix3D{
		right.X,
		up.X,
		-forward.X, // negative since the camera points in negative z-direction

		right.Y,
		up.Y,
		-forward.Y, // negative since the camera points in negative z-direction

		right.Z,
		up.Z,
		-forward.Z, // negative since the camera points in negative z-direction
	}.toHomogenous()
	// fmt.Printf("point towards %s matrix : %s, det %.3f\n", forward, m, m.Determinant())
	return m
}
