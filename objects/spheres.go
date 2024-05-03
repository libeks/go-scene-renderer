package objects

import (
	"fmt"
	"math"
	"slices"

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

// implements BasicObject
type Sphere struct {
	Center  geometry.Point
	Radius  float64
	Forward geometry.Vector3D
	Up      geometry.Vector3D

	cached bool
	bb     BoundingBox
}

func (s Sphere) ApplyMatrix(m geometry.HomogeneusMatrix) BasicObject {
	fmt.Printf("Applying matrix %s to sphere %s\n", m, s)
	center, ok := m.MultVect(s.Center.ToHomogenous()).ToPoint()
	if !ok {
		panic(fmt.Errorf("could not apply matrix %s to point %s", m, s.Center))
	}
	m3D := m.Slice3DMatrix()
	forward := m3D.MultVect(s.Forward.Unit()).Unit()
	up := m3D.MultVect(s.Up.Unit()).Unit()
	radius := m3D.MultVect(s.Forward.ScalarMultiply(s.Radius)).Mag()

	ns := Sphere{
		Center:  center,
		Radius:  radius,
		Forward: forward,
		Up:      up,
	}
	fmt.Printf("result sphere %s\n", ns)
	return &ns
}

func (s *Sphere) GetBoundingBox() BoundingBox {
	if s.cached {
		return s.bb
	}
	xs, ys := []float64{}, []float64{}
	rasterLines := s.GetWireframe()
	for _, line := range rasterLines {
		xs = append(xs, line.A.X, line.B.X)
		ys = append(ys, line.A.Y, line.B.Y)
	}
	// fmt.Printf("ys: %v, min %.3f, max %.3f\n", ys, slices.Min(ys), slices.Max(ys))
	bb := BoundingBox{
		TopLeft: geometry.Pixel{
			X: max(slices.Min(xs), -1.0),
			Y: max(slices.Min(ys), -1.0),
		},
		BottomRight: geometry.Pixel{
			X: min(slices.Max(xs), 1.0),
			Y: min(slices.Max(ys), 1.0),
		},
		MinZDepth: max(0, -(s.Center.Z + s.Radius)),
		MaxZDepth: -(s.Center.Z - s.Radius),
	}
	fmt.Printf("bb %s\n", bb)
	s.cached = true
	s.bb = bb
	return bb
}

func sceneLinesAroundPoint(pt geometry.Point) []geometry.Line {
	uD := geometry.V3(0, .2, 0)
	rD := geometry.V3(.2, 0, 0)
	bD := geometry.V3(0, 0, .2)

	return []geometry.Line{
		{A: pt, B: geometry.Point(pt.Vector().AddVector(uD))},
		{A: pt, B: geometry.Point(pt.Vector().AddVector(rD))},
		geometry.Line{A: pt, B: geometry.Point(pt.Vector().AddVector(bD))},
	}
}

func (s Sphere) GetWireframe() []geometry.RasterLine {
	if s.Center.Z > 0 {
		// center behind fulcrum, no wireframe
		return []geometry.RasterLine{}
	}
	up := geometry.V3(0, 1, 0).ScalarMultiply(s.Radius)
	left := geometry.V3(-1, 0, 0).ScalarMultiply(s.Radius)
	away := geometry.V3(0, 0, -1).ScalarMultiply(s.Radius)
	c000 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(1).AddVector(up.ScalarMultiply(-1)).AddVector(away.ScalarMultiply(-1))))
	c100 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(-1).AddVector(up.ScalarMultiply(-1)).AddVector(away.ScalarMultiply(-1))))
	c010 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(1).AddVector(up.ScalarMultiply(1)).AddVector(away.ScalarMultiply(-1))))
	c110 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(-1).AddVector(up.ScalarMultiply(1)).AddVector(away.ScalarMultiply(-1))))

	c001 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(1).AddVector(up.ScalarMultiply(-1)).AddVector(away.ScalarMultiply(1))))
	c101 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(-1).AddVector(up.ScalarMultiply(-1)).AddVector(away.ScalarMultiply(1))))
	c011 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(1).AddVector(up.ScalarMultiply(1)).AddVector(away.ScalarMultiply(1))))
	c111 := geometry.Point(s.Center.Vector().AddVector(left.ScalarMultiply(-1).AddVector(up.ScalarMultiply(1)).AddVector(away.ScalarMultiply(1))))

	sceneLines := []geometry.Line{
		{A: c000, B: c100},
		{A: c110, B: c100},
		{A: c000, B: c010},
		{A: c110, B: c010},

		{A: c001, B: c101},
		{A: c111, B: c101},
		{A: c001, B: c011},
		{A: c111, B: c011},

		{A: c000, B: c001},
		{A: c100, B: c101},
		{A: c010, B: c011},
		{A: c110, B: c111},
	}

	ret := []geometry.RasterLine{}
	for _, l := range sceneLines {
		rasterLine := l.CropToScreenView()
		if rasterLine != nil {
			ret = append(ret, *rasterLine)
		}
	}
	fmt.Printf("Wireframe %s\n", ret)
	return ret
}

