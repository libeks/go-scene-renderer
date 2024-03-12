package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
)

type TriangleScene struct {
	t geometry.Triangle
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
	t geometry.Triangle
}

func SpinningTriangle(tri geometry.Triangle) DynamicTriangle {
	return DynamicTriangle{
		tri,
	}
}

func (s DynamicTriangle) GetFrame(t float64) Object {
	matrix := geometry.TranslationMatrix(geometry.Vector3D{
		0, 0, -2,
	}).MatrixMult(geometry.RotateMatrixY(t * (2 * math.Pi)))
	// fmt.Printf("At t=%.3f the matrix is %s\n", t, matrix)
	return s.t.ApplyMatrix(matrix)
}

type TransformedObject struct {
	x        TransformableObject
	matrixFn func(t float64) geometry.HomogeneusMatrix
}

func (o TransformedObject) GetFrame(t float64) Object {
	m := o.matrixFn(t)
	return o.x.ApplyMatrix(m)
}

type ComplexObject struct {
	triangles []geometry.Triangle
}

func (o ComplexObject) GetColorDepth(x, y float64) (*color.Color, float64) {
	minZ := math.MaxFloat64
	var closestColor *color.Color
	for _, obj := range o.triangles {
		c, depth := obj.GetColorDepth(x, y)
		if c != nil && depth < minZ {
			minZ = depth
			closestColor = c
		}
	}
	if closestColor != nil {
		return closestColor, minZ
	}
	return nil, 0
}

func (o ComplexObject) ApplyMatrix(m geometry.HomogeneusMatrix) TransformableObject {
	newTriangles := make([]geometry.Triangle, len(o.triangles))
	for i, triangle := range o.triangles {
		newTriangles[i] = triangle.ApplyMatrix(m)
	}
	return ComplexObject{
		triangles: newTriangles,
	}
}

func UnitCube() TransformableObject {
	return ComplexObject{
		triangles: []geometry.Triangle{
			geometry.Triangle{
				geometry.Point{0, 0, 0},
				geometry.Point{0, 1, 0},
				geometry.Point{1, 0, 0},
				color.Black,
				color.Green,
				color.Red,
			},
			geometry.Triangle{
				geometry.Point{1, 1, 0},
				geometry.Point{0, 1, 0},
				geometry.Point{1, 0, 0},
				color.Yellow,
				color.Green,
				color.Red,
			},

			geometry.Triangle{
				geometry.Point{0, 0, 0},
				geometry.Point{0, 0, 1},
				geometry.Point{1, 0, 0},
				color.Black,
				color.Blue,
				color.Red,
			},
			geometry.Triangle{
				geometry.Point{1, 0, 1},
				geometry.Point{0, 0, 1},
				geometry.Point{1, 0, 0},
				color.Magenta,
				color.Blue,
				color.Red,
			},

			geometry.Triangle{
				geometry.Point{0, 0, 0},
				geometry.Point{0, 0, 1},
				geometry.Point{0, 1, 0},
				color.Black,
				color.Blue,
				color.Green,
			},
			geometry.Triangle{
				geometry.Point{0, 1, 1},
				geometry.Point{0, 0, 1},
				geometry.Point{0, 1, 0},
				color.Cyan,
				color.Blue,
				color.Green,
			},

			// halfway

			geometry.Triangle{
				geometry.Point{0, 0, 1},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 0, 1},
				color.Blue,
				color.Cyan,
				color.Magenta,
			},
			geometry.Triangle{
				geometry.Point{1, 1, 1},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 0, 1},
				color.White,
				color.Cyan,
				color.Magenta,
			},

			geometry.Triangle{
				geometry.Point{0, 1, 0},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 1, 0},
				color.Green,
				color.Cyan,
				color.Yellow,
			},
			geometry.Triangle{
				geometry.Point{1, 1, 1},
				geometry.Point{0, 1, 1},
				geometry.Point{1, 1, 0},
				color.White,
				color.Cyan,
				color.Yellow,
			},

			geometry.Triangle{
				geometry.Point{1, 0, 0},
				geometry.Point{1, 0, 1},
				geometry.Point{1, 1, 0},
				color.Red,
				color.Magenta,
				color.Yellow,
			},
			geometry.Triangle{
				geometry.Point{1, 1, 1},
				geometry.Point{1, 0, 1},
				geometry.Point{1, 1, 0},
				color.White,
				color.Magenta,
				color.Yellow,
			},
		},
	}.ApplyMatrix(geometry.TranslationMatrix(
		geometry.Vector3D{
			-0.5, -0.5, -0.5,
		},
	))
}

