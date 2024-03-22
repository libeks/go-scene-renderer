package scenes

import (
	"fmt"
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
)

type TriangleScene struct {
	t objects.Triangle
}

func (s TriangleScene) GetFrameColor(x, y, t float64) colors.Color {
	triangleColor, _ := s.t.GetColorDepth(x, y)
	if triangleColor != nil {
		return *triangleColor
	}
	return colors.White
}

func (s TriangleScene) GetColorPalette(t float64) []colors.Color {
	return []colors.Color{colors.White, colors.Black}
}

type DynamicTriangle struct {
	t objects.Triangle
}

func SpinningTriangle(tri objects.Triangle) DynamicTriangle {
	return DynamicTriangle{
		tri,
	}
}

func (s DynamicTriangle) GetFrame(t float64) objects.Object {
	matrix := geometry.TranslationMatrix(geometry.Vector3D{
		0, 0, -2,
	}).MatrixMult(geometry.RotateMatrixY(t * maths.Rotation))
	// fmt.Printf("At t=%.3f the matrix is %s\n", t, matrix)
	return s.t.ApplyMatrix(matrix)
}

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

	fmt.Printf("DiagonalCube: %s\n", diagonalCube)
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, math.Sqrt(3) / 2, -3}),
						geometry.RotateMatrixY(t*maths.Rotation),
					)
				},
			},
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, -math.Sqrt(3) / 2, -3}),
						geometry.RotateMatrixY(-t*maths.Rotation),
					)
				},
			},
		},
		Background: background,
	}
}

func SpinningMulticube(background DynamicBackground) DynamicScene {
	initialCube := UnitRGBCube()
	diagonalCube := initialCube.ApplyMatrix(geometry.RotateMatrixX(-0.615).MatrixMult(
		geometry.RotateMatrixZ(math.Pi / 4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
	)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	// diagonalCube := initialCube
	// spacing := math.Sqrt(3)
	spacing := 2.0

	column := objects.ComplexObject{
		[]objects.TransformableObject{
			diagonalCube.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, -2 * spacing, 0})),
			diagonalCube.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, -spacing, 0})),
			diagonalCube.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
			diagonalCube.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, spacing, 0})),
			diagonalCube.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 2 * spacing, 0})),
		},
	}

	slice := objects.ComplexObject{
		Objs: []objects.TransformableObject{
			column.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{-2 * spacing, 0, 0})),
			column.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0})),
			column.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
			column.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0})),
			column.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{2 * spacing, 0, 0})),
		},
	}

	multiCube := objects.ComplexObject{
		Objs: []objects.TransformableObject{
			slice.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2 * spacing})),
			slice.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -spacing})),
			slice.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
			slice.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, spacing})),
			slice.ApplyMatrix(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 2 * spacing})),
		},
	}

	// scene :=

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: multiCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -10}),
						geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					)
				},
			},
		},
		Background: background,
	}
}

func NoiseTest() DynamicScene {
	texture := colors.NewPerlinNoise(colors.Grayscale)
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: objects.Parallelogram(geometry.Point{0, 0, -5}, geometry.Point{2, 0, -5}, geometry.Point{0, 2, -5}, texture),
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct()
				},
			},
		},
		Background: BackgroundFromTexture(colors.StaticTexture(texture)),
	}
}

func SpinningIndividualMulticube(background DynamicBackground) DynamicScene {
	// initialCube := UnitRGBCube()
	// texture := color.SquareGradientTexture(color.White, color.Red, color.Black, color.Blue)
	texture := colors.NewPerlinNoise(colors.Grayscale)
	initialCube := UnitTextureCube(
		texture,
		texture,
		texture,
		texture,
		texture,
		texture,
	)
	diagonalCube := initialCube.ApplyMatrix(
		geometry.MatrixProduct(
			geometry.RotateMatrixX(-0.615),
			geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
		)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	spacing := 2.0

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),       // position within the scene
						geometry.RotateMatrixY(t*maths.Rotation),                      // rotation around common center
						geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0}), // position within the group
						geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),         // rotation around own axis
					)
				},
			},
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}), // position within the scene
						geometry.RotateMatrixY(t*maths.Rotation),                // rotation around common center
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),  // position within the group
						geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),   // rotation around own axis
					)
				},
			},
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),      // position within the scene
						geometry.RotateMatrixY(t*maths.Rotation),                     // rotation around common center
						geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0}), // position within the group
						geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),        // rotation around own axis
					)
				},
			},
		},
		Background: background,
	}
}

func DummySpinningCube(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: UnitRGBCube(),
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2}),
						geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
						geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
						geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
					)
				},
			},
		},
		Background: background,
	}
}

func DummySpinningTriangle() DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: objects.GradientTriangle(
					geometry.Point{-0.5, -0.5, -1.0},
					geometry.Point{-0.5, 0.5, -1.0},
					geometry.Point{0.5, -0.5, -1.0},
					colors.Hex("#6CB4F5"),
					colors.Hex("#EBF56C"),
					colors.Black,
				),
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2}),
						geometry.RotateMatrixY(t*maths.Rotation),
					)
				},
			},
		},
		Background: BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Black})),
	}
}