func (s Sphere) GetWireframe3() []geometry.RasterLine {
	if s.Center.Z > 0 {
		// center behind fulcrum, no wireframe
		return []geometry.RasterLine{}
	}
	w := geometry.OriginPoint.Subtract(s.Center)
	unitW := w.Unit().ScalarMultiply(-1)   // unit vector from center of sphere to the camera
	alpha := math.Acos(s.Radius / w.Mag()) // angle between center of sphere to the camera and the normal circle
	fmt.Printf("acos(%.3f/%.3f) is %.3f\n", s.Radius, w.Mag(), alpha)

	amountToCamera := (s.Radius * s.Radius) / (w.Mag())
	fmt.Printf("sin(%.3f) is %.3f\n", alpha, math.Sin(alpha))
	nRadius := s.Radius * math.Sin(alpha)
	nCenter := geometry.Point(s.Center.Add(geometry.Point(unitW.ScalarMultiply(-amountToCamera))))
	rightVect := unitW.CrossProduct(geometry.V3(0, 1, 0)).Unit()
	upVect := w.CrossProduct(rightVect).Unit()
	fmt.Printf("up: %s, right: %s\n", upVect, rightVect)

	left := geometry.Point(nCenter.Vector().AddVector(rightVect.ScalarMultiply(-nRadius)))
	right := geometry.Point(nCenter.Vector().AddVector(rightVect.ScalarMultiply(nRadius)))
	up := geometry.Point(nCenter.Vector().AddVector(upVect.ScalarMultiply(nRadius)))
	down := geometry.Point(nCenter.Vector().AddVector(upVect.ScalarMultiply(-nRadius)))

	// distances to center
	fmt.Printf("Distances to center are %.3f, %.3f, %.3f, %.3f\n", s.Center.Subtract(left).Mag(), s.Center.Subtract(right).Mag(), s.Center.Subtract(up).Mag(), s.Center.Subtract(down).Mag())

	sceneLines := []geometry.Line{}
	sceneLines = append(sceneLines,
		sceneLinesAroundPoint(left)...,
	)
	sceneLines = append(sceneLines,
		sceneLinesAroundPoint(right)...,
	)
	sceneLines = append(sceneLines,
		sceneLinesAroundPoint(up)...,
	)
	sceneLines = append(sceneLines,
		sceneLinesAroundPoint(down)...,
	)
	ret := []geometry.RasterLine{}
	for _, l := range sceneLines {
		rasterLine := l.CropToScreenView()
		if rasterLine != nil {
			ret = append(ret, *rasterLine)
		}
	}
	fmt.Printf("Wireframe %s\n", ret)
	return ret
}

