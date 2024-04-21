package scenes

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
)

type DynamicScene interface {
	GetFrame(t float64) StaticScene
}

type StaticScene interface {
	Flatten() ([]objects.BasicObject, Background)
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

func (s ObjectScene) Flatten() ([]objects.BasicObject, Background) {
	tris := []objects.BasicObject{}
	inverseMatrix := s.CameraDirection.InverseHomoMatrix()
	fmt.Printf("Inverse camera matrix: %s\n", inverseMatrix)
	for _, obj := range s.Objects {
		// obj = obj.ApplyMatrix(inverseMatrix)
		basics := obj.Flatten()
		movedBasics := []objects.BasicObject{}
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
	// direction = geometry.Direction{
	// 	Origin: geometry.Point{0, 0, 4},
	// 	Orientation: geometry.EulerDirection{
	// 		geometry.V3(0, 0, -1),
	// 		geometry.V3(0, 1, 0),
	// 		geometry.V3(1, 0, 0)},
	// }
	fmt.Printf("%.3f direction: %v\n", t, direction)
	return ObjectScene{
		Objects:         frameObjects,
		Background:      s.Background.GetFrame(t),
		CameraDirection: direction,
	}
}
