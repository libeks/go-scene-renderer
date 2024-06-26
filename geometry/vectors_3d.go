package geometry

import (
	"fmt"
	"math"
)

var (
	NilVector3D = Vector3D{} // default zero values
)

// helper to not have to write field names
func V3(x, y, z float64) Vector3D {
	return Vector3D{
		X: x,
		Y: y,
		Z: z,
	}
}

type Vector3D struct {
	X, Y, Z float64
}

func (v Vector3D) String() string {
	return fmt.Sprintf("V[%0.3f, %0.3f, %0.3f]", v.X, v.Y, v.Z)
}

func (v Vector3D) ToHomogenous() HomogenousVector {
	return HomogenousVector{
		v.X, v.Y, v.Z, 1,
	}
}

func (v Vector3D) AddVector(w Vector3D) Vector3D {
	return Vector3D{
		v.X + w.X,
		v.Y + w.Y,
		v.Z + w.Z,
	}
}

func (v Vector3D) ScalarMultiply(r float64) Vector3D {
	return Vector3D{
		v.X * r,
		v.Y * r,
		v.Z * r,
	}
}

func (v Vector3D) DotProduct(w Vector3D) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vector3D) ScalarProject(w Vector3D) float64 {
	return v.DotProduct(w) / (w.Mag() * w.Mag())
}

func (v Vector3D) CrossProduct(w Vector3D) Vector3D {
	return Vector3D{
		v.Y*w.Z - v.Z*w.Y,
		v.Z*w.X - v.X*w.Z,
		v.X*w.Y - v.Y*w.X,
	}
}

func (v Vector3D) Mag() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v Vector3D) Unit() Vector3D {
	d := v.Mag()
	if d == 0 {
		return NilVector3D
	}
	return v.ScalarMultiply(1 / d)
}
