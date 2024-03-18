package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
)

type TriangleGradientTexture struct {
	ColorA color.Color
	ColorB color.Color
	ColorC color.Color
}

// given coordinates from the A point towards B and C (each in the range of (0,1))
// return what color it should be
func (t *TriangleGradientTexture) GetTextureColor(b, c float64) color.Color {
	if c == 1 {
		return t.ColorB
	}
	abGradient := color.SimpleGradient{t.ColorA, t.ColorB}
	abColor := abGradient.Interpolate(b / (1 - c))
	triangleGradient := color.SimpleGradient{abColor, t.ColorC}
	cColor := triangleGradient.Interpolate(c)
	return cColor
}

func GradientTriangle(a, b, c geometry.Point, colorA, colorB, colorC color.Color) Triangle {
	return Triangle{
		A: a,
		B: b,
		C: c,
		Colorer: &TriangleGradientTexture{
			ColorA: colorA,
			ColorB: colorB,
			ColorC: colorC,
		},
	}
}

type Triangle struct {
	A geometry.Point
	B geometry.Point
	C geometry.Point
	// Colorer will be evaluated with two parameters (b,c), each from (0,1), but a+b<1.0
	// it describes the coordinates on the triangle from A towards B and C, respectively
	Colorer color.Texture
}

func (t Triangle) GetColorDepth(x, y float64) (*color.Color, float64) {
	b, c, depth, intersect := t.rayIntersectLocalCoords(Ray{geometry.OriginPoint, geometry.Vector3D{x, y, -1.0}})
	if !intersect {
		return nil, 0
	}
	color := t.Colorer.GetTextureColor(b, c)
	return &color, depth
}

func (t Triangle) ApplyMatrix(m geometry.HomogeneusMatrix) TransformableObject {
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
		t.Colorer,
	}
}

func (t Triangle) GetWireframe() []geometry.Line {
	return []geometry.Line{
		geometry.Line{t.A, t.B},
		geometry.Line{t.A, t.C},
		geometry.Line{t.B, t.C},
	}
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle %s %s %s", t.A, t.B, t.C)
}

func (t Triangle) plane() Plane {
	bVector := t.bVect()
	cVector := t.cVect()
	nVector := bVector.CrossProduct(cVector)
	return Plane{nVector, t.A.Vector().DotProduct(nVector)}
}

func (t Triangle) bVect() geometry.Vector3D {
	return t.B.Subtract(t.A)
}

func (t Triangle) cVect() geometry.Vector3D {
	return t.C.Subtract(t.A)
}

// return the intersection in triangle-local coordinates, in direction of A->B and A->C
// bool signifies whether intersection is inside the triange
// third float is the depth, in positive values
func (t Triangle) rayIntersectLocalCoords(r Ray) (float64, float64, float64, bool) {
	intersectDot := t.plane().IntersectPoint(r)
	if intersectDot == nil {
		return 0, 0, 0, false
	}
	iVect := intersectDot.Subtract(t.A)
	iMag := geometry.OriginPoint.Subtract(*intersectDot).Mag()

	bVect := t.bVect()
	cVect := t.cVect()
	b := iVect.ScalarProject(bVect)
	c := iVect.ScalarProject(cVect)
	// check if vector (b,c) is inside the triangle [(0,0), (1,0), (0,1)]
	if b < 0.0 || b > 1.0 || c < 0.0 || c > 1.0 {
		// outside the unit square
		return b, c, iMag, false
	}
	if b+c > 1.0 {
		// fmt.Printf("Outside unit hypotenuse\n")
		// inside unit square, but on far side of hypotenuse
		return b, c, iMag, false
	}
	// inside unit square and inside the hypotenuse
	return b, c, iMag, true
}

type Ray struct {
	P geometry.Point    // origin point
	D geometry.Vector3D // direction vector describing the ray
}

type Plane struct {
	N geometry.Vector3D // normal vector
	D float64           // d parameter, describing plane
}

func (p Plane) String() string {
	return fmt.Sprintf("Plane(->%s, at %f)", p.N, p.D)
}

func (p Plane) IntersectPoint(r Ray) *geometry.Point {
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
	point := geometry.Point(r.P.Vector().AddVector(r.D.ScalarMultiply(t)))
	return &point
}
