package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
)

type DynamicScene interface {
	GetFrame(t float64) StaticScene
}

type StaticScene interface {
	Flatten() ([]objects.StaticBasicObject, Background)
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
	CameraDirection geometry.Direction
}

func (s ObjectScene) Flatten() ([]objects.StaticBasicObject, Background) {
	tris := []objects.StaticBasicObject{}
	inverseMatrix := s.CameraDirection.InverseHomoMatrix()
	for _, obj := range s.Objects {
		basics := obj.Flatten()
		movedBasics := []objects.StaticBasicObject{}
		for _, b := range basics {
			movedBasics = append(movedBasics, b.ApplyMatrix(inverseMatrix))
		}
		tris = append(tris, movedBasics...)
	}

	return tris, s.Background
}

// implements DynamicScene
type CombinedDynamicScene struct {
	Objects    []objects.DynamicObjectInt
	CameraPath geometry.Path
	Background DynamicBackground
}

func (s CombinedDynamicScene) GetFrame(t float64) StaticScene {
	frameObjects := make([]objects.StaticObject, len(s.Objects))
	for i, object := range s.Objects {
		obj := object.Frame(t)
		frameObjects[i] = obj
	}

	direction := geometry.OriginPosition
	if s.CameraPath != nil {
		direction = s.CameraPath.GetDirection(t)
	}
	return ObjectScene{
		Objects:         frameObjects,
		Background:      s.Background.GetFrame(t),
		CameraDirection: direction,
	}
}
