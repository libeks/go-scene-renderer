package objects

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
)

// fourth coordinate is inferred from the first three
func GradientParallelogram(a, b, c geometry.Point, colorA, colorB, colorC, colorD colors.Color) ComplexObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return ComplexObject{
		Objs: []TransformableObject{
			&Triangle{
				A:       a,
				B:       b,
				C:       c,
				Colorer: colors.TriangleGradientTexture(colorA, colorB, colorC),
			},
			&Triangle{
				A:       d,
				B:       c,
				C:       b,
				Colorer: colors.TriangleGradientTexture(colorD, colorC, colorB),
			},
		},
	}
}

func Parallelogram(a, b, c geometry.Point, texture colors.Texture) ComplexObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return ComplexObject{
		Objs: []TransformableObject{
			&Triangle{
				A:       a,
				B:       b,
				C:       c,
				Colorer: texture,
			},
			&Triangle{
				A:       d,
				B:       c,
				C:       b,
				Colorer: colors.RotateTexture180(texture),
			},
		},
	}
}
