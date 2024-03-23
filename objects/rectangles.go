package objects

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
)

// fourth coordinate is inferred from the first three
func GradientParallelogram(a, b, c geometry.Point, colorA, colorB, colorC, colorD colors.Color) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromTriangles(
		DynamicTriangle{
			Triangle: Triangle{
				A: a,
				B: b,
				C: c,
			},
			Colorer: colors.StaticTexture(colors.TriangleGradientTexture(colorA, colorB, colorC)),
		},
		DynamicTriangle{
			Triangle: Triangle{
				A: d,
				B: c,
				C: b,
			},
			Colorer: colors.StaticTexture(colors.TriangleGradientTexture(colorD, colorC, colorB)),
		},
	)
}

func Parallelogram(a, b, c geometry.Point, texture colors.DynamicTexture) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromTriangles(
		DynamicTriangle{
			Triangle: Triangle{
				A: a,
				B: b,
				C: c,
			},
			Colorer: texture,
		},
		DynamicTriangle{
			Triangle: Triangle{
				A: d,
				B: c,
				C: b,
			},
			Colorer: colors.RotateDynamicTexture180(texture),
		},
	)
}
