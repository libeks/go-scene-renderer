package geometry

import (
	"fmt"
	"math"
)

var (
	NilVector2D = Vector2D{} // default zero values
)

type Vector2D struct {
	X, Y float64
}

func (v Vector2D) String() string {
	return fmt.Sprintf("V[%0.3f, %0.3f]", v.X, v.Y)
}

// func (v Vector2D) ToHomogenous() HomogenousVector {
// 	return HomogenousVector{
// 		v.X, v.Y, v.Z, 1,
// 	}
// }

func (v Vector2D) AddVector(w Vector2D) Vector2D {
	return Vector2D{
		v.X + w.X,
		v.Y + w.Y,
	}
}

func (v Vector2D) ScalarMultiply(r float64) Vector2D {
	return Vector2D{
		v.X * r,
		v.Y * r,
	}
}

func (v Vector2D) DotProduct(w Vector2D) float64 {
	return v.X*w.X + v.Y*w.Y
}

// func (v Vector2D) ScalarProject(w Vector2D) float64 {
// 	return v.DotProduct(w) / (w.Mag() * w.Mag())
// }

// func (v Vector2D) CrossProduct(w Vector2D) Vector2D {
// 	return Vector2D{
// 		v.Y*w.Z - v.Z*w.Y,
// 		v.Z*w.X - v.X*w.Z,
// 		v.X*w.Y - v.Y*w.X,
// 	}
// }

func (v Vector2D) Mag() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v Vector2D) Unit() Vector2D {
	d := v.Mag()
	if d == 0 {
		return NilVector2D
	}
	return v.ScalarMultiply(1 / d)
}
