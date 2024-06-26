package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/sampler"
	"github.com/libeks/go-scene-renderer/textures"
)

func DummySpinningCubes(background DynamicBackground) DynamicScene {
	initialCube := UnitRGBCube()
	diagonalCube := initialCube

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, math.Sqrt(3)/2, -3)),
					geometry.RotateMatrixY(t*maths.Rotation),
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, -math.Sqrt(3)/2, -3)),
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

	spacing := 2.0

	column := objects.CombineDynamicObjects(
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.V3(0, -2*spacing, 0))),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.V3(0, -spacing, 0))),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, 0))),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.V3(0, spacing, 0))),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 2*spacing, 0))),
	)

	slice := objects.CombineDynamicObjects(
		column.WithTransform(geometry.TranslationMatrix(geometry.V3(-2*spacing, 0, 0))),
		column.WithTransform(geometry.TranslationMatrix(geometry.V3(-spacing, 0, 0))),
		column.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, 0))),
		column.WithTransform(geometry.TranslationMatrix(geometry.V3(spacing, 0, 0))),
		column.WithTransform(geometry.TranslationMatrix(geometry.V3(2*spacing, 0, 0))),
	)

	multiCube := objects.CombineDynamicObjects(
		slice.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, -2*spacing))),
		slice.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, -spacing))),
		slice.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, 0))),
		slice.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, spacing))),
		slice.WithTransform(geometry.TranslationMatrix(geometry.V3(0, 0, 2*spacing))),
	)

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			multiCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -10)),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
				)
			},
			),
		},
		Background: background,
	}
}

func NoiseTest() DynamicScene {
	texture := textures.StaticTexture(textures.NewPerlinNoiseTexture(colors.Grayscale))
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.Parallelogram(geometry.Pt(0, 0, -5), geometry.Pt(2, 0, -5), geometry.Pt(0, 2, -5), textures.OpaqueDynamicTexture(texture)),
		},
		Background: BackgroundFromTexture(texture),
	}
}

func SpinningIndividualMulticube(background DynamicBackground) DynamicScene {
	texture := textures.OpaqueDynamicTexture(textures.StaticTexture(textures.NewPerlinNoiseTexture(colors.Grayscale)))
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
		Objects: []objects.DynamicObjectInt{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -5)),       // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                // rotation around common center
					geometry.TranslationMatrix(geometry.V3(-spacing, 0, 0)), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),   // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -5)),     // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),              // rotation around common center
					geometry.TranslationMatrix(geometry.V3(0, 0, 0)),      // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)), // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -5)),      // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),               // rotation around common center
					geometry.TranslationMatrix(geometry.V3(spacing, 0, 0)), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),  // rotation around own axis
				)
			}),
		},
		Background: background,
	}
}

func SpinningIndividualMulticubeWithHoles(background DynamicBackground) DynamicScene {
	texture := textures.StaticTexture(textures.NewPerlinNoiseTexture(colors.Grayscale))
	initialCube := UnitTextureCubeWithTransparency(
		texture,
		texture,
		texture,
		texture,
		texture,
		texture,
		textures.DynamicFromAnimatedTransparency(
			textures.CircleCutout{Radius: 0.8},
		),
	)
	diagonalCube := initialCube.WithTransform(
		geometry.MatrixProduct(
			geometry.RotateMatrixX(-0.615),
			geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
		)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	// spacing := 2.0

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -1.5)), // position within the scene
					geometry.ScaleMatrix(1),                             // scale around own center
					geometry.RotateMatrixY(-2*t*maths.Rotation),         // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -1.5)), // position within the scene
					geometry.ScaleMatrix(0.7),                           // scale around own center
					geometry.RotateMatrixZ(-2*t*maths.Rotation),         // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -1.5)), // position within the scene
					geometry.ScaleMatrix(0.49),                          // scale around own center
					geometry.RotateMatrixX(-2*t*maths.Rotation),         // rotation around own axis
				)
			}),
		},
		Background: background,
	}
}

