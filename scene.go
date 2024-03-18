package main

import (
	"github.com/libeks/go-scene-renderer/scenes"
)

var (
	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#6CB4F5"),
	// 		color.Hex("#EBF56C"),
	// 		color.Hex("#F5736C"),
	// 	},
	// }
	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#F590C1"), // pink
	// 		color.Hex("#000"),
	// 		color.Hex("#90E8F5"), // light blue
	// 		color.Hex("#000"),
	// 		color.Hex("#F590C1"), // pink
	// 	},
	// }

	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#FFF"), // black
	// 		color.Hex("#DDF522"),
	// 		color.Hex("#A0514C"),
	// 		color.Hex("#000"), // white
	// 	},
	// }
	// scene = scenes.SineWaveWCross{
	// 	XYRatio:      0.0001,
	// 	SigmoidRatio: 2.0,
	// 	SinCycles:    3,
	// 	TScale:       0.3,
	// 	// TOffset:      0.0,
	// 	// Gradient:     color.Grayscale,
	// 	Gradient: gradient,
	// }
	scene = scenes.DummySpinningCube(scenes.Random{})
	// scene = scenes.SpinningMulticube(scenes.SineWave{
	// 	XYRatio:      0.1,
	// 	SigmoidRatio: 2,
	// 	SinCycles:    3,
	// 	Gradient:     color.Grayscale,
	// })
	// scene = scenes.NewPerlinNoise(color.Grayscale)
	// scene = scenes.DummyTriangle()

	// scene = scenes.HorizGradient{
	// 	Gradient: gradient,
	// }

	// scene = scenes.SineWave{
	// 	XYRatio:      0.1,
	// 	SigmoidRatio: 2,
	// 	SinCycles:    3,
	// 	Gradient:     color.Grayscale,
	// }
)
