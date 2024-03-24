package geometry

import (
	"fmt"
	"math"
)

var (
	NilHomogenousVector = HomogenousVector{} // default zero values
)

type HomogenousVector struct {
	X, Y, Z, T float64
}

func (v HomogenousVector) String() string {
	return fmt.Sprintf("V[%0.3f, %0.3f, %0.3f, %0.3f]", v.X, v.Y, v.Z, v.T)
}

func (v HomogenousVector) To3D() (Vector3D, bool) {
	if v.T == 0 {
		return Vector3D{}, false
	}
	return Vector3D{
		// normalize by 4th coordinate
		v.X / v.T, v.Y / v.T, v.Z / v.T,
	}, true
}

func (v HomogenousVector) ToPoint() (Point, bool) {
	vect, ok := v.To3D()
	return Point(vect), ok
}

func (v HomogenousVector) AddVector(w HomogenousVector) HomogenousVector {
	// is this actually correct? What is the sum of two
	return HomogenousVector{
		v.X + w.X,
		v.Y + w.Y,
		v.Z + w.Z,
		v.T + w.T,
	}
}

func (v HomogenousVector) ScalarMultiply(r float64) HomogenousVector {
	return HomogenousVector{
		v.X * r,
		v.Y * r,
		v.Z * r,
		v.T, // don't change 4th coordinate
	}
}

func (v HomogenousVector) DotProduct(w HomogenousVector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z + v.T*w.T
}

func (v HomogenousVector) ScalarProject(w HomogenousVector) float64 {
	return v.DotProduct(w) / (w.Mag() * w.Mag())
}

func (v HomogenousVector) Mag() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v HomogenousVector) Unit() HomogenousVector {
	d := v.Mag()
	if d == 0 {
		return NilHomogenousVector
	}
	return v.ScalarMultiply(1 / d)
}
