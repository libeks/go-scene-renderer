package objects

import (
	"fmt"
	"math"
	"slices"

	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/textures"
)

func DynamicSphere(t Sphere, colorer textures.DynamicTransparentTexture) dynamicBasicObject {
	return dynamicBasicObject{
		BasicObject: &t,
		Colorer:     colorer,
	}
}

func UnitSphere() Sphere {
	return Sphere{
		Center:  geometry.Point{X: 0, Y: 0, Z: 0},
		Radius:  1,
		Forward: geometry.V3(0, 0, 1),
		Up:      geometry.V3(0, 1, 0),
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
	return &ns
}

// return the bounding box of the cube surrounding the sphere, which is an overestimate, but good enough
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
	if len(xs) == 0 {
		s.cached = true
		s.bb = EmptyBB
		return s.bb
	}
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
	s.cached = true
	s.bb = bb
	return bb
}

// Return the wireframe of the cube surrounding the sphere
func (s Sphere) GetWireframe() []geometry.RasterLine {
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
		return nil
	}
	intersections := []intersection{}
	for _, root := range tRoots {
		intersectDot := r.PointAt(root)
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
	return intersections
}
