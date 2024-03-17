package objects

import (
	"fmt"
	"math"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
)

type Object interface {
	// returns the color of the object at a ray
	// emanating from the camera at (0,0,0), pointed in the direction
	// (x,y, -1), with perspective
	// and a z-index. The bigger the index, the farther the object.
	// If there is no intersection, return (nil, 0.0)
	GetColorDepth(x, y float64) (*color.Color, float64)

	String() string
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

func (o ComplexObject) GetColorDepth(x, y float64) (*color.Color, float64) {
	minZ := math.MaxFloat64
	var closestColor *color.Color
	for _, obj := range o.Objs {
		c, depth := obj.GetColorDepth(x, y)
		if c != nil && depth < minZ {
			minZ = depth
			closestColor = c
		}
	}
	if closestColor != nil {
		return closestColor, minZ
	}
	return nil, 0
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

func (o ComplexObject) String() string {
	return fmt.Sprintf("ComplexObject: []{%v}", o.Objs)
}
