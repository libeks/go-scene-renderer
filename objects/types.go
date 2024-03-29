package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/geometry"
)

var (
	identityTransform = func(t float64) geometry.HomogeneusMatrix {
		return geometry.HomogeneusIdentity
	}
)

type StaticObject struct {
	triangles []StaticTriangle
}

func (ob StaticObject) Flatten() []*StaticTriangle {
	triPointers := make([]*StaticTriangle, len(ob.triangles))
	for i, tri := range ob.triangles {
		triPointers[i] = &tri
	}
	return triPointers
}

func DynamicObjectFromTriangles(tris ...DynamicTriangle) DynamicObject {
	newTriangles := make([]triWithTransform, len(tris))
	for i, tri := range tris {
		newTriangles[i] = triWithTransform{

			tri: tri,
			fn:  identityTransform,
		}
	}
	return DynamicObject{
		newTriangles,
	}

}

func CombineDynamicObjects(objs ...DynamicObject) DynamicObject {
	newTriangles := []triWithTransform{}
	for _, obj := range objs {
		newTriangles = append(newTriangles, obj.triangles...)
	}
	return DynamicObject{
		newTriangles,
	}
}

type triWithTransform struct {
	tri DynamicTriangle
	fn  func(float64) geometry.HomogeneusMatrix
}

type DynamicObject struct {
	triangles []triWithTransform
}

func (ob DynamicObject) Frame(t float64) StaticObject {
	staticTriangles := make([]StaticTriangle, 0, len(ob.triangles))
	for _, dyTriangle := range ob.triangles {
		transformedTriangle := dyTriangle.tri.ApplyMatrix(dyTriangle.fn(t))
		if transformedTriangle != nil {
			staticTriangles = append(staticTriangles, transformedTriangle.Frame(t)) // set texture to be static
		}
	}
	return StaticObject{
		triangles: staticTriangles,
	}
}

func (ob DynamicObject) WithTransform(m geometry.HomogeneusMatrix) DynamicObject {
	newTriangles := make([]triWithTransform, 0, len(ob.triangles))
	for _, tri := range ob.triangles {
		newTriangles = append(newTriangles, triWithTransform{
			tri: tri.tri,
			fn: func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(m, tri.fn(t))
			},
		})
	}
	return DynamicObject{
		triangles: newTriangles,
	}
}

func (ob DynamicObject) WithDynamicTransform(f func(float64) geometry.HomogeneusMatrix) DynamicObject {
	newTriangles := make([]triWithTransform, 0, len(ob.triangles))
	for _, tri := range ob.triangles {
		newTriangles = append(newTriangles, triWithTransform{
			tri: tri.tri,
			fn: func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(f(t), tri.fn(t))
			},
		})
	}
	return DynamicObject{
		triangles: newTriangles,
	}
}

type BoundingBox struct {
	TopLeft     geometry.Pixel
	BottomRight geometry.Pixel
	MinDepth    float64
	MaxDepth    float64
	empty       bool
}

func (bb BoundingBox) IsEmpty() bool {
	return bb.empty
}

func (bb BoundingBox) String() string {
	return fmt.Sprintf("BB(%s %s)", bb.TopLeft, bb.BottomRight)
}

type ray struct {
	P geometry.Point    // origin point
	D geometry.Vector3D // direction vector describing the ray
}

type plane struct {
	N geometry.Vector3D // normal vector
	D float64           // d parameter, describing plane
}

func (p plane) String() string {
	return fmt.Sprintf("Plane(->%s, at %f)", p.N, p.D)
}

func (p plane) IntersectPoint(r ray) *geometry.Point {
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
