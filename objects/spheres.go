package objects

import (
	"fmt"
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
)

func DynamicSphere(t Sphere, colorer colors.DynamicTransparentTexture) dynamicBasicObject {
	return dynamicBasicObject{
		BasicObject: &t,
		Colorer:     colorer,
	}
}

func UnitSphere() Sphere {
	return Sphere{
		Center:  geometry.Point{0, 0, 0},
		Radius:  1,
		Forward: geometry.Vector3D{0, 0, 1},
		Up:      geometry.Vector3D{0, 1, 0},
	}
}

// // DynamicTriangle is a Triangle with a DynamicTexture, which can be evaluated for a specific frame
// type dynamicSphere struct {
// 	Sphere
// 	Colorer colors.DynamicTransparentTexture
// }

// func (t dynamicSphere) Frame(f float64) staticTriangle {
// 	return staticTriangle{
// 		Triangle: &t.Sphere,
// 		Colorer:  t.Colorer.GetFrame(f),
// 	}
// }

// func (t dynamicSphere) ApplyMatrix(m geometry.HomogeneusMatrix) dynamicSphere {
// 	newSphere := t.Sphere.ApplyMatrix(m)
// 	return dynamicSphere{
// 		Sphere:  newSphere,
// 		Colorer: t.Colorer,
// 	}
// }

// func (t dynamicSphere) GetBoundingBox() BoundingBox {
// 	return t.Sphere.GetBoundingBox()
// }

// func (t dynamicSphere) String() string {
// 	return fmt.Sprintf("DynamicSphere: %s with %s", t.Sphere, t.Colorer)
// }

// // return all the lines that describe the triangle, without any fill, used to generate wireframe images
// func (t dynamicSphere) GetWireframe() []geometry.RasterLine {
// 	return t.Sphere.GetWireframe()
// }

// implements BasicObject
type Sphere struct {
	Center  geometry.Point
	Radius  float64
	Forward geometry.Vector3D
	Up      geometry.Vector3D

	cached bool
	bb     BoundingBox
}

// func (s Sphere) GetColorDepth(x, y float64) (*colors.Color, float64) {
// 	b, c, depth, intersect := s.rayIntersectLocalCoords(ray{geometry.OriginPoint, geometry.V3(x, y, -1)})
// 	if !intersect {
// 		return nil, 0
// 	}
// 	colorPtr := s.Colorer.GetTextureColor(b, c)
// 	if colorPtr == nil {
// 		return nil, 0
// 	}
// 	return colorPtr, depth
// }

func (s Sphere) ApplyMatrix(m geometry.HomogeneusMatrix) BasicObject {
	fmt.Printf("Applying matrix %s to sphere %s\n", m, s)
	center, ok := m.MultVect(s.Center.ToHomogenous()).ToPoint()
	if !ok {
		panic(fmt.Errorf("could not apply matrix %s to point %s", m, s.Center))
	}
	m3D := m.Slice3DMatrix()
	forward := m3D.MultVect(s.Forward.Unit())
	up := m3D.MultVect(s.Up.Unit())
	radius := m3D.MultVect(s.Forward.ScalarMultiply(s.Radius)).Mag()

	ns := Sphere{
		Center:  center,
		Radius:  radius,
		Forward: forward,
		Up:      up,
		// Colorer: s.Colorer,
	}
	fmt.Printf("result sphere %s\n", ns)
	return &ns
}

func (s *Sphere) GetBoundingBox() BoundingBox {
	if s.cached {
		return s.bb
	}
	w := geometry.OriginPoint.Subtract(s.Center)
	winvUnit := w.ScalarMultiply(-1).Unit()
	// fmt.Printf("Sphere %s\n", s)
	minZDepth := -s.Center.Z - s.Radius
	maxZDepth := -s.Center.Z + s.Radius
	if maxZDepth < 0 {
		return EmptyBB
	}
	alpha := math.Asin(s.Radius / w.Mag())
	t := math.Cos(alpha) * w.Mag()
	// fmt.Printf("alpha: %.3f, t: %.3f\n", alpha, t)
	// v := geometry.RotateRoll(-alpha).MultVect(w.Unit())
	leftPt := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(alpha).MultVect(winvUnit)}.PointAt(t)
	rightPt := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(-alpha).MultVect(winvUnit)}.PointAt(t)
	upPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(alpha).MultVect(winvUnit)}.PointAt(t)
	downPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(-alpha).MultVect(winvUnit)}.PointAt(t)
	// fmt.Printf("lr up points: %s %s %s %s\n", leftPt, rightPt, upPt, downPt)
	left, _ := leftPt.ToPixel()
	right, _ := rightPt.ToPixel()
	up, _ := upPt.ToPixel()
	down, _ := downPt.ToPixel()
	// fmt.Printf("lr ul: %s %s %s %s, %v %v %v %v\n", left, right, up, down, ok1, ok2, ok3, ok4)
	bb := BoundingBox{
		TopLeft: geometry.Pixel{
			X: max(left.X, -1.0),
			Y: max(up.Y, -1.0),
		},
		BottomRight: geometry.Pixel{
			X: min(right.X, 1.0),
			Y: min(down.Y, 1.0),
		},
		MinZDepth: max(0, minZDepth),
		MaxZDepth: maxZDepth,
	}
	s.cached = true
	s.bb = bb
	// fmt.Printf("Sphere bbox: %s\n", bb)
	return bb
}

