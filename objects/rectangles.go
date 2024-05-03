package objects

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
)

// fourth coordinate is inferred from the first three
func GradientParallelogram(a, b, c geometry.Point, colorA, colorB, colorC, colorD colors.Color) DynamicObject {
	d := geometry.Point(c.Add(geometry.Point(b.Subtract(a))))

	return DynamicObjectFromBasics(
		DynamicBasicObject(
			&Triangle{
				A: a,
				B: b,
				C: c,
			},
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.TriangleGradientTexture(colorA, colorB, colorC))),
		),
		DynamicBasicObject(
			&Triangle{
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

	return DynamicObjectFromBasics(
		DynamicBasicObject(
			&Triangle{
				A: a,
				B: b,
				C: c,
			},
			texture,
		),
		DynamicBasicObject(
			&Triangle{
				A: d,
				B: c,
				C: b,
			},
			colors.RotateDynamicTexture180(texture),
		),
	)
}

func RectanglesAlongPath(path geometry.BezierPath, n int, size float64, texture colors.DynamicTransparentTexture) DynamicObject {
	// upVector := geometry.Vector3D{X: 0, Y: 1, Z: 0}
	objects := []DynamicObject{}
	for i := range n {
		t := float64(i) / float64(n-1)
		direction := path.GetDirection(t)
		a := geometry.Point(direction.Origin.Vector().AddVector(direction.Orientation.RightVector.ScalarMultiply(-size)).AddVector(direction.Orientation.UpVector.ScalarMultiply(-size)))
		b := geometry.Point(direction.Origin.Vector().AddVector(direction.Orientation.RightVector.ScalarMultiply(size)).AddVector(direction.Orientation.UpVector.ScalarMultiply(-size)))
		c := geometry.Point(direction.Origin.Vector().AddVector(direction.Orientation.RightVector.ScalarMultiply(-size)).AddVector(direction.Orientation.UpVector.ScalarMultiply(size)))
		// var tt colors.DynamicTransparentTexture
		// if i == n-1 {
		// 	tt = colors.OpaqueDynamicTexture(colors.StaticTexture(colors.Uniform{Color: colors.Red}))
		// 	// scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Black}))
		// } else if i == 0 {
		// 	tt = colors.OpaqueDynamicTexture(colors.StaticTexture(colors.Uniform{Color: colors.Blue}))
		// } else {
		tt := texture
		// }
		objects = append(objects, Parallelogram(a, b, c, tt))

	}
	return CombineDynamicObjects(objects...)
}

// AxisAlignedPointer returns three trianges:
// green points towards +x
// red points towards +y
// blue points towards -z
func AxisAlignedPointer() DynamicObject {
	offset := 0.2
	return DynamicObjectFromBasics(
		DynamicBasicObject(
			&Triangle{
				A: geometry.Pt(0, 0, 0),
				B: geometry.Pt(1, 0, 0),
				C: geometry.Pt(0, offset, 0),
			},
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.HorizontalGradient{Gradient: colors.SimpleGradient{Start: colors.Green, End: colors.Green}})),
		),
		DynamicBasicObject(
			&Triangle{
				A: geometry.Pt(0, 0, 0),
				B: geometry.Pt(0, 1, 0),
				C: geometry.Pt(0, 0, -offset),
			},
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.HorizontalGradient{Gradient: colors.SimpleGradient{Start: colors.Red, End: colors.Red}})),
		),
		DynamicBasicObject(
			&Triangle{
				A: geometry.Pt(0, 0, 0),
				B: geometry.Pt(0, 0, -1),
				C: geometry.Pt(offset, 0, 0),
			},
			colors.OpaqueDynamicTexture(colors.StaticTexture(colors.HorizontalGradient{Gradient: colors.SimpleGradient{Start: colors.Blue, End: colors.Blue}})),
		),
	)
}