func MulticubeDance(g1, g2, g3 colors.Gradient, background DynamicBackground) DynamicScene {
	spacing := 2.0
	cubes := 19
	ratio := 1 / float64(cubes)
	objs := make([]objects.DynamicObjectInt, 0, cubes*cubes*cubes)
	timeRatio := 1.5
	for d1 := range cubes {
		for d2 := range cubes {
			for d3 := range cubes {
				x1 := float64(d1)
				x2 := float64(d2)
				x3 := float64(d3)

				cube := UnitGradientCube(
					g1.Interpolate((x1+0)*ratio).Add(g2.Interpolate((x2+0)*ratio)).Add(g3.Interpolate((x3+0)*ratio)),
					g1.Interpolate((x1+1)*ratio).Add(g2.Interpolate((x2+0)*ratio)).Add(g3.Interpolate((x3+0)*ratio)),
					g1.Interpolate((x1+1)*ratio).Add(g2.Interpolate((x2+1)*ratio)).Add(g3.Interpolate((x3+0)*ratio)),
					g1.Interpolate((x1+0)*ratio).Add(g2.Interpolate((x2+1)*ratio)).Add(g3.Interpolate((x3+0)*ratio)),

					g1.Interpolate((x1+0)*ratio).Add(g2.Interpolate((x2+0)*ratio)).Add(g3.Interpolate((x3+1)*ratio)),
					g1.Interpolate((x1+1)*ratio).Add(g2.Interpolate((x2+0)*ratio)).Add(g3.Interpolate((x3+1)*ratio)),
					g1.Interpolate((x1+1)*ratio).Add(g2.Interpolate((x2+1)*ratio)).Add(g3.Interpolate((x3+1)*ratio)),
					g1.Interpolate((x1+0)*ratio).Add(g2.Interpolate((x2+1)*ratio)).Add(g3.Interpolate((x3+1)*ratio)),
				)

				objs = append(objs, cube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(0, 0, -2*float64(cubes))), // position within the scene
						geometry.RotateMatrixY(t*maths.Rotation),
						geometry.TranslationMatrix(geometry.V3(
							(float64(d1-cubes/2))*(spacing-min(t*timeRatio, 1)),
							(float64(d2-cubes/2))*(spacing-min(t*timeRatio, 1)),
							(float64(d3-cubes/2))*(spacing-min(t*timeRatio, 1)),
						)), // position within the scene

					)
				}))
			}
		}
	}

	return CombinedDynamicScene{
		Objects:    objs,
		Background: background,
	}
}

func DummySpinningCube(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			UnitRGBCube().WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -2)),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
					geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
				)
			}),
		},
		Background: background,
	}
}

func DummyTextureSpinningCube(t textures.DynamicTransparentTexture, background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			UnitTextureCube(t, t, t, t, t, t).WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -1.5)),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
					geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
				)
			}),
		},
		Background: background,
	}
}

func SingleSpinningTriangle(background DynamicBackground) DynamicScene {
	transform := func(t float64) geometry.HomogeneusMatrix {
		return geometry.MatrixProduct(
			geometry.TranslationMatrix(geometry.V3(1.5, 0, -3)),
			geometry.RotateMatrixZ(t*maths.Rotation),
		)
	}
	transform2 := func(t float64) geometry.HomogeneusMatrix {
		return geometry.MatrixProduct(
			geometry.TranslationMatrix(geometry.V3(-1.5, 0, -3)),
			geometry.RotateMatrixZ(t*maths.Rotation),
		)
	}
	a99, a90, a91 := 0.0, 1.0, 0.0
	a09, a00, a01 := 1.0, 0.5, 1.0
	a19, a10, a11 := 0.0, 1.0, 0.0
	gradient := colors.LinearGradient{Points: []colors.Color{colors.Black, colors.Grayscale.Interpolate(0.75), colors.Red}}
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(1, 1, 0),
						geometry.Pt(1, 0, 0),
						geometry.Pt(0, 1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a11, B: a10, C: a01, D: a00,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 0, 0),
						geometry.Pt(1, 0, 0),
						geometry.Pt(0, 1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a00, B: a10, C: a01, D: a11,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(1, -1, 0),
						geometry.Pt(1, 0, 0),
						geometry.Pt(0, -1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a19, B: a10, C: a09, D: a00,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 0, 0),
						geometry.Pt(1, 0, 0),
						geometry.Pt(0, -1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a00, B: a10, C: a09, D: a19,
							},
						),
					),
				),
			).WithDynamicTransform(transform),

			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(-1, 1, 0),
						geometry.Pt(-1, 0, 0),
						geometry.Pt(0, 1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a91, B: a90, C: a01, D: a00,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 0, 0),
						geometry.Pt(-1, 0, 0),
						geometry.Pt(0, 1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a00, B: a90, C: a01, D: a91,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(-1, -1, 0),
						geometry.Pt(-1, 0, 0),
						geometry.Pt(0, -1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a99, B: a90, C: a09, D: a00,
							},
						),
					),
				),
			).WithDynamicTransform(transform),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 0, 0),
						geometry.Pt(-1, 0, 0),
						geometry.Pt(0, -1, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a00, B: a90, C: a09, D: a99,
							},
						),
					),
				),
			).WithDynamicTransform(transform),

			// triangles towards the middle
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(-1, 0, 0),
						geometry.Pt(-1, -1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a90, B: a99, C: a00, D: a09,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, -1, 0),
						geometry.Pt(-1, -1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a09, B: a99, C: a00, D: a90,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(1, 0, 0),
						geometry.Pt(1, -1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a10, B: a19, C: a00, D: a09,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, -1, 0),
						geometry.Pt(1, -1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a09, B: a19, C: a00, D: a10,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),

			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(1, 0, 0),
						geometry.Pt(1, 1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a10, B: a11, C: a00, D: a01,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 1, 0),
						geometry.Pt(1, 1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a01, B: a11, C: a00, D: a10,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(-1, 0, 0),
						geometry.Pt(-1, 1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a90, B: a91, C: a00, D: a01,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
			objects.DynamicObjectFromBasics(
				objects.DynamicBasicObject(
					objects.Tri(
						geometry.Pt(0, 1, 0),
						geometry.Pt(-1, 1, 0),
						geometry.Pt(0, 0, 0),
					),
					textures.OpaqueDynamicTexture(
						textures.StaticTexture(
							textures.TriangleGradientInterpolationTexture{
								Gradient: gradient,

								A: a01, B: a91, C: a00, D: a90,
							},
						),
					),
				),
			).WithDynamicTransform(transform2),
		},
		Background: background,
	}
}

