package objects

import (
	"fmt"
	"math"
	"slices"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/textures"
)

func GradientTriangle(a, b, c geometry.Point, colorA, colorB, colorC colors.Color) dynamicBasicObject {
	return DynamicBasicObject(
		&Triangle{
			A: a,
			B: b,
			C: c,
		},
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.TriangleGradientTexture(colorA, colorB, colorC))),
	)
}

func Tri(a, b, c geometry.Point) *Triangle {
	return &Triangle{
		A: a,
		B: b,
		C: c,
	}
}

// A Triangle describes an uncolored object in the space
type Triangle struct {
	A geometry.Point
	B geometry.Point
	C geometry.Point

	// the below are cached values for efficiency. They are created at the top of rayIntersectLocalCoords
	cached      bool
	plane       plane
	bVect       geometry.Vector3D
	cVect       geometry.Vector3D
	normal      geometry.Vector3D
	normalMagSq float64

	cachedBoundingBox bool
	bbox              BoundingBox
}

func (t Triangle) ApplyMatrix(m geometry.HomogeneusMatrix) BasicObject {
	a, ok := m.MultVect(t.A.ToHomogenous()).ToPoint()
	if !ok {
		panic(fmt.Errorf("could not apply matrix %s to point %s", m, t.A))
	}
	b, ok := m.MultVect(t.B.ToHomogenous()).ToPoint()
	if !ok {
		panic(fmt.Errorf("could not apply matrix %s to point %s", m, t.B))
	}
	c, ok := m.MultVect(t.C.ToHomogenous()).ToPoint()
	if !ok {
		panic(fmt.Errorf("could not apply matrix %s to point %s", m, t.C))
	}
	return &Triangle{
		A: a, B: b, C: c,
	}
}

func (t Triangle) Flatten() []*Triangle {
	return []*Triangle{&t}
}

