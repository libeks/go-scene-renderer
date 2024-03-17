package geometry

import "fmt"

var (
	OriginPoint = Point{0.0, 0.0, 0.0}
)

// A point is a vector, but I don't want to get confused
type Point Vector3D

func (p Point) String() string {
	return fmt.Sprintf("P(%.3f,%.3f,%.3f)", p.X, p.Y, p.Z)
}

func (p Point) Subtract(q Point) Vector3D {
	return Vector3D{
		p.X - q.X,
		p.Y - q.Y,
		p.Z - q.Z,
	}
}

func (p Point) Add(q Point) Vector3D {
	return Vector3D{
		p.X + q.X,
		p.Y + q.Y,
		p.Z + q.Z,
	}
}

func (p Point) Vector() Vector3D {
	return Vector3D(p)
}

func (p Point) ToHomogenous() HomogenousVector {
	return HomogenousVector{p.X, p.Y, p.Z, 1}
}
