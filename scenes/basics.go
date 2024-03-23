package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
)

// returns a unit square in the x-y plane, with colors arranged as indicated by x,y colors in color parameter names
func UnitSquare(c00, c10, c11, c01 colors.Color) objects.DynamicObject {
	return objects.DynamicObjectFromTriangles(
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
	)
}

func UnitRGBCube() objects.DynamicObject {
	return UnitGradientCube(
		colors.Black,
		colors.Red,
		colors.Yellow,
		colors.Green,
		colors.Blue,
		colors.Magenta,
		colors.White,
		colors.Cyan,
	)
}

// returns a unit cube, with colors arranged as indicated by x,y,z colors in color parameter names
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitGradientCube(c000, c100, c110, c010, c001, c101, c111, c011 colors.Color) objects.DynamicObject {
	return objects.CombineDynamicObjects(
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
			c001, c011, c101, c111,
		),
		objects.GradientParallelogram(
			geometry.Point{0, 1, 0},
			geometry.Point{0, 1, 1},
			geometry.Point{1, 1, 0},
			c010, c011, c110, c111,
		),
		objects.GradientParallelogram(
			geometry.Point{1, 0, 0},
			geometry.Point{1, 0, 1},
			geometry.Point{1, 1, 0},
			c100, c101, c110, c111,
		),
	).WithTransform(geometry.TranslationMatrix(
		geometry.Vector3D{
			-0.5, -0.5, -0.5,
		},
	))
}

// returns a unit cube, with textures applied
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitTextureCube(t1, t2, t3, t4, t5, t6 colors.DynamicTexture) objects.DynamicObject {
	return objects.CombineDynamicObjects(
		objects.Parallelogram(
			geometry.Point{0, 0, 0},
			geometry.Point{0, 1, 0},
			geometry.Point{1, 0, 0},
			t1,
		),
		objects.Parallelogram(
			geometry.Point{0, 0, 0},
			geometry.Point{0, 0, 1},
			geometry.Point{1, 0, 0},
			t2,
		),
		objects.Parallelogram(
			geometry.Point{0, 0, 0},
			geometry.Point{0, 0, 1},
			geometry.Point{0, 1, 0},
			t3,
		),

		// halfway

		objects.Parallelogram(
			geometry.Point{0, 0, 1},
			geometry.Point{0, 1, 1},
			geometry.Point{1, 0, 1},
			t4,
		),
		objects.Parallelogram(
			geometry.Point{0, 1, 0},
			geometry.Point{0, 1, 1},
			geometry.Point{1, 1, 0},
			t5,
		),
		objects.Parallelogram(
			geometry.Point{1, 0, 0},
			geometry.Point{1, 0, 1},
			geometry.Point{1, 1, 0},
			t6,
		),
	).WithTransform(geometry.TranslationMatrix(
		geometry.Vector3D{
			-0.5, -0.5, -0.5,
		},
	))
}