func (t *Triangle) GetBoundingBox() BoundingBox {
	if t.cachedBoundingBox {
		return t.bbox
	}
	wireframe := t.getSceneWireframe()
	points := make([]geometry.Pixel, 0, 6)
	pointsX := make([]float64, 0, 6)
	pointsY := make([]float64, 0, 6)
	zdepths := make([]float64, 0, 6)
	for _, line := range wireframe {
		pointA, zdepthA := line.A.ToPixel()
		pointB, zdepthB := line.B.ToPixel()
		if pointA == nil || pointB == nil {
			panic(fmt.Errorf("line should already be in front of camera %s", line))
		}
		points = append(points, *pointA)
		pointsX = append(pointsX, pointA.X)
		pointsY = append(pointsY, pointA.Y)
		points = append(points, *pointB)
		pointsX = append(pointsX, pointB.X)
		pointsY = append(pointsY, pointB.Y)
		zdepths = append(zdepths, zdepthA, zdepthB)
	}

	if len(points) == 0 {
		return BoundingBox{
			empty: true,
		}
	}

	bb := BoundingBox{
		TopLeft: geometry.Pixel{
			X: max(slices.Min(pointsX), -1.0),
			Y: max(slices.Min(pointsY), -1.0),
		},
		BottomRight: geometry.Pixel{
			X: min(slices.Max(pointsX), 1.0),
			Y: min(slices.Max(pointsY), 1.0),
		},
		MinZDepth: max(0, slices.Min(zdepths)),
		MaxZDepth: min(math.MaxFloat64, slices.Max(zdepths)),
	}
	t.bbox = bb
	t.cachedBoundingBox = true
	return bb
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
// note that the resulting lines may not exactly match the triangle, as they are cropped to what is
// in front of the camera
func (t Triangle) getSceneWireframe() []geometry.Line {
	minDepth := -0.01 // minimum z-coordinate to keep on screen
	lineAB := geometry.Line{A: t.A, B: t.B}.CropToFrontOfCamera(minDepth)
	lineAC := geometry.Line{A: t.A, B: t.C}.CropToFrontOfCamera(minDepth)
	lineBC := geometry.Line{A: t.B, B: t.C}.CropToFrontOfCamera(minDepth)
	if lineAB != nil && lineAC != nil && lineBC != nil {
		// all lines are in front of the camera
		return []geometry.Line{
			*lineAB, *lineAC, *lineBC,
		}
	}
	// all lines are behind camera
	if lineAB == nil && lineAC == nil && lineBC == nil {
		// all points are behind camera
		return []geometry.Line{}
	}
	// if one line is missing, replace with the endpoints of the other lines
	if !t.A.IsInFrontOfCamera(minDepth) && !t.B.IsInFrontOfCamera(minDepth) {
		// only C is in front of the screen
		return []geometry.Line{
			*lineAC,
			*lineBC,
			{
				A: lineAC.A, B: lineBC.A,
			},
		}
	} else if !t.A.IsInFrontOfCamera(minDepth) && !t.C.IsInFrontOfCamera(minDepth) {
		// only B is in front of the screen
		return []geometry.Line{
			*lineAB,
			*lineBC,
			{
				A: lineBC.B, B: lineAB.A,
			},
		}
	} else if !t.B.IsInFrontOfCamera(minDepth) && !t.C.IsInFrontOfCamera(minDepth) {
		// only A is in front of the screen
		return []geometry.Line{
			*lineAB,
			*lineAC,
			{
				A: lineAB.B, B: lineAC.B,
			},
		}
	}

	// if two lines are missing, not sure what that means
	panic(fmt.Errorf("two lines are missing, not sure what to do \n%s \n%s \n%s \n%v \n%v \n%v", t.A, t.B, t.C, lineAB, lineAC, lineBC))
}

func (t Triangle) GetWireframe() []geometry.RasterLine {
	sceneLines := t.getSceneWireframe()
	ret := []geometry.RasterLine{}
	for _, l := range sceneLines {
		rasterLine := l.CropToScreenView()
		if rasterLine != nil {
			ret = append(ret, *rasterLine)
		}
	}
	return ret
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
	return t.B.Subtract(t.A)
}

func (t Triangle) getCVect() geometry.Vector3D {
	return t.C.Subtract(t.A)
}

// return the intersection in triangle-local coordinates, in direction of A->B and A->C
// bool signifies whether intersection is inside the triange
// third float is the z-depth, in positive values
func (t *Triangle) RayIntersectLocalCoords(r ray) []intersection {
	// cache the vectors AB and AC, as well as the plane, this is 37% more efficient
	if !t.cached {
		t.bVect = t.getBVect()
		t.cVect = t.getCVect()
		t.plane = t.getPlane()
		t.normal = t.plane.N
		t.normalMagSq = t.normal.Mag() * t.normal.Mag()
		t.cached = true
	}
	intersectDot, doesIntersect := t.plane.IntersectPoint(r)
	if !doesIntersect {
		return nil
		// return 0, 0, 0, false
	}
	iVect := intersectDot.Subtract(t.A)
	// iMag := geometry.OriginPoint.Subtract(intersectDot).Mag()
	zDepth := -intersectDot.Z

	bVect := t.bVect
	cVect := t.cVect
	normal := t.normal
	normalMagSq := t.normalMagSq
	b := iVect.CrossProduct(cVect).DotProduct(normal) / normalMagSq
	c := bVect.CrossProduct(iVect).DotProduct(normal) / normalMagSq
	// check if vector (b,c) is inside the triangle [(0,0), (1,0), (0,1)]
	if b < 0 || b > 1 || c < 0 || c > 1 {
		// outside the unit square
		// return b, c, zDepth, false
		return nil
	}
	if b+c > 1 {
		// inside unit square, but on far side of hypotenuse
		return nil
		// return b, c, zDepth, false
	}
	// inside unit square and inside the hypotenuse
	return []intersection{{b, c, zDepth}}
	// return b, c, zDepth, true
}