func CheckerboardSquare(background DynamicBackground) DynamicScene {
	texture := textures.OpaqueDynamicTexture(textures.StaticTexture(textures.Checkerboard{Squares: 16}))
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.Parallelogram(
				geometry.Pt(0, 0, 0),
				geometry.Pt(2, 0, 0),
				geometry.Pt(-2, 2, 0),
				texture).WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
						geometry.RotateMatrixX(t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
	}
}

func CheckerboardSquareWithRoundHole(background DynamicBackground) DynamicScene {
	texture := textures.StaticTexture(textures.Checkerboard{Squares: 16})
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.Parallelogram(
				geometry.Pt(0, 0, 0),
				geometry.Pt(2, 0, 0),
				geometry.Pt(0, 2, 0),
				textures.GetDynamicTransparentTexture(
					texture,
					textures.DynamicFromAnimatedTransparency(
						textures.CircleCutout{Radius: 0.8},
					),
				),
			).WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
						geometry.RotateMatrixY(t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
	}
}

func HeightMap(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.NewDynamicObject(
				objects.HeightMap{
					Height: sampler.UnitCircleClamper{
						DynamicSampler: sampler.RotatingSampler{
							DynamicSampler: sampler.DynamicFromAnimated(sampler.Sigmoid{Sampler: sampler.NewPerlinNoise(), Ratio: 5}),
							Rotations:      1,
							Radius:         0.25,
							OffsetX:        0.5,
							OffsetY:        0.5,
							OffsetT:        10,
						},
						MaxRadius: 0.95,
						Decay:     9,
					},
					Gradient: colors.LinearGradient{Points: []colors.Color{colors.Black, colors.Grayscale.Interpolate(0.75), colors.White}},
					N:        100,
				},
			).WithDynamicTransform(

				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.RotateMatrixX(0.2),
						geometry.TranslationMatrix(geometry.V3(0, -0.8, -1.5)),
						geometry.ScaleMatrix(1),
						// geometry.RotateMatrixY(-t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
	}
}

func SquaresAlongPath(background DynamicBackground) DynamicScene {
	path := geometry.BezierPath{
		Points: []geometry.Point{
			{X: 0, Y: 0, Z: 0},
			{X: 3, Y: 0, Z: 0},
			{X: 3, Y: 3, Z: 0},
		},
	}
	checkerTexture := textures.DynamicTexture(textures.StaticTexture(textures.Checkerboard{Squares: 8}))
	texture := textures.GetDynamicTransparentTexture(
		checkerTexture,
		textures.DynamicFromAnimatedTransparency(
			textures.CircleCutout{Radius: 0.8},
		),
	)

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.RectanglesAlongPath(path, 100, 1, texture).
				WithDynamicTransform(
					func(t float64) geometry.HomogeneusMatrix {
						return geometry.MatrixProduct(
							geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
							geometry.RotateMatrixY(-t*maths.Rotation),
						)
					},
				),
		},
		Background: background,
	}
}

