package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/objects"
)

type DynamicScene interface {
	GetFrame(t float64) StaticScene
}

type StaticScene interface {
	Flatten() ([]*objects.StaticTriangle, Background)
}

type Background interface {
	GetColor(x, y float64) colors.Color
}

type DynamicBackground interface {
	GetFrame(t float64) Background
}

type ObjectScene struct {
	Objects []objects.StaticObject
	Background
}

func (s ObjectScene) Flatten() ([]*objects.StaticTriangle, Background) {
	tris := []*objects.StaticTriangle{}
	for _, obj := range s.Objects {
		tris = append(tris, obj.Flatten()...)
	}
	return tris, s.Background
}

// implements DynamicScene
type CombinedDynamicScene struct {
	Objects    []objects.DynamicObjectInt
	Background DynamicBackground
}

func (s CombinedDynamicScene) GetFrame(t float64) StaticScene {
	frameObjects := make([]objects.StaticObject, len(s.Objects))
	for i, object := range s.Objects {
		frameObjects[i] = object.Frame(t)
	}
	return ObjectScene{
		Objects:    frameObjects,
		Background: s.Background.GetFrame(t),
	}
}
