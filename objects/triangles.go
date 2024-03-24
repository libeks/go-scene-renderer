package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
)

func GradientTriangle(a, b, c geometry.Point, colorA, colorB, colorC colors.Color) DynamicTriangle {
	return DynamicTriangle{
		Triangle: Triangle{
			A: a,
			B: b,
			C: c,
		},
		Colorer: colors.StaticTexture(colors.TriangleGradientTexture(colorA, colorB, colorC)),
	}
}

// DynamicTriangle is a Triangle with a DynamicTexture, which can be evaluated for a specific frame
type DynamicTriangle struct {
	Triangle
	Colorer colors.DynamicTexture
}

func (t DynamicTriangle) Frame(f float64) StaticTriangle {
	return StaticTriangle{
		t.Triangle,
		t.Colorer.GetFrame(f),
	}
}

func (t DynamicTriangle) ApplyMatrix(m geometry.HomogeneusMatrix) *DynamicTriangle {
	newTriangle := t.Triangle.ApplyMatrix(m)
	if newTriangle == nil {
		return nil
	}
	return &DynamicTriangle{
		Triangle: *newTriangle,
		Colorer:  t.Colorer,
	}
}

func (t DynamicTriangle) GetBoundingBox() BoundingBox {
	return t.Triangle.GetBoundingBox()
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
func (t DynamicTriangle) GetWireframe() []geometry.Line {
	return t.Triangle.GetWireframe()
}

// StaticTriangle is a Triangle with a Texture applied to it
type StaticTriangle struct {
	Triangle
	// Colorer will be evaluated with two parameters (b,c), each from (0,1), but b+c<1.0
	// it describes the coordinates on the triangle from A towards B and C, respectively
	Colorer colors.Texture
}

// returns the color of the triangle at a ray
// emanating from the camera at (0,0,0), pointed in the direction
// (x,y, -1), with perspective
// and a z-index. The bigger the index, the farther the object.
func (t *StaticTriangle) GetColorDepth(x, y float64) (*colors.Color, float64) {
	b, c, depth, intersect := t.rayIntersectLocalCoords(ray{geometry.OriginPoint, geometry.Vector3D{x, y, -1.0}})
	if !intersect {
		return nil, 0
	}
	color := t.Colorer.GetTextureColor(b, c)
	return &color, depth
}

func (t StaticTriangle) GetBoundingBox() BoundingBox {
	return t.Triangle.GetBoundingBox()
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
func (t StaticTriangle) GetWireframe() []geometry.Line {
	return t.Triangle.GetWireframe()
}

// A Triangle describes an uncolored object in the space
type Triangle struct {
	A geometry.Point
	B geometry.Point
	C geometry.Point

	// the below are cached values for efficiency. They are created at the top of rayIntersectLocalCoords
	cached bool
	plane  plane
	bVect  geometry.Vector3D
	cVect  geometry.Vector3D
	normal geometry.Vector3D

	cachedBoundingBox bool
	bbox              BoundingBox
}

func (t Triangle) ApplyMatrix(m geometry.HomogeneusMatrix) *Triangle {
	a, ok := m.MultVect(t.A.ToHomogenous()).ToPoint()
	if !ok {
		return nil
	}
	b, ok := m.MultVect(t.B.ToHomogenous()).ToPoint()
	if !ok {
		return nil
	}
	c, ok := m.MultVect(t.C.ToHomogenous()).ToPoint()
	if !ok {
		return nil
	}
	return &Triangle{
		A: a, B: b, C: c,
	}
}

func (t Triangle) Flatten() []*Triangle {
	return []*Triangle{&t}
}

func (t *Triangle) GetBoundingBox() BoundingBox {
	// TODO: cache the bounding box?
	// fmt.Printf("%s\n", t)
	if t.cachedBoundingBox {
		return t.bbox
	}
	a, ad := t.A.ToPixel()
	b, bd := t.B.ToPixel()
	c, cd := t.C.ToPixel()
	// fmt.Printf("pixels %s %s %s\n", a, b, c)
	if a == nil || b == nil || c == nil {
		return BoundingBox{
			TopLeft: geometry.Pixel{
				-2, -2,
			},
			BottomRight: geometry.Pixel{
				2, 2,
			},
		}
	}
	minx := min(a.X, b.X, c.X)
	miny := min(a.Y, b.Y, c.Y)
	maxx := max(a.X, b.X, c.X)
	maxy := max(a.Y, b.Y, c.Y)
	mindepth := min(ad, bd, cd)
	maxdepth := max(ad, bd, cd)
	bb := BoundingBox{
		TopLeft: geometry.Pixel{
			minx, miny,
		},
		BottomRight: geometry.Pixel{
			maxx, maxy,
		},
		MinDepth: mindepth,
		MaxDepth: maxdepth,
	}
	t.bbox = bb
	// t.cachedBoundingBox = true
	// fmt.Printf("bounding box %s\n", bb)
	return bb
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
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

func (t Triangle) getPlane() plane {
	bVector := t.bVect
	cVector := t.cVect
	nVector := bVector.CrossProduct(cVector)
	return plane{nVector, t.A.Vector().DotProduct(nVector)}
}

func (t Triangle) getBVect() geometry.Vector3D {
	// fmt.Printf("BVect %s\n", t.B.Subtract(t.A))
	return t.B.Subtract(t.A)
}

func (t Triangle) getCVect() geometry.Vector3D {
	// fmt.Printf("CVect %s\n", t.C.Subtract(t.A))
	return t.C.Subtract(t.A)
}

// return the intersection in triangle-local coordinates, in direction of A->B and A->C
// bool signifies whether intersection is inside the triange
// third float is the depth, in positive values
func (t *Triangle) rayIntersectLocalCoords(r ray) (float64, float64, float64, bool) {
	// cache the vectors AB and AC, as well as the plane, this is 37% more efficient
	if !t.cached {
		t.bVect = t.getBVect()
		t.cVect = t.getCVect()
		t.plane = t.getPlane()
		t.normal = t.plane.N
		t.cached = true
	}
	intersectDot := t.plane.IntersectPoint(r)
	if intersectDot == nil {
		return 0, 0, 0, false
	}
	iVect := intersectDot.Subtract(t.A)
	iMag := geometry.OriginPoint.Subtract(*intersectDot).Mag()

	bVect := t.bVect
	cVect := t.cVect
	normal := t.normal
	normalMag := normal.Mag()
	b := iVect.CrossProduct(cVect).DotProduct(normal) / (normalMag * normalMag)
	c := bVect.CrossProduct(iVect).DotProduct(normal) / (normalMag * normalMag)
	// check if vector (b,c) is inside the triangle [(0,0), (1,0), (0,1)]
	if b < 0.0 || b > 1.0 || c < 0.0 || c > 1.0 {
		// outside the unit square
		return b, c, iMag, false
	}
	if b+c > 1.0 {
		// inside unit square, but on far side of hypotenuse
		return b, c, iMag, false
	}
	// inside unit square and inside the hypotenuse
	return b, c, iMag, true
}
