package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/textures"
)

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

func BackgroundScene(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Background: background,
	}
}

// returns a unit cube, with textures applied
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitTextureCube(t1, t2, t3, t4, t5, t6 textures.DynamicTransparentTexture) objects.DynamicObject {
	return objects.CombineDynamicObjects(
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 1, 0),
			geometry.Pt(1, 0, 0),
			t1,
		),
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 0, 1),
			geometry.Pt(1, 0, 0),
			t2,
		),
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 0, 1),
			geometry.Pt(0, 1, 0),
			t3,
		),

		// halfway

		objects.Parallelogram(
			geometry.Pt(0, 0, 1),
			geometry.Pt(0, 1, 1),
			geometry.Pt(1, 0, 1),
			t4,
		),
		objects.Parallelogram(
			geometry.Pt(0, 1, 0),
			geometry.Pt(0, 1, 1),
			geometry.Pt(1, 1, 0),
			t5,
		),
		objects.Parallelogram(
			geometry.Pt(1, 0, 0),
			geometry.Pt(1, 0, 1),
			geometry.Pt(1, 1, 0),
			t6,
		),
	).WithTransform(geometry.TranslationMatrix(
		geometry.V3(
			-0.5, -0.5, -0.5,
		),
	))
}

// returns a unit cube, with textures applied
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitTextureCubeWithTransparency(t1, t2, t3, t4, t5, t6 textures.DynamicTexture, alpha textures.DynamicTransparency) objects.DynamicObject {
	return objects.CombineDynamicObjects(
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 1, 0),
			geometry.Pt(1, 0, 0),
			textures.GetDynamicTransparentTexture(
				t1,
				alpha,
			),
		),
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 0, 1),
			geometry.Pt(1, 0, 0),
			textures.GetDynamicTransparentTexture(
				t2,
				alpha,
			),
		),
		objects.Parallelogram(
			geometry.Pt(0, 0, 0),
			geometry.Pt(0, 0, 1),
			geometry.Pt(0, 1, 0),
			textures.GetDynamicTransparentTexture(
				t3,
				alpha,
			),
		),

		// halfway

		objects.Parallelogram(
			geometry.Pt(0, 0, 1),
			geometry.Pt(0, 1, 1),
			geometry.Pt(1, 0, 1),
			textures.GetDynamicTransparentTexture(
				t4,
				alpha,
			),
		),
		objects.Parallelogram(
			geometry.Pt(0, 1, 0),
			geometry.Pt(0, 1, 1),
			geometry.Pt(1, 1, 0),
			textures.GetDynamicTransparentTexture(
				t5,
				alpha,
			),
		),
		objects.Parallelogram(
			geometry.Pt(1, 0, 0),
			geometry.Pt(1, 0, 1),
			geometry.Pt(1, 1, 0),
			textures.GetDynamicTransparentTexture(
				t6,
				alpha,
			),
		),
	).WithTransform(geometry.TranslationMatrix(
		geometry.V3(
			-0.5, -0.5, -0.5,
		),
	))
}

// returns a unit cube, with textures applied
// it is centered on the origin point, having sizes of length 1
// so one corner is (-0.5, -0.5, -0.5) and the opposite one is (0.5, 0.5, 0.5), etc
func UnitGradientCube(c000, c100, c110, c010, c001, c101, c111, c011 colors.Color) objects.DynamicObject {
	return UnitTextureCube(
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c000, c010, c100, c110))),
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c000, c001, c100, c101))),
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c000, c001, c010, c011))),

		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c001, c011, c101, c111))),
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c010, c011, c110, c111))),
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.SquareGradientTexture(c100, c101, c110, c111))),
	)
}
