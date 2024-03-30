package objects

import (
	"fmt"
	"slices"

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
func (t DynamicTriangle) GetWireframe() []geometry.RasterLine {
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
	b, c, depth, intersect := t.rayIntersectLocalCoords(ray{geometry.OriginPoint, geometry.Vector3D{x, y, -1}})
	if !intersect {
		return nil, 0
	}
	color := t.Colorer.GetTextureColor(b, c)
	return &color, depth
}

func (t StaticTriangle) ApplyMatrix(m geometry.HomogeneusMatrix) *StaticTriangle {
	newTriangle := t.Triangle.ApplyMatrix(m)
	if newTriangle == nil {
		return nil
	}
	return &StaticTriangle{
		Triangle: *newTriangle,
		Colorer:  t.Colorer,
	}
}

func (t StaticTriangle) GetBoundingBox() BoundingBox {
	return t.Triangle.GetBoundingBox()
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
func (t StaticTriangle) GetWireframe() []geometry.RasterLine {
	return t.Triangle.GetWireframe()
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
	if t.cachedBoundingBox {
		return t.bbox
	}
	wireframe := t.getSceneWireframe()
	points := []geometry.Pixel{}
	pointsX := []float64{}
	pointsY := []float64{}
	depths := []float64{}
	for _, line := range wireframe {
		pointA, depthA := line.A.ToPixel()
		pointB, depthB := line.B.ToPixel()
		if pointA == nil || pointB == nil {
			panic(fmt.Errorf("Line should already be in front of camera %s", line))
		}
		points = append(points, *pointA)
		pointsX = append(pointsX, pointA.X)
		pointsY = append(pointsY, pointA.Y)
		points = append(points, *pointB)
		pointsX = append(pointsX, pointB.X)
		pointsY = append(pointsY, pointB.Y)
		depths = append(depths, depthA, depthB)
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
	}
	t.bbox = bb
	t.cachedBoundingBox = true
	// fmt.Printf("bounding box %s\n", bb)
	return bb
}

// return all the lines that describe the triangle, without any fill, used to generate wireframe images
// note that the resulting lines may not exactly match the triangle, as they are cropped to what is
// in front of the camera
func (t Triangle) getSceneWireframe() []geometry.Line {
	minDepth := -0.01 // minimum z-coordinate to keep on screen
	lineAB := geometry.Line{t.A, t.B}.CropToFrontOfCamera(minDepth)
	lineAC := geometry.Line{t.A, t.C}.CropToFrontOfCamera(minDepth)
	lineBC := geometry.Line{t.B, t.C}.CropToFrontOfCamera(minDepth)
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
			geometry.Line{
				lineAC.A, lineBC.A,
			},
		}
	} else if !t.A.IsInFrontOfCamera(minDepth) && !t.C.IsInFrontOfCamera(minDepth) {
		// only B is in front of the screen
		return []geometry.Line{
			*lineAB,
			*lineBC,
			geometry.Line{
				lineBC.B, lineAB.A,
			},
		}
	} else if !t.B.IsInFrontOfCamera(minDepth) && !t.C.IsInFrontOfCamera(minDepth) {
		// only A is in front of the screen
		return []geometry.Line{
			*lineAB,
			*lineAC,
			geometry.Line{
				lineAB.B, lineAC.B,
			},
		}
	}

	// if two lines are missing, not sure what that means
	panic(fmt.Errorf("two lines are missing, not sure what to do \n%s \n%s \n%s \n%v \n%v \n%v", t.A, t.B, t.C, lineAB, lineAC, lineBC))
	return []geometry.Line{}
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
// third float is the depth, in positive values
func (t *Triangle) rayIntersectLocalCoords(r ray) (float64, float64, float64, bool) {
	// cache the vectors AB and AC, as well as the plane, this is 37% more efficient
	if !t.cached {
		t.bVect = t.getBVect()
		t.cVect = t.getCVect()
		t.plane = t.getPlane()
		t.normal = t.plane.N
		t.normalMagSq = t.normal.Mag() * t.normal.Mag()
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
	normalMagSq := t.normalMagSq
	b := iVect.CrossProduct(cVect).DotProduct(normal) / normalMagSq
	c := bVect.CrossProduct(iVect).DotProduct(normal) / normalMagSq
	// check if vector (b,c) is inside the triangle [(0,0), (1,0), (0,1)]
	if b < 0 || b > 1 || c < 0 || c > 1 {
		// outside the unit square
		return b, c, iMag, false
	}
	if b+c > 1 {
		// inside unit square, but on far side of hypotenuse
		return b, c, iMag, false
	}
	// inside unit square and inside the hypotenuse
	return b, c, iMag, true
}
