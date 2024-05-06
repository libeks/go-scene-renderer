package objects

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/textures"
)

func DynamicBasicObject(b BasicObject, colorer textures.DynamicTransparentTexture) dynamicBasicObject {
	return dynamicBasicObject{
		BasicObject: b,
		Colorer:     colorer,
	}
}

// dynamicBasicObject is a BasicObject with a DynamicTexture, which can be evaluated for a specific frame
type dynamicBasicObject struct {
	BasicObject
	Colorer textures.DynamicTransparentTexture
}

func (t dynamicBasicObject) Frame(f float64) StaticBasicObject {
	return StaticBasicObject{
		BasicObject: t.BasicObject,
		Colorer:     t.Colorer.GetFrame(f),
	}
}

func (t dynamicBasicObject) ApplyMatrix(m geometry.HomogeneusMatrix) *dynamicBasicObject {
	newBasicObject := t.BasicObject.ApplyMatrix(m)
	return &dynamicBasicObject{
		BasicObject: newBasicObject,
		Colorer:     t.Colorer,
	}
}

// func (t dynamicBasicObject) GetBoundingBox() BoundingBox {
// 	return t.BasicObject.GetBoundingBox()
// }

func (t dynamicBasicObject) String() string {
	return fmt.Sprintf("dynamicBasicObject: %s with %s", t.BasicObject, t.Colorer)
}

// // return all the lines that describe the BasicObject, without any fill, used to generate wireframe images
// func (t dynamicBasicObject) GetWireframe() []geometry.RasterLine {
// 	return t.BasicObject.GetWireframe()
// }

func NewStaticBasicObject(t BasicObject, colorer textures.TransparentTexture) StaticBasicObject {
	return StaticBasicObject{
		BasicObject: t,
		Colorer:     colorer,
	}
}

// StaticBasicObject is a BasicObject with a Texture applied to it
type StaticBasicObject struct {
	BasicObject
	// Colorer will be evaluated with two parameters (b,c), each from (0,1), but b+c<1.0
	// it describes the coordinates on the BasicObject from A towards B and C, respectively
	Colorer textures.TransparentTexture
}

// returns the color of the BasicObject at a ray
// emanating from the camera at (0,0,0), pointed in the direction
// (x,y, -1), with perspective
// and a z-index. The bigger the index, the farther the object.
func (t StaticBasicObject) GetColorDepth(x, y float64) (*colors.Color, float64) {
	intersections := t.RayIntersectLocalCoords(ray{geometry.OriginPoint, geometry.V3(x, y, -1)})
	if len(intersections) == 0 {
		return nil, 0
	}
	for _, int := range intersections {
		colorPtr := t.Colorer.GetTextureColor(int.b, int.c)
		if colorPtr != nil {
			return colorPtr, int.zDepth
		}
	}
	return nil, 0
}

func (t StaticBasicObject) ApplyMatrix(m geometry.HomogeneusMatrix) StaticBasicObject {
	newBasicObject := t.BasicObject.ApplyMatrix(m)
	return StaticBasicObject{
		BasicObject: newBasicObject,
		Colorer:     t.Colorer,
	}
}

// func (t staticBasicObject) GetBoundingBox() BoundingBox {
// 	return t.BasicObject.GetBoundingBox()
// }

// // return all the lines that describe the BasicObject, without any fill, used to generate wireframe images
// func (t staticBasicObject) GetWireframe() []geometry.RasterLine {
// 	return t.BasicObject.GetWireframe()
// }

func (t StaticBasicObject) String() string {
	return fmt.Sprintf("StaticBasicObject: %s with %s", t.BasicObject, t.Colorer)
}
