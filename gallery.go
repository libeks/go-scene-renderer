package main

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/sampler"
	"github.com/libeks/go-scene-renderer/scenes"
)

var (
	EinsteinOnTheBeach = scenes.BackgroundScene(
		scenes.BackgroundFromTexture(colors.FuzzyDynamic{
			Texture: colors.StaticTexture(
				colors.VerticalGradient{
					Gradient: colors.LinearGradient{
						Points: []colors.Color{
							// einstein on the beach kind of gradient
							colors.Black,
							colors.Hex("#737C93"),
							colors.Hex("#7C8099"),
							colors.Hex("#72727A"),
							colors.Hex("#9E9796"),
							colors.Hex("#FCFDF7"),
							colors.Hex("#FBEBE1"),
							colors.Hex("#999596"),
							colors.Hex("#5E6A6A"),
							colors.Black,
						},
					},
				},
			),
			StdDev: 0.003,
		},
		),
	)

	SwivelLines = scenes.CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.NewDynamicObject(
				objects.Parallelogram(
					geometry.Point{X: -1, Y: -1, Z: -2},
					geometry.Point{X: -1, Y: 1, Z: -2},
					geometry.Point{X: 1, Y: -1, Z: -2},
					colors.OpaqueDynamicTexture(
						colors.DynamicFromAnimatedTexture(
							colors.GetAniTextureFromSampler(
								sampler.Sigmoid{
									Sampler: sampler.Wiggle{
										Sampler: sampler.Rotated{
											Sampler: sampler.SineWave{
												Factor: 150,
											},
											Angle: math.Pi / 2,
										},
										NWiggles: 4,
										Angle:    math.Pi / 50,
									},
									Ratio: 5,
								},
								colors.SimpleGradient{Start: colors.White, End: colors.Black},
							),
						),
					),
				),
			),
		},
		Background: scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.GetAniTextureFromSampler(
					sampler.Sigmoid{
						Sampler: sampler.Wiggle{
							Sampler: sampler.SineWave{
								Factor: 400,
							},
							NWiggles: 4,
							Angle:    math.Pi / 50,
						},
						Ratio: 5,
					},
					colors.SimpleGradient{Start: colors.White, End: colors.Black},
				),
			),
		),
	}

	CharMap = scenes.BackgroundScene(
		scenes.BackgroundFromTexture(colors.DynamicFromAnimatedTexture(
			colors.NewDynamicSubtexturer(
				colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
				100,
				sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 10},
			),
		),
		),
	)

	MinecraftCube = scenes.DummyTextureSpinningCube(
		colors.OpaqueDynamicTexture(colors.DynamicFromAnimatedTexture(
			colors.NewDynamicSubtexturer(
				colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
				8,
				sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 5},
			),
		)),
		scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.NewDynamicSubtexturer(
					colors.GetSpecialMapper(colors.White, colors.Black, 0.2),
					32,
					sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 5},
				),
			),
		),
	)

	RoundedSquare = scenes.BackgroundScene(
		scenes.BackgroundFromTexture(
			colors.StaticTexture(
				colors.RoundedSquare{
					On:        colors.White,
					Off:       colors.Black,
					HalfWidth: 0.9,
					Radius:    0.1,
				},
			),
		),
	)

	MulticubeContracting = scenes.MulticubeDance(
		colors.SimpleGradient{Start: colors.Black, End: colors.Red},
		colors.SimpleGradient{Start: colors.Black, End: colors.Green},
		colors.SimpleGradient{Start: colors.Black, End: colors.Blue},
		scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.GetAniTextureFromSampler(
					sampler.SineWaveAnimation{
						XYRatio:      0.1,
						SigmoidRatio: 2,
						SinCycles:    3,
					},
					colors.Subsample(colors.Grayscale, 0.4, 0.6),
				),
			),
		),
	)

	SpinningMulticube = scenes.SpinningMulticube(
		scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.GetAniTextureFromSampler(
					sampler.SineWaveAnimation{
						XYRatio:      0.1,
						SigmoidRatio: 2,
						SinCycles:    3,
					},
					colors.Subsample(colors.Grayscale, 0.4, 0.6),
				),
			),
		),
	)
	Checkckerboard   = scenes.CheckerboardSquare(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Blue})))
	SpinningTriangle = scenes.SingleSpinningTriangle(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Blue})))
	SpinningHolyCube = scenes.SpinningIndividualMulticubeWithHoles(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Blue})))
	HeightMap        = scenes.HeightMap(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Black})))

	SpinningTriangleWithHole = scenes.CheckerboardSquareWithRoundHole(
		scenes.BackgroundFromTexture(
			colors.DynamicFromAnimatedTexture(
				colors.GetAniTextureFromSampler(
					sampler.SineWaveAnimation{
						XYRatio:      0.1,
						SigmoidRatio: 2,
						SinCycles:    3,
					},
					colors.Subsample(colors.Grayscale, 0.4, 0.6),
				),
			),
		),
	)

	Noise                      = scenes.NoiseTest()
	SquaresAlongPath           = scenes.SquaresAlongPath(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Black})))
	SquaresAlongPathWithCamera = scenes.CameraThroughSquaresAlongPath(scenes.BackgroundFromTexture(colors.StaticTexture(colors.Uniform{Color: colors.Black})))
	// Perlin = scenes.NewPerlinNoise(color.Grayscale)
)
