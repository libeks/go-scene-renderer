package geometry

import "fmt"

var (
	OriginPoint = Point{}
)

func Pt(x, y, z float64) Point {
	return Point{X: x, Y: y, Z: z}
}

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

func (p Point) IsInFrontOfCamera(minDepth float64) bool {
	return p.Z < minDepth
}

// return Screen pixel, as well as the z-depth (not actual depth)
func (p Point) ToPixel() (*Pixel, float64) {
	if p.Z > 0 {
		return nil, 0
	}
	return &Pixel{
		p.X / -p.Z,
		p.Y / -p.Z,
	}, -p.Z
}

func (p Point) LinesAroundPoint() []Line {
	uD := V3(0, .2, 0)
	rD := V3(.2, 0, 0)
	bD := V3(0, 0, .2)

	return []Line{
		{A: p, B: Point(p.Vector().AddVector(uD))},
		{A: p, B: Point(p.Vector().AddVector(rD))},
		{A: p, B: Point(p.Vector().AddVector(bD))},
	}
}