func CameraThroughSquaresAlongPath(background DynamicBackground) DynamicScene {
	path := geometry.BezierPath{
		Points: []geometry.Point{
			{X: 0, Y: 0, Z: 0},
			{X: 1, Y: 10, Z: 0},
			{X: -1, Y: 0, Z: 8},
			{X: 0, Y: 10, Z: 16},
			{X: 2, Y: 10, Z: 16},
			// {X: 0, Y: 13, Z: 16},
		},
	}
	outerSphereTexture := textures.GetDynamicTransparentTexture(
		textures.StaticTexture(
			textures.Uniform(colors.Black),
		),
		textures.StaticTransparency(textures.MiddleBand{Min: 0.0, Max: 0.7}),
	)
	middleSphereTexture := textures.GetDynamicTransparentTexture(
		textures.StaticTexture(
			textures.Uniform(colors.Hex("#ff4500")),
		),
		textures.StaticTransparency(textures.MiddleBand{Min: 0.0, Max: 0.6}),
	)
	innerSphereTexture := textures.OpaqueDynamicTexture(
		textures.StaticTexture(
			textures.Uniform(colors.White),
		),
	)
	spheres := objects.DynamicObjectFromBasics(
		objects.DynamicSphere(objects.UnitSphere(), outerSphereTexture),
		*objects.DynamicSphere(objects.UnitSphere(), middleSphereTexture).ApplyMatrix(geometry.ScaleMatrix(0.98)),
		*objects.DynamicSphere(objects.UnitSphere(), innerSphereTexture).ApplyMatrix(geometry.ScaleMatrix(0.60)),
	).WithTransform(geometry.ScaleMatrix(0.6))
	cameraPath := geometry.SamplePath(path, 0, 0.8)
	spherePath := geometry.SamplePath(path, 0.2, 1)
	gateTexture := textures.StaticTexture(textures.BinarySamplerWithColors{
		StaticSampler: sampler.ConcentricCircles(0.01),
		On:            colors.Black,
		Off:           colors.White,
	})
	texture := textures.GetDynamicTransparentTexture(
		gateTexture,
		textures.DynamicFromAnimatedTransparency(
			textures.CircleCutout{Radius: 0.12},
		),
	)

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.RectanglesAlongPath(path, 50, 5, texture),
			spheres.WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					sphereLocation := spherePath.GetDirection(t).Origin
					cameraLocation := cameraPath.GetDirection(t).Origin
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D(sphereLocation)),
						geometry.PointTowards(geometry.Point(sphereLocation.Subtract(cameraLocation))),
						geometry.RotatePitch(-math.Pi/2),
					)
				},
			),
		},
		Background: background,
		CameraPath: cameraPath,
	}
}

func CameraWithAxisTriangles(background DynamicBackground) DynamicScene {
	path := geometry.BezierPath{
		Points: []geometry.Point{
			{X: 0, Y: 0, Z: 0},
			{X: 1, Y: 10, Z: 0},
			{X: -1, Y: 0, Z: 8},
			{X: 0, Y: 10, Z: 16},
		},
	}

	cameraPath := geometry.SamplePath(path, 0, 0.8)
	spherePath := geometry.SamplePath(path, 0.2, 1)

	actualCameraPath := geometry.BezierPath{
		Points: []geometry.Point{
			{X: 0, Y: 0, Z: -4},
			{X: 0, Y: 0, Z: -3},
		},
	}

	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.AxisAlignedPointer().WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D(spherePath.GetDirection(t).Origin)),
						geometry.PointTowards(cameraPath.GetDirection(t).Origin),
						geometry.RotateYaw(math.Pi/2),
					)
				},
			),
			objects.AxisAlignedPointer().WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D(cameraPath.GetDirection(t).Origin)),
						// geometry.PointTowards(cameraPath.GetDirection(t).Origin),
						// geometry.RotateYaw(math.Pi/2),
					)
				},
			),
		},
		Background: background,
		CameraPath: actualCameraPath,
	}
}

