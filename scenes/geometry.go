package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
)

func DummySpinningCubes(background DynamicBackground) DynamicScene {
	initialCube := UnitRGBCube()
	// initialCube := UnitCube(
	// 	color.Black,
	// 	color.White,
	// 	color.White,
	// 	color.White,
	// 	color.White,
	// 	color.White,
	// 	color.Black,
	// 	color.White,
	// )
	// diagonalCube := initialCube.ApplyMatrix(geometry.RotateMatrixX(-0.615).MatrixMult(
	// 	geometry.RotateMatrixZ(math.Pi / 4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
	// )) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)
	diagonalCube := initialCube

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, math.Sqrt(3) / 2, -3}),
					geometry.RotateMatrixY(t*maths.Rotation),
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, -math.Sqrt(3) / 2, -3}),
					geometry.RotateMatrixY(-t*maths.Rotation),
				)
			}),
		},
		Background: background,
	}
}

func SpinningMulticube(background DynamicBackground) DynamicScene {
	initialCube := UnitRGBCube()
	diagonalCube := initialCube.WithTransform(geometry.MatrixProduct(
		geometry.RotateMatrixX(-0.615),
		geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
	)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	// diagonalCube := initialCube
	// spacing := math.Sqrt(3)
	spacing := 2.0

	column := objects.CombineDynamicObjects(
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, -2 * spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, -spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 2 * spacing, 0})),
	)

	slice := objects.CombineDynamicObjects(
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{-2 * spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{2 * spacing, 0, 0})),
	)

	multiCube := objects.CombineDynamicObjects(
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2 * spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 2 * spacing})),
	)

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			multiCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -10}),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
				)
			},
			),
		},
		Background: background,
	}
}

func NoiseTest() DynamicScene {
	texture := colors.StaticTexture(colors.NewPerlinNoise(colors.Grayscale))
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.Parallelogram(geometry.Point{0, 0, -5}, geometry.Point{2, 0, -5}, geometry.Point{0, 2, -5}, texture),
		},
		Background: BackgroundFromTexture(texture),
	}
}

func SpinningIndividualMulticube(background DynamicBackground) DynamicScene {
	// initialCube := UnitRGBCube()
	// texture := color.SquareGradientTexture(color.White, color.Red, color.Black, color.Blue)
	texture := colors.StaticTexture(colors.NewPerlinNoise(colors.Grayscale))
	initialCube := UnitTextureCube(
		texture,
		texture,
		texture,
		texture,
		texture,
		texture,
	)
	diagonalCube := initialCube.WithTransform(
		geometry.MatrixProduct(
			geometry.RotateMatrixX(-0.615),
			geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
		)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	spacing := 2.0

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),       // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                      // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0}), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),         // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}), // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),  // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),   // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),      // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                     // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0}), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),        // rotation around own axis
				)
			}),
		},
		Background: background,
	}
}

func DummySpinningCube(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			UnitRGBCube().WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 1, -2}),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
					geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
				)
			}),
		},
		Background: background,
	}
}

func CheckerboardSquare(background DynamicBackground) DynamicScene {
	texture := colors.StaticTexture(colors.Checkerboard{8})
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			// objects.DynamicObjectFromTriangles(
			// 	objects.DynamicTriangle{
			// 		Triangle: objects.Triangle{
			// 			A: geometry.Point{0, 0, 0},
			// 			B: geometry.Point{1, 0, 0},
			// 			C: geometry.Point{0, 1, 0},
			// 		},
			// 		Colorer: texture,
			// 	},
			// ).WithDynamicTransform(
			// 	func(t float64) geometry.HomogeneusMatrix {
			// 		return geometry.MatrixProduct(
			// 			geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),
			// 			geometry.RotateMatrixX(t*maths.Rotation),
			// 			// geometry.TranslationMatrix(geometry.Vector3D{0, 0, 5}),
			// 		)
			// 	},
			// ),
			// objects.DynamicObjectFromTriangles(
			// 	objects.DynamicTriangle{
			// 		Triangle: objects.Triangle{
			// 			A: geometry.Point{0, 0, 0},
			// 			B: geometry.Point{1, 0, 0},
			// 			C: geometry.Point{1, 1, 0},
			// 		},
			// 		Colorer: texture,
			// 	},
			// ).WithDynamicTransform(
			// 	func(t float64) geometry.HomogeneusMatrix {
			// 		return geometry.MatrixProduct(
			// 			geometry.TranslationMatrix(geometry.Vector3D{2, 0, -5}),
			// 			geometry.RotateMatrixX(t*maths.Rotation),
			// 			// geometry.TranslationMatrix(geometry.Vector3D{0, 0, 5}),
			// 		)
			// 	},
			// ),
			// objects.DynamicObjectFromTriangles(
			// 	objects.DynamicTriangle{
			// 		Triangle: objects.Triangle{
			// 			A: geometry.Point{0, 0, 0},
			// 			B: geometry.Point{1, 0, 0},
			// 			C: geometry.Point{-1, 1, 0},
			// 		},
			// 		Colorer: texture,
			// 	},
			// ).WithDynamicTransform(
			// 	func(t float64) geometry.HomogeneusMatrix {
			// 		return geometry.MatrixProduct(
			// 			geometry.TranslationMatrix(geometry.Vector3D{-2, 0, -5}),
			// 			geometry.RotateMatrixX(t*maths.Rotation),
			// 			// geometry.TranslationMatrix(geometry.Vector3D{0, 0, 5}),
			// 		)
			// 	},
			// ),
			objects.Parallelogram(
				geometry.Point{0, 0, 0},
				geometry.Point{2, 0, 0},
				geometry.Point{-2, 2, 0},
				texture).WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),
						geometry.RotateMatrixX(t*maths.Rotation),
						// geometry.TranslationMatrix(geometry.Vector3D{0, 0, 5}),
					)
				},
			),
		},
		Background: background,
	}
}