func DummySpinningCube(background DynamicScene) DynamicScene {
	return CombinedDynamicScene{
		Objects: []DynamicObject{
			TransformedObject{
				UnitCube(),
				func(t float64) geometry.HomogeneusMatrix {
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
		Objects: []DynamicObject{
			TransformedObject{
				ComplexObject{
					[]geometry.Triangle{
						geometry.Triangle{
							geometry.Point{-0.5, -0.5, -1.0},
							geometry.Point{-0.5, 0.5, -1.0},
							geometry.Point{0.5, -0.5, -1.0},
							color.Hex("#6CB4F5"),
							color.Hex("#EBF56C"),
							color.Black,
						},
					},
				},
				func(t float64) geometry.HomogeneusMatrix {
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
		Objects: []Object{
			geometry.Triangle{
				geometry.Point{-0.5, -0.5, -1.0},
				geometry.Point{-0.5, 0.5, -1.0},
				geometry.Point{0.5, -0.5, -1.0},
				color.Hex("#6CB4F5"),
				color.Hex("#EBF56C"),
				color.Black,
			},
			geometry.Triangle{
				geometry.Point{-1.0, -1.0, -1.2},
				geometry.Point{-1.0, 1.0, -1.2},
				geometry.Point{1.0, -1.0, -1.2},
				color.Red,
				color.Hex("#EBF56C"),
				color.White,
			},
			geometry.Triangle{
				geometry.Point{1.0, 1.0, -0.9},
				geometry.Point{-1.0, 1.0, -1.1},
				geometry.Point{1.0, -1.0, -1.1},
				color.Hex("#90E8F5"),
				color.Hex("#EBF56C"),
				color.Hex("#F590C1"),
			},
			geometry.Triangle{
				geometry.Point{0.75, 0.75, -1.0},
				geometry.Point{-0.75, 0.75, -1.0},
				geometry.Point{0.75, -0.75, -1.0},
				color.Hex("#0F0"),
				color.Hex("#F00"),
				color.Hex("#00F"),
			},
			geometry.Triangle{
				geometry.Point{0.5, 0.5, -0.5},
				geometry.Point{-0.5, 0.5, -2.0},
				geometry.Point{0.5, -0.5, -2.0},
				color.Hex("#0Ff"),
				color.Hex("#F00"),
				color.Hex("#00F"),
			},
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

type CombinedScene struct {
	Objects    []Object
	Background Frame
}

func (s CombinedScene) GetColor(x, y float64) color.Color {
	minZ := math.MaxFloat64
	var closestColor *color.Color
	for _, obj := range s.Objects {
		c, depth := obj.GetColorDepth(x, y)
		if c != nil && depth < minZ {
			minZ = depth
			closestColor = c
		}
	}
	if closestColor != nil {
		return *closestColor
	}
	return s.Background.GetColor(x, y)
}

type CombinedDynamicScene struct {
	Objects    []DynamicObject
	Background DynamicScene
}

func (s CombinedDynamicScene) GetFrame(t float64) Frame {
	frameObjects := make([]Object, len(s.Objects))
	for i, object := range s.Objects {
		frameObjects[i] = object.GetFrame(t)
	}
	return CombinedScene{
		Objects:    frameObjects,
		Background: s.Background.GetFrame(t),
	}
}