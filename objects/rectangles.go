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
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.TriangleGradientTexture(colorA, colorB, colorC))),
		),
		DynamicTriangle(
			Triangle{
				A: d,
				B: c,
				C: b,
			},
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.TriangleGradientTexture(colorD, colorC, colorB))),
		),
	)
}

func Parallelogram(a, b, c geometry.Point, texture colors.DynamicTransparentTexture) DynamicObject {
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

func RectanglesAlongPath(path geometry.BezierPath, n int, size float64, texture colors.DynamicTransparentTexture) DynamicObject {
	upVector := geometry.Vector3D{X: 0, Y: 1, Z: 0}
	objects := []DynamicObject{}
	for i := range n {
		t := float64(i) / float64(n-1)
		direction := path.GetDirection(t)
		normalVector := direction.ForwardVector.CrossProduct(upVector).Unit()
		relaltiveUpVector := normalVector.CrossProduct(direction.ForwardVector).Unit()
		a := geometry.Point(direction.Origin.Vector().AddVector(normalVector.ScalarMultiply(-size)).AddVector(relaltiveUpVector.ScalarMultiply(-size)))
		b := geometry.Point(direction.Origin.Vector().AddVector(normalVector.ScalarMultiply(size)).AddVector(relaltiveUpVector.ScalarMultiply(-size)))
		c := geometry.Point(direction.Origin.Vector().AddVector(normalVector.ScalarMultiply(-size)).AddVector(relaltiveUpVector.ScalarMultiply(size)))
		objects = append(objects, Parallelogram(a, b, c, texture))

	}
	return CombineDynamicObjects(objects...)
}
