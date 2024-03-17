package scenes

import (
	"fmt"
	"math"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
)

type TriangleScene struct {
	t objects.Triangle
}

func (s TriangleScene) GetFrameColor(x, y, t float64) color.Color {
	triangleColor, _ := s.t.GetColorDepth(x, y)
	if triangleColor != nil {
		return *triangleColor
	}
	return color.White
}

func (s TriangleScene) GetColorPalette(t float64) []color.Color {
	return []color.Color{color.White, color.Black}
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
	}).MatrixMult(geometry.RotateMatrixY(t * (2 * math.Pi)))
	// fmt.Printf("At t=%.3f the matrix is %s\n", t, matrix)
	return s.t.ApplyMatrix(matrix)
}

func DummySpinningCubes(background DynamicScene) DynamicScene {
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
					return geometry.TranslationMatrix(geometry.Vector3D{
						0, math.Sqrt(3) / 2, -3,
					}).MatrixMult(
						geometry.RotateMatrixY(t * (2 * math.Pi)),
					)
				},
			},
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.TranslationMatrix(geometry.Vector3D{
						0, -math.Sqrt(3) / 2, -3,
					}).MatrixMult(
						geometry.RotateMatrixY(-t * (2 * math.Pi)),
					)
				},
			},
		},
		Background: background,
	}
}

func DummySpinningCubes2(background DynamicScene) DynamicScene {
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

	// scene :=

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: diagonalCube,
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.TranslationMatrix(geometry.Vector3D{
						0, 0, -3,
					}).MatrixMult(
						geometry.RotateMatrixY(t * (2 * math.Pi)),
					)
				},
			},
		},
		Background: background,
	}
}

func DummySpinningCube(background DynamicScene) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.TransformedObject{
				Object: UnitRGBCube(),
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.TranslationMatrix(geometry.Vector3D{
						0, 0, -2,
					}).MatrixMult(
						geometry.RotateMatrixY(t * (2 * math.Pi)),
					).MatrixMult(
						// arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
						geometry.RotateMatrixX(-0.615).MatrixMult(
							geometry.RotateMatrixZ(math.Pi / 4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
						),
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
					color.Hex("#6CB4F5"),
					color.Hex("#EBF56C"),
					color.Black,
				),
				MatrixFn: func(t float64) geometry.HomogeneusMatrix {
					return geometry.TranslationMatrix(geometry.Vector3D{
						0, 0, -2,
					}).MatrixMult(geometry.RotateMatrixY(t * (2 * math.Pi)))
				},
			},
		},
		Background: Uniform{color.Black},
	}
}

func DummyTriangle() CombinedScene {
	return CombinedScene{
		Objects: []objects.Object{
			objects.GradientTriangle(
				geometry.Point{-0.5, -0.5, -1.0},
				geometry.Point{-0.5, 0.5, -1.0},
				geometry.Point{0.5, -0.5, -1.0},
				color.Hex("#6CB4F5"),
				color.Hex("#EBF56C"),
				color.Black,
			),
			objects.GradientTriangle(
				geometry.Point{-1.0, -1.0, -1.2},
				geometry.Point{-1.0, 1.0, -1.2},
				geometry.Point{1.0, -1.0, -1.2},
				color.Red,
				color.Hex("#EBF56C"),
				color.White,
			),
			objects.GradientTriangle(
				geometry.Point{1.0, 1.0, -0.9},
				geometry.Point{-1.0, 1.0, -1.1},
				geometry.Point{1.0, -1.0, -1.1},
				color.Hex("#90E8F5"),
				color.Hex("#EBF56C"),
				color.Hex("#F590C1"),
			),
			objects.GradientTriangle(
				geometry.Point{0.75, 0.75, -1.0},
				geometry.Point{-0.75, 0.75, -1.0},
				geometry.Point{0.75, -0.75, -1.0},
				color.Hex("#0F0"),
				color.Hex("#F00"),
				color.Hex("#00F"),
			),
			objects.GradientTriangle(
				geometry.Point{0.5, 0.5, -0.5},
				geometry.Point{-0.5, 0.5, -2.0},
				geometry.Point{0.5, -0.5, -2.0},
				color.Hex("#0Ff"),
				color.Hex("#F00"),
				color.Hex("#00F"),
			),
		},
		Background: Uniform{color.White},
		// Background: SineWaveWCross{
		// 	XYRatio:      0.0001,
		// 	SigmoidRatio: 2.0,
		// 	SinCycles:    3,
		// 	TScale:       0.3,
		// 	Gradient: color.LinearGradient{
		// 		Points: []color.Color{
		// 			color.Hex("#FFF"), // black
		// 			color.Hex("#DDF522"),
		// 			color.Hex("#A0514C"),
		// 			color.Hex("#000"), // white
		// 		},
		// 	},
		// },
	}
}
