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
	matrix := geometry.RotateMatrixZ(t * (2 * math.Pi))
	// fmt.Printf("At t=%.3f the matrix is %s\n", t, matrix)
	return s.t.ApplyMatrix(matrix)
}

func DummySpinningTriangle() DynamicScene {
	return CombinedDynamicScene{
		Objects: []DynamicObject{SpinningTriangle(
			geometry.Triangle{
				geometry.Point{-0.5, -0.5, -1.0},
				geometry.Point{-0.5, 0.5, -1.0},
				geometry.Point{0.5, -0.5, -1.0},
				color.Hex("#6CB4F5"),
				color.Hex("#EBF56C"),
				color.Black,
			},
		)},
		Background: Uniform{color.White},
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
