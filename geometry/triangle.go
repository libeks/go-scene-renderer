package geometry

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/color"
)

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
	return Vector3D{p.X, p.Y, p.Z}
}

func (p Point) ToHomogenous() HomogenousVector {
	return HomogenousVector{p.X, p.Y, p.Z, 1}
}

type Triangle struct {
	A      Point
	B      Point
	C      Point
	ColorA color.Color
	ColorB color.Color
	ColorC color.Color
}

func (t Triangle) GetColorDepth(x, y float64) (*color.Color, float64) {
	b, c, depth, intersect := t.rayIntersectLocalCoords(Ray{OriginPoint, Vector3D{x, y, -1.0}})
	if !intersect {
		return nil, 0
	}
	abGradient := color.SimpleGradient{t.ColorA, t.ColorB}
	abColor := abGradient.Interpolate(b)
	triangleGradient := color.SimpleGradient{abColor, t.ColorC}
	cColor := triangleGradient.Interpolate(c)
	// fmt.Printf("Depth %0.3f,%0.3f at %0.3f\n", x, y, depth)
	return &cColor, depth
}

func (t Triangle) ApplyMatrix(m HomogeneusMatrix) Triangle {
	a, ok := m.MultVect(t.A.ToHomogenous()).ToPoint()
	if !ok {
		return Triangle{}
	}
	b, ok := m.MultVect(t.B.ToHomogenous()).ToPoint()
	if !ok {
		return Triangle{}
	}
	c, ok := m.MultVect(t.C.ToHomogenous()).ToPoint()
	if !ok {
		return Triangle{}
	}
	return Triangle{
		a, b, c,
		t.ColorA, t.ColorB, t.ColorC,
	}
}

func (t Triangle) plane() Plane {
	bVector := t.bVect()
	cVector := t.cVect()
	// fmt.Printf("BVector %s\n", bVector)
	// fmt.Printf("CVector %s\n", cVector)
	nVector := bVector.CrossProduct(cVector)
	return Plane{nVector, t.A.Vector().DotProduct(nVector)}
}

func (t Triangle) bVect() Vector3D {
	return t.B.Subtract(t.A)
}

func (t Triangle) cVect() Vector3D {
	return t.C.Subtract(t.A)
}

// return the intersection in triangle-local coordinates, in direction of A->B and A->C
// bool signifies whether intersection is inside the triange
// third float is the depth, in positive values
func (t Triangle) rayIntersectLocalCoords(r Ray) (float64, float64, float64, bool) {
	// fmt.Printf("plane %s\n", t.Plane())
	// fmt.Printf("ray %s\n", r)
	intersectDot := t.plane().IntersectPoint(r)
	// fmt.Printf("intersection dot %s\n", intersectDot)
	if intersectDot == nil {
		return 0, 0, 0, false
	}
	iVect := intersectDot.Subtract(t.A)
	iMag := OriginPoint.Subtract(*intersectDot).Mag()
	// fmt.Printf("intersection vector %s\n", iVect)

	bVect := t.bVect()
	cVect := t.cVect()
	// fmt.Printf("bVect %s\n", bVect)
	b := iVect.ScalarProject(bVect)
	c := iVect.ScalarProject(cVect)
	// fmt.Printf("plane coords %0.3f %0.3f\n", b, c)
	// check if vector (b,c) is inside the triangle [(0,0), (1,0), (0,1)]
	if b < 0.0 || b > 1.0 || c < 0.0 || c > 1.0 {
		// outside the unit square
		// fmt.Printf("Outside unit square\n")
		return b, c, iMag, false
	}
	if b+c > 1.0 {
		// fmt.Printf("Outside unit hypotenuse\n")
		// inside unit square, but on far side of hypotenuse
		return b, c, iMag, false
	}
	// inside unit square and inside the hypotenuse
	// fmt.Printf("inside triangle\n")
	return b, c, iMag, true
}

type Ray struct {
	P Point    // origin point
	D Vector3D // direction vector describing the ray
}

type Plane struct {
	N Vector3D // normal vector
	D float64  // d parameter, describing plane
}

func (p Plane) String() string {
	return fmt.Sprintf("Plane(->%s, at %f)", p.N, p.D)
}

func (p Plane) IntersectPoint(r Ray) *Point {
	denominator := p.N.DotProduct(r.D)
	// fmt.Printf("denominator %0.3f\n", denominator)
	if denominator == 0.0 {
		return nil // ray is parallel to plane, no intersection
	}
	t := (p.D - p.N.DotProduct(r.P.Vector())) / denominator
	if t < 0.0 {
		// fmt.Printf("object behind camera %0.3f\n", t)
		return nil // ray intersects plane before ray's starting point
	}
	point := Point(r.P.Vector().AddVector(r.D.ScalarMultiply(t)))
	return &point
}
