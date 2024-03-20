package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
)

// returns a unit square in the x-y plane, with colors arranged as indicated by x,y colors in color parameter names
func UnitSquare(c00, c10, c11, c01 color.Color) []*objects.Triangle {
	return []*objects.Triangle{
		objects.GradientTriangle(
			geometry.Point{0, 0, 0},
			geometry.Point{0, 1, 0},
			geometry.Point{1, 0, 0},
			c00,
			c01,
			c10,
		),
		objects.GradientTriangle(
			geometry.Point{1, 1, 0},
			geometry.Point{0, 1, 0},
			geometry.Point{1, 0, 0},
			c11,
			c01,
			c10,
		),
	}
}

func UnitRGBCube() objects.TransformableObject {
	return UnitCube(
		color.Black,
		color.Red,
		color.Yellow,
		color.Green,
		color.Blue,
		color.Magenta,
		color.White,
		color.Cyan,
	)
}

// returns a unit cube, with colors arranged as indicated by x,y,z colors in color parameter names
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitCube(c000, c100, c110, c010, c001, c101, c111, c011 color.Color) objects.TransformableObject {
	return objects.ComplexObject{
		Objs: []objects.TransformableObject{
			objects.GradientParallelogram(
				geometry.Point{0, 0, 0},
				geometry.Point{0, 1, 0},
				geometry.Point{1, 0, 0},
				c000, c010, c100, c110,
			),
			objects.GradientParallelogram(
				geometry.Point{0, 0, 0},
				geometry.Point{0, 0, 1},
				geometry.Point{1, 0, 0},
				c000, c001, c100, c101,
			),
			objects.GradientParallelogram(
				geometry.Point{0, 0, 0},
				geometry.Point{0, 0, 1},
				geometry.Point{0, 1, 0},
				c000, c001, c010, c011,
			),

			// halfway

			objects.GradientParallelogram(
				geometry.Point{0, 0, 1},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 0, 1},
				c001,
				c011,
				c101,
				c111,
			),
			objects.GradientParallelogram(
				geometry.Point{0, 1, 0},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 1, 0},
				c010,
				c011,
				c110,
				c111,
			),
			objects.GradientParallelogram(
				geometry.Point{1, 0, 0},
				geometry.Point{1, 0, 1},
				geometry.Point{1, 1, 0},
				c100,
				c101,
				c110,
				c111,
			),
		},
	}.ApplyMatrix(geometry.TranslationMatrix(
		geometry.Vector3D{
			-0.5, -0.5, -0.5,
		},
	))
}
