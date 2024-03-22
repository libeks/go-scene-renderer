package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/geometry"
)

type Object interface {
	// return all the triangles that are part of this object, to simplify computation
	Flatten() []*Triangle

	// String() string
}

type TransformableObject interface {
	Object
	ApplyMatrix(m geometry.HomogeneusMatrix) TransformableObject
}

type DynamicObject interface {
	GetFrame(t float64) Object
}

// TransformedObject implements DynamicObject
type TransformedObject struct {
	Object   TransformableObject
	MatrixFn func(t float64) geometry.HomogeneusMatrix
}

func (o TransformedObject) GetFrame(t float64) Object {
	m := o.MatrixFn(t)
	return o.Object.ApplyMatrix(m)
}

// ComplexObject implements TransformableObject
type ComplexObject struct {
	Objs []TransformableObject
}

func (o ComplexObject) ApplyMatrix(m geometry.HomogeneusMatrix) TransformableObject {
	newTriangles := make([]TransformableObject, len(o.Objs))
	for i, obj := range o.Objs {
		newTriangles[i] = obj.ApplyMatrix(m)
	}
	return ComplexObject{
		Objs: newTriangles,
	}
}

func (o ComplexObject) Flatten() []*Triangle {
	tris := []*Triangle{}
	for _, obj := range o.Objs {
		tris = append(tris, obj.Flatten()...)
	}
	return tris
}

func (o ComplexObject) String() string {
	return fmt.Sprintf("ComplexObject: []{%v}", o.Objs)
}

type BoundingBox struct {
	TopLeft     geometry.Pixel
	BottomRight geometry.Pixel
	MinDepth    float64
	MaxDepth    float64
}
