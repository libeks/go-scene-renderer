package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/geometry"
)

var (
	identityTransform = func(t float64) geometry.HomogeneusMatrix {
		return geometry.HomogeneusIdentity
	}

	EmptyBB = BoundingBox{empty: true}
)

type DynamicObjectInt interface {
	Frame(float64) StaticObject
}

type BasicObject interface {
	// GetColorDepth(x, y float64) (*colors.Color, float64)
	ApplyMatrix(m geometry.HomogeneusMatrix) BasicObject
	GetBoundingBox() BoundingBox
	GetWireframe() []geometry.RasterLine
	RayIntersectLocalCoords(ray) []intersection
}

type StaticObject struct {
	basics []StaticBasicObject
}

func (ob StaticObject) Flatten() []StaticBasicObject {
	return ob.basics
}

func DynamicObjectFromBasics(basics ...dynamicBasicObject) DynamicObject {
	newObjs := make([]objWithTransform, len(basics))
	for i, tri := range basics {
		newObjs[i] = objWithTransform{
			obj: dynamicTriangleWrapper{tri},
			fn:  identityTransform,
		}
	}
	return DynamicObject{
		newObjs,
	}

}

func CombineDynamicObjects(objs ...DynamicObject) DynamicObject {
	newObjs := []objWithTransform{}
	for _, obj := range objs {
		newObjs = append(newObjs, obj.objs...)
	}
	return DynamicObject{
		newObjs,
	}
}

type dynamicTriangleWrapper struct {
	tri dynamicBasicObject
}

func (d dynamicTriangleWrapper) Frame(t float64) StaticObject {
	return StaticObject{[]StaticBasicObject{d.tri.Frame(t)}}
}

type objWithTransform struct {
	obj DynamicObjectInt
	fn  func(float64) geometry.HomogeneusMatrix
}

func NewDynamicObject(obj DynamicObjectInt) DynamicObject {
	return DynamicObject{objs: []objWithTransform{{obj: obj, fn: identityTransform}}}
}

type DynamicObject struct {
	objs []objWithTransform
}

func (ob DynamicObject) Frame(t float64) StaticObject {
	staticTriangles := []StaticBasicObject{}
	for _, dyObj := range ob.objs {
		staticTris := dyObj.obj.Frame(t)
		for _, tri := range staticTris.basics {
			transformedTriangle := tri.ApplyMatrix(dyObj.fn(t))
			staticTriangles = append(staticTriangles, transformedTriangle) // set texture to be static
		}
	}
	return StaticObject{
		basics: staticTriangles,
	}
}

func (ob DynamicObject) WithTransform(m geometry.HomogeneusMatrix) DynamicObject {
	newTriangles := make([]objWithTransform, 0, len(ob.objs))
	for _, tri := range ob.objs {
		newTriangles = append(newTriangles, objWithTransform{
			obj: tri.obj,
			fn: func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(m, tri.fn(t))
			},
		})
	}
	return DynamicObject{
		objs: newTriangles,
	}
}

func (ob DynamicObject) WithDynamicTransform(f func(float64) geometry.HomogeneusMatrix) DynamicObject {
	newTriangles := make([]objWithTransform, 0, len(ob.objs))
	for _, tri := range ob.objs {
		newTriangles = append(newTriangles, objWithTransform{
			obj: tri.obj,
			fn: func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(f(t), tri.fn(t))
			},
		})
	}
	return DynamicObject{
		objs: newTriangles,
	}
}

type BoundingBox struct {
	TopLeft     geometry.Pixel
	BottomRight geometry.Pixel
	MinZDepth   float64
	MaxZDepth   float64
	empty       bool
}

func (bb BoundingBox) IsEmpty() bool {
	return bb.empty
}

func (bb BoundingBox) String() string {
	return fmt.Sprintf("BB(%s %s zmin:%.3f zmax:%.3f)", bb.TopLeft, bb.BottomRight, bb.MinZDepth, bb.MaxZDepth)
}

type ray struct {
	P geometry.Point    // origin point
	D geometry.Vector3D // direction vector describing the ray
}

func (r ray) PointAt(t float64) geometry.Point {
	return geometry.Point(r.P.Vector().AddVector(r.D.ScalarMultiply(t)))
}

type plane struct {
	N geometry.Vector3D // normal vector
	D float64           // d parameter, describing plane
}

func (p plane) String() string {
	return fmt.Sprintf("Plane(->%s, at %f)", p.N, p.D)
}

func (p plane) IntersectPoint(r ray) (geometry.Point, bool) {
	denominator := p.N.DotProduct(r.D)
	if denominator == 0.0 {
		return geometry.Point{}, false // ray is parallel to plane, no intersection
	}
	t := (p.D - p.N.DotProduct(r.P.Vector())) / denominator
	if t < 0.0 {
		return geometry.Point{}, false // ray intersects plane before ray's starting point
	}
	point := r.PointAt(t)
	return point, true
}

type intersection struct {
	b      float64
	c      float64
	zDepth float64
}
