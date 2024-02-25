package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/geometry"
)

type TriangleScene struct {
	// XYRatio      float64
	// SigmoidRatio float64
	// SinCycles    int
	// Gradient     color.Gradient
	t geometry.Triangle
}

func (s TriangleScene) GetColor(x, y, t float64) color.Color {
	triangleColor := s.t.GetColor(x, y)
	if triangleColor != nil {
		return *triangleColor
	}
	return color.White
}

func (s TriangleScene) GetColorPalette(t float64) []color.Color {
	return []color.Color{color.White, color.Black}
}

// func DummyTriangle() TriangleScene {
// 	return TriangleScene{
// 		geometry.Triangle{
// 			geometry.Point{-0.5, -0.5, -1.0},
// 			geometry.Point{-0.5, 0.5, -1.0},
// 			geometry.Point{0.5, -0.5, -1.0},
// 			// color.Black,
// 			// color.Black,
// 			color.Hex("#6CB4F5"),
// 			color.Hex("#EBF56C"),
// 			color.Black,
// 		},
// 	}
// }

func DummyTriangle() Scene {
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
		},
		// Background: Uniform{color.White},
		Background: SineWaveWCross{
			XYRatio:      0.0001,
			SigmoidRatio: 2.0,
			SinCycles:    3,
			TScale:       0.3,
			Gradient: color.LinearGradient{
				Points: []color.Color{
					color.Hex("#FFF"), // black
					color.Hex("#DDF522"),
					color.Hex("#A0514C"),
					color.Hex("#000"), // white
				},
			},
		},
	}
}

type CombinedScene struct {
	Objects    []Object
	Background Scene
}

func (s CombinedScene) GetColor(x, y, t float64) color.Color {
	for _, obj := range s.Objects {
		c := obj.GetColor(x, y)
		if c != nil {
			return *c
		}
	}
	return s.Background.GetColor(x, y, t)
}

func (s CombinedScene) GetColorPalette(t float64) []color.Color {
	return []color.Color{color.White, color.Black} // TODO UPDATE
}
