package main

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/scenes"
)

var (
	//	scene = scenes.SineWaveWCross{
	//		XYRatio:      0.0001,
	//		SigmoidRatio: 2.0,
	//		SinCycles:    3,
	//		TScale:       0.3,
	//		// TOffset:      0.0,
	//		// Gradient:     color.Grayscale,
	//		Gradient: gradient,
	//	}
	// scene = scenes.DummySpinningCube(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Black})))

	scene = scenes.SpinningMulticube(
		scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.SineWaveAnimation{
					XYRatio:      0.1,
					SigmoidRatio: 2,
					SinCycles:    3,
					Gradient:     colors.Grayscale,
				}),
		),
	)
	// scene = scenes.CheckerboardSquare(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Blue})))
)

// scene = scenes.NoiseTest()
// scene = scenes.NewPerlinNoise(color.Grayscale)
// scene = scenes.DummyTriangle()

//	scene = scenes.SineWave{
//		XYRatio:      0.1,
//		SigmoidRatio: 2,
//		SinCycles:    3,
//		Gradient:     color.Grayscale,
//	}