func ThreeSpheres(background DynamicBackground) DynamicScene {
	checkerTexture := textures.DynamicTexture(textures.StaticTexture(textures.Checkerboard{Squares: 16}))
	whiteTexture := textures.StaticTexture(textures.Uniform(colors.White))
	texture := textures.GetDynamicTransparentTexture(
		checkerTexture,
		textures.DynamicFromAnimatedTransparency(
			textures.CircleCutout{Radius: 0.6},
		),
	)
	transparentCheckerTexture := textures.GetDynamicTransparentTexture(
		whiteTexture,
		textures.StaticTransparency(textures.Checkerboard{Squares: 16}),
	)
	redTexture := textures.OpaqueDynamicTexture(
		textures.StaticTexture(
			textures.VerticalGradientTexture{
				Gradient: colors.LinearGradient{
					Points: []colors.Color{colors.Black, colors.Red, colors.Red, colors.Black},
				},
			},
		),
	)

	dynamicSphere := objects.DynamicObjectFromBasics(objects.DynamicSphere(objects.UnitSphere(), texture))
	checkerSphere := objects.DynamicObjectFromBasics(objects.DynamicSphere(objects.UnitSphere(), transparentCheckerTexture))
	redSphere := objects.DynamicObjectFromBasics(objects.DynamicSphere(objects.UnitSphere(), redTexture))
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			dynamicSphere.WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(-3, 0, -5)),
						geometry.RotateMatrixY(-2*t*maths.Rotation),
					)
				},
			),
			redSphere.WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
						geometry.RotateMatrixZ(-3*t*maths.Rotation),
					)
				},
			),
			checkerSphere.WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.V3(3, 0, -5)),
						geometry.RotateMatrixX(4*t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
		// CameraPath: path,
	}
}

func NineSpheres(background DynamicBackground) DynamicScene {
	redTexture := textures.OpaqueDynamicTexture(textures.StaticTexture(textures.Uniform(colors.Red)))
	yellowSquare := objects.Parallelogram(geometry.Pt(-1, -1, 0), geometry.Pt(1, -1, 0), geometry.Pt(-1, 1, 0),
		textures.OpaqueDynamicTexture(textures.StaticTexture(textures.Uniform(colors.Yellow))),
	)

	redSphere := objects.DynamicObjectFromBasics(objects.DynamicSphere(objects.UnitSphere(), redTexture))
	objects := []objects.DynamicObjectInt{}
	objects = append(objects,
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, 3, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 3, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, 3, -5)),
				)
			},
		),

		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, 0, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, 0, -5)),
				)
			},
		),

		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, -3, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, -3, -5)),
				)
			},
		),
		redSphere.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, -3, -5)),
				)
			},
		),
	)
	objects = append(objects,
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, 3, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 3, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, 3, -5)),
				)
			},
		),

		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, 0, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, 0, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, 0, -5)),
				)
			},
		),

		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(3, -3, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(0, -3, -5)),
				)
			},
		),
		yellowSquare.WithDynamicTransform(
			func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.V3(-3, -3, -5)),
				)
			},
		),
	)
	return CombinedDynamicScene{
		Objects:    objects,
		Background: background,
		// CameraPath: path,
	}
}

func OneBigSphere(background DynamicBackground) DynamicScene {
	redTexture := textures.OpaqueDynamicTexture(textures.StaticTexture(textures.Uniform(colors.Red)))
	location := func(t float64) geometry.HomogeneusMatrix {
		return geometry.MatrixProduct(
			geometry.RotateRoll(t*maths.Rotation),
			geometry.TranslationMatrix(geometry.V3(3, 0, -5)),
		)
	}

	redSphere := objects.DynamicObjectFromBasics(objects.DynamicSphere(objects.UnitSphere(), redTexture))
	objects := []objects.DynamicObjectInt{}
	objects = append(objects,

		redSphere.WithDynamicTransform(
			location,
		),
	)
	return CombinedDynamicScene{
		Objects:    objects,
		Background: background,
	}
}

func HeightMapCross(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObjectInt{
			objects.NewDynamicObject(
				objects.HeightMap{
					Height:   sampler.RotatingCross(0.1),
					Gradient: colors.LinearGradient{Points: []colors.Color{colors.Black, colors.Grayscale.Interpolate(0.75), colors.White}},
					N:        50,
				},
			).WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.RotateMatrixX(0.2),
						geometry.TranslationMatrix(geometry.V3(-0.5, -1.3, -1.5)),
						geometry.ScaleMatrix(1),
						// geometry.RotateMatrixY(-t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
	}
}
