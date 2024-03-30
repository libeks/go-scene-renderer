package main

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/scenes"
)

var (
	// scene = scenes.BackgroundScene(
	// 	scenes.BackgroundFromTexture(colors.FuzzyDynamic{
	// 		colors.StaticTexture(
	// 			colors.VerticalGradient{
	// 				colors.LinearGradient{
	// 					[]colors.Color{
	// 						// einstein on the beach kind of gradient
	// 						colors.Black,
	// 						colors.Hex("#737C93"),
	// 						colors.Hex("#7C8099"),
	// 						colors.Hex("#72727A"),
	// 						colors.Hex("#9E9796"),
	// 						colors.Hex("#FCFDF7"),
	// 						colors.Hex("#FBEBE1"),
	// 						colors.Hex("#999596"),
	// 						colors.Hex("#5E6A6A"),
	// 						colors.Black,
	// 					},
	// 				},
	// 			},
	// 		),
	// 		0.003,
	// 	},
	// 	),
	// )
	// scene = scenes.BackgroundScene(
	// 	scenes.BackgroundFromTexture(colors.DynamicFromAnimatedTexture(
	// 		colors.NewPerlinNoiseTexture(colors.SimpleGradient{colors.White, colors.Black}),
	// 	),
	// 	),
	// )
	// scene = scenes.BackgroundScene(
	// 	scenes.BackgroundFromTexture(colors.DynamicFromAnimatedTexture(
	// 		colors.DynamicSubtexturer{
	// 			colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
	// 			100,
	// 			colors.NewPerlinNoise(),
	// 		},
	// 	),
	// 	),
	// )
	// RotatingLine
	// scene = scenes.SineWaveWCross{
	// 	XYRatio:      0.0001,
	// 	SigmoidRatio: 2.0,
	// 	SinCycles:    3,
	// 	TScale:       0.3,
	// 	// TOffset:      0.0,
	// 	// Gradient:     color.Grayscale,
	// 	Gradient: gradient,
	// }
	// scene = scenes.DummyTextureSpinningCube(
	// 	colors.DynamicFromAnimatedTexture(colors.DynamicSubtexturer{
	// 		colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
	// 		8,
	// 		colors.NewRandomPerlinNoise(),
	// 	}),
	// 	scenes.BackgroundFromTexture(colors.DynamicFromAnimatedTexture(colors.DynamicSubtexturer{
	// 		colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
	// 		32,
	// 		colors.NewRandomPerlinNoise(),
	// 	})),
	// )

	// scene = scenes.MulticubeDance(
	// 	// scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Black})),
	// 	colors.SimpleGradient{colors.Black, colors.Red},
	// 	colors.SimpleGradient{colors.Black, colors.Green},
	// 	colors.SimpleGradient{colors.Black, colors.Blue},
	// 	scenes.BackgroundFromTexture(
	// 		colors.DynamicFromAnimatedTexture(
	// 			colors.SineWaveAnimation{
	// 				XYRatio:      0.1,
	// 				SigmoidRatio: 2,
	// 				SinCycles:    3,
	// 				Gradient:     colors.Grayscale.Subsample(0.4, 0.6),
	// 			}),
	// 	),
	// )

	// scene = scenes.SpinningMulticube(
	// 	scenes.BackgroundFromTexture(
	// 		colors.DynamicFromAnimatedTexture(
	// 			colors.SineWaveAnimation{
	// 				XYRatio:      0.1,
	// 				SigmoidRatio: 2,
	// 				SinCycles:    3,
	// 				Gradient:     colors.Grayscale,
	// 			}),
	// 	),
	// )
	// scene = scenes.CheckerboardSquare(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Blue})))
	// scene = scenes.SingleSpinningTriangle(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Blue})))
	scene = scenes.HeightMap(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{colors.Blue})))
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
