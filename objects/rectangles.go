package objects

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
)

// fourth coordinate is inferred from the first three
func GradientParallelogram(a, b, c geometry.Point, colorA, colorB, colorC, colorD colors.Color) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromTriangles(
		DynamicTriangle(
			Triangle{
				A: a,
				B: b,
				C: c,
			},
			colors.StaticTexture(colors.TriangleGradientTexture(colorA, colorB, colorC)),
		),
		DynamicTriangle(
			Triangle{
				A: d,
				B: c,
				C: b,
			},
			colors.StaticTexture(colors.TriangleGradientTexture(colorD, colorC, colorB)),
		),
	)
}

func Parallelogram(a, b, c geometry.Point, texture colors.DynamicTexture) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromTriangles(
		DynamicTriangle(
			Triangle{
				A: a,
				B: b,
				C: c,
			},
			texture,
		),
		DynamicTriangle(
			Triangle{
				A: d,
				B: c,
				C: b,
			},
			colors.RotateDynamicTexture180(texture),
		),
	)
}

func ParallelogramWithTransparency(a, b, c geometry.Point, texture colors.DynamicTexture, transparency colors.DynamicTransparency) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromTriangles(
		DynamicTriangleWithTransparency(
			Triangle{
				A: a,
				B: b,
				C: c,
			},
			texture,
			transparency,
		),
		DynamicTriangleWithTransparency(
			Triangle{
				A: d,
				B: c,
				C: b,
			},
			colors.RotateDynamicTexture180(texture),
			transparency, // TODO: flip transparency around, flip it around
		),
	)
}