func (s Sphere) GetWireframe() []geometry.RasterLine {
	// what if center is behind camera?

	if s.Center.Z > 0 {
		// center behind fulcrum, no wireframe
		return []geometry.RasterLine{}
	}
	w := geometry.OriginPoint.Subtract(s.Center)

	alpha := math.Asin(s.Radius / w.Mag())
	t := math.Cos(alpha) * w.Mag()
	// v := geometry.RotateRoll(-alpha).MultVect(w.Unit())
	leftPt := ray{P: geometry.OriginPoint, D: geometry.RotateRoll3D(-alpha).MultVect(w.Unit())}.PointAt(t)
	rightPt := ray{P: geometry.OriginPoint, D: geometry.RotateRoll3D(alpha).MultVect(w.Unit())}.PointAt(t)
	upPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(-alpha).MultVect(w.Unit())}.PointAt(t)
	downPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(alpha).MultVect(w.Unit())}.PointAt(t)
	topLine := geometry.Line{A: geometry.Point{leftPt.X, upPt.Y, leftPt.Z}, B: geometry.Point{rightPt.X, upPt.Y, leftPt.Z}}
	bottomLine := geometry.Line{A: geometry.Point{leftPt.X, downPt.Y, leftPt.Z}, B: geometry.Point{rightPt.X, downPt.Y, leftPt.Z}}
	leftLine := geometry.Line{A: geometry.Point{leftPt.X, upPt.Y, leftPt.Z}, B: geometry.Point{leftPt.X, downPt.Y, leftPt.Z}}
	rightLine := geometry.Line{A: geometry.Point{rightPt.X, upPt.Y, leftPt.Z}, B: geometry.Point{rightPt.X, downPt.Y, leftPt.Z}}
	sceneLines := []geometry.Line{topLine, bottomLine, leftLine, rightLine}

	ret := []geometry.RasterLine{}
	for _, l := range sceneLines {
		rasterLine := l.CropToScreenView()
		if rasterLine != nil {
			ret = append(ret, *rasterLine)
		}
	}
	return ret
}

func (s Sphere) String() string {
	return fmt.Sprintf("Sphere at %s with radius %.3f, up: %s: forward: %s", s.Center, s.Radius, s.Up, s.Forward)
}

// return the intersection in triangle-local coordinates, in direction of A->B and A->C
// bool signifies whether intersection is inside the triange
// third float is the depth, in positive values
func (s Sphere) RayIntersectLocalCoords(r ray) []intersection {
	// let w be the vector from sphere center to ray origin, it'll make the math simpler
	w := r.P.Subtract(s.Center)
	v := r.D
	// solve the quadratic equation
	// t^2* v * v + t * 2 * v * w + w * w - r ^ 2 = 0
	tRoots := maths.QuadraticRoots(v.DotProduct(v), 2*v.DotProduct(w), w.DotProduct(w)-s.Radius*s.Radius)
	tRoots = maths.KeepOnlyPositives(tRoots)
	if len(tRoots) == 0 {
		// fmt.Printf("No intersection at %s\n", r)
		return nil
	}
	intersections := []intersection{}
	// fmt.Printf("Roots at %s: %v\n", r, tRoots)
	for _, root := range tRoots {
		intersectDot := r.PointAt(root)
		// iMag := geometry.OriginPoint.Subtract(intersectDot).Mag()
		iVect := intersectDot.Subtract(s.Center)
		// vertical component c
		cComp := iVect.DotProduct(s.Up) / s.Radius
		c := math.Acos(cComp)
		iVect = iVect.AddVector(s.Up.ScalarMultiply(-s.Radius * cComp))
		// iVect should now be in the horizontal plane of the circle, should have a zero y-coord
		b := math.Atan2(iVect.CrossProduct(s.Forward).DotProduct(s.Up), iVect.DotProduct(s.Forward))
		if b < 0 {
			b += 2 * math.Pi
		}
		c = c / (math.Pi)     // so it's from 0 to 1
		b = b / (2 * math.Pi) // so it's from 0 to 1
		intersections = append(intersections, intersection{b, c, -intersectDot.Z})
	}
	// fmt.Printf("intersections: %v\n", intersections)
	return intersections
}
