package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/objects"
)

// A scene that has several frames, indexed into by the t parameter, ranging from 0.0 -> 1.0
type AnimatedScene interface {
	// GetColor operates on a frame where x, y range from -1.0 -> 1.0, centered on (0.0, 0.0)
	GetFrameColor(x, y float64, t float64) color.Color
}

type DynamicScene interface {
	GetFrame(t float64) Frame
}

type Frame interface {
	GetColor(x, y float64) color.Color
}

type CombinedScene struct {
	Objects    []objects.Object
	Background Frame
}

func (s CombinedScene) GetColor(x, y float64) color.Color {
	minZ := math.MaxFloat64
	var closestColor *color.Color
	for _, obj := range s.Objects {
		c, depth := obj.GetColorDepth(x, y)
		if c != nil && depth < minZ {
			minZ = depth
			closestColor = c
		}
	}
	if closestColor != nil {
		return *closestColor
	}
	return s.Background.GetColor(x, y)
}

type CombinedDynamicScene struct {
	Objects    []objects.DynamicObject
	Background DynamicScene
}

func (s CombinedDynamicScene) GetFrame(t float64) Frame {
	frameObjects := make([]objects.Object, len(s.Objects))
	for i, object := range s.Objects {
		frameObjects[i] = object.GetFrame(t)
	}
	return CombinedScene{
		Objects:    frameObjects,
		Background: s.Background.GetFrame(t),
	}
}
