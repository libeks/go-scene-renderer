package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
)

// A scene that has several frames, indexed into by the t parameter, ranging from 0.0 -> 1.0
type AnimatedScene interface {
	// GetColor operates on a frame where x, y range from -1.0 -> 1.0, centered on (0.0, 0.0)
	GetFrameColor(x, y float64, t float64) color.Color
}

type Object interface {
	// returns the color of the object at a ray
	// emanating from the camera at (0,0,0), pointed in the direction
	// (x,y, -1), with perspective
	// if there is no intersection, return nil
	// and a z-index. The bigger the index, the farther the object
	GetColorDepth(x, y float64) (*color.Color, float64)
}

type TransformableObject interface {
	Object

	ApplyMatrix(m geometry.HomogeneusMatrix) TransformableObject
}

type DynamicObject interface {
	GetFrame(t float64) Object
}

type DynamicScene interface {
	GetFrame(t float64) Frame
}

type Frame interface {
	GetColor(x, y float64) color.Color
}