func (s Sphere) GetWireframe2() []geometry.RasterLine {
	// what if center is behind camera?

	if s.Center.Z > 0 {
		// center behind fulcrum, no wireframe
		return []geometry.RasterLine{}
	}
	w := geometry.OriginPoint.Subtract(s.Center)
	unitW := w.ScalarMultiply(-1)

	alpha := math.Asin(s.Radius / -s.Center.Z)
	fmt.Printf("Alpha %.3f, wunit: %v\n", alpha, w.Unit())
	tt := math.Cos(alpha) * math.Cos(alpha) * w.Mag()
	t := math.Cos(alpha) * w.Mag()
	uleft := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(alpha).MatrixMult(geometry.RotatePitch3D(-alpha)).MultVect(unitW)}.PointAt(tt)
	uright := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(-alpha).MatrixMult(geometry.RotatePitch3D(-alpha)).MultVect(unitW)}.PointAt(tt)
	dleft := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(alpha).MatrixMult(geometry.RotatePitch3D(alpha)).MultVect(unitW)}.PointAt(tt)
	dright := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(-alpha).MatrixMult(geometry.RotatePitch3D(alpha)).MultVect(unitW)}.PointAt(tt)

	uD := geometry.V3(0, .2, 0)
	rD := geometry.V3(.2, 0, 0)
	bD := geometry.V3(0, 0, .2)
	sceneLines := []geometry.Line{
		{A: uleft, B: geometry.Point(uleft.Vector().AddVector(uD))},
		{A: uleft, B: geometry.Point(uleft.Vector().AddVector(rD))},
		{A: uleft, B: geometry.Point(uleft.Vector().AddVector(bD))},

		{A: uright, B: geometry.Point(uright.Vector().AddVector(uD))},
		{A: uright, B: geometry.Point(uright.Vector().AddVector(rD))},
		{A: uright, B: geometry.Point(uright.Vector().AddVector(bD))},

		{A: dleft, B: geometry.Point(dleft.Vector().AddVector(uD))},
		{A: dleft, B: geometry.Point(dleft.Vector().AddVector(rD))},
		{A: dleft, B: geometry.Point(dleft.Vector().AddVector(bD))},

		{A: dright, B: geometry.Point(dright.Vector().AddVector(uD))},
		{A: dright, B: geometry.Point(dright.Vector().AddVector(rD))},
		{A: dright, B: geometry.Point(dright.Vector().AddVector(bD))},
	}

	leftPt := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(alpha).MultVect(unitW)}.PointAt(t)
	rightPt := ray{P: geometry.OriginPoint, D: geometry.RotateYaw3D(-alpha).MultVect(unitW)}.PointAt(t)
	upPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(-alpha).MultVect(unitW)}.PointAt(t)
	downPt := ray{P: geometry.OriginPoint, D: geometry.RotatePitch3D(alpha).MultVect(unitW)}.PointAt(t)
	sceneLines = append(sceneLines,
		geometry.Line{A: leftPt, B: geometry.Point(leftPt.Vector().AddVector(uD))},
		geometry.Line{A: leftPt, B: geometry.Point(leftPt.Vector().AddVector(rD))},
		geometry.Line{A: leftPt, B: geometry.Point(leftPt.Vector().AddVector(bD))},

		geometry.Line{A: rightPt, B: geometry.Point(rightPt.Vector().AddVector(uD))},
		geometry.Line{A: rightPt, B: geometry.Point(rightPt.Vector().AddVector(rD))},
		geometry.Line{A: rightPt, B: geometry.Point(rightPt.Vector().AddVector(bD))},

		geometry.Line{A: upPt, B: geometry.Point(upPt.Vector().AddVector(uD))},
		geometry.Line{A: upPt, B: geometry.Point(upPt.Vector().AddVector(rD))},
		geometry.Line{A: upPt, B: geometry.Point(upPt.Vector().AddVector(bD))},

		geometry.Line{A: downPt, B: geometry.Point(downPt.Vector().AddVector(uD))},
		geometry.Line{A: downPt, B: geometry.Point(downPt.Vector().AddVector(rD))},
		geometry.Line{A: downPt, B: geometry.Point(downPt.Vector().AddVector(bD))},
	)
	// // fmt.Printf("Wireframe left: %s, right %s, up %s, down %s\n", leftPt, rightPt, upPt, downPt)
	// // left, _ := leftPt.ToPixel()
	// // right, _ := rightPt.ToPixel()
	// // up, _ := upPt.ToPixel()
	// // down, _ := downPt.ToPixel()
	topLine := geometry.Line{A: uleft, B: uright}
	bottomLine := geometry.Line{A: dleft, B: dright}
	leftLine := geometry.Line{A: uleft, B: dleft}
	rightLine := geometry.Line{A: uright, B: uleft}
	// sceneLines := []geometry.Line{topLine, bottomLine, leftLine, rightLine}
	sceneLines = append(sceneLines, topLine, bottomLine, leftLine, rightLine)
	fmt.Printf("SceneLines: %v\n", sceneLines)

	ret := []geometry.RasterLine{}
	for _, l := range sceneLines {
		rasterLine := l.CropToScreenView()
		if rasterLine != nil {
			ret = append(ret, *rasterLine)
		}
	}
	fmt.Printf("Wireframe %s\n", ret)
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
		if iVect.Mag()-s.Radius > 0.000001 {
			fmt.Printf("ivect %s mag %v, radius %v\n", iVect, iVect.Mag(), s.Radius)
			panic("Sphere vector is larger than radius")
		}
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
