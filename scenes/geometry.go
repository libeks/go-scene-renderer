package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
)

func DummySpinningCubes(background DynamicBackground) DynamicScene {
	initialCube := UnitRGBCube()
	diagonalCube := initialCube

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, math.Sqrt(3) / 2, -3}),
					geometry.RotateMatrixY(t*maths.Rotation),
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, -math.Sqrt(3) / 2, -3}),
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
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, -2 * spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, -spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, spacing, 0})),
		diagonalCube.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 2 * spacing, 0})),
	)

	slice := objects.CombineDynamicObjects(
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{-2 * spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0})),
		column.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{2 * spacing, 0, 0})),
	)

	multiCube := objects.CombineDynamicObjects(
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2 * spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, -spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, spacing})),
		slice.WithTransform(geometry.TranslationMatrix(geometry.Vector3D{0, 0, 2 * spacing})),
	)

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			multiCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -10}),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
				)
			},
			),
		},
		Background: background,
	}
}

func NoiseTest() DynamicScene {
	texture := colors.StaticTexture(colors.NewPerlinNoiseTexture(colors.Grayscale))
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.Parallelogram(geometry.Point{0, 0, -5}, geometry.Point{2, 0, -5}, geometry.Point{0, 2, -5}, texture),
		},
		Background: BackgroundFromTexture(texture),
	}
}

func SpinningIndividualMulticube(background DynamicBackground) DynamicScene {
	// initialCube := UnitRGBCube()
	// texture := color.SquareGradientTexture(color.White, color.Red, color.Black, color.Blue)
	texture := colors.StaticTexture(colors.NewPerlinNoiseTexture(colors.Grayscale))
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
		Objects: []objects.DynamicObject{
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),       // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                      // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{-spacing, 0, 0}), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),         // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}), // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),  // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),   // rotation around own axis
				)
			}),
			diagonalCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),      // position within the scene
					geometry.RotateMatrixY(t*maths.Rotation),                     // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{spacing, 0, 0}), // position within the group
					geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),        // rotation around own axis
				)
			}),
		},
		Background: background,
	}
}

func MulticubeDance(background DynamicBackground) DynamicScene {
	colorCube := UnitRGBCube()
	// texture := color.SquareGradientTexture(color.White, color.Red, color.Black, color.Blue)
	// texture := colors.StaticTexture(colors.NewPerlinNoiseTexture(colors.Grayscale))
	redGradient := colors.SimpleGradient{colors.Black, colors.Red}
	red03 := redGradient.Interpolate(0) // black, but it makes the math simpler
	red13 := redGradient.Interpolate(0.33333)
	red23 := redGradient.Interpolate(0.66666)
	red33 := redGradient.Interpolate(1) // red, but it makes the math simpler
	greenGradient := colors.SimpleGradient{colors.Black, colors.Green}
	green03 := greenGradient.Interpolate(0) // black, but it makes the math simpler
	green13 := greenGradient.Interpolate(0.33333)
	green23 := greenGradient.Interpolate(0.66666)
	green33 := greenGradient.Interpolate(1) // green, but it makes the math simpler
	blueGradient := colors.SimpleGradient{colors.Black, colors.Green}
	blue03 := blueGradient.Interpolate(0) // black, but it makes the math simpler
	blue13 := blueGradient.Interpolate(0.33333)
	// blue23 := blueGradient.Interpolate(0.66666)
	// blue33 := redGradient.Interpolate(1) // blue, but it makes the math simpler
	c11 := UnitGradientCube(
		red03.Add(blue03).Add(green03),
		red13.Add(blue03).Add(green03),
		red13.Add(blue03).Add(green13),
		red03.Add(blue03).Add(green13),

		red03.Add(blue13).Add(green03),
		red13.Add(blue13).Add(green03),
		red13.Add(blue13).Add(green13),
		red03.Add(blue13).Add(green13),
	)

	c12 := UnitGradientCube(
		red13.Add(blue03).Add(green03),
		red23.Add(blue03).Add(green03),
		red23.Add(blue03).Add(green13),
		red13.Add(blue03).Add(green13),

		red13.Add(blue13).Add(green03),
		red23.Add(blue13).Add(green03),
		red23.Add(blue13).Add(green13),
		red13.Add(blue13).Add(green13),
	)

	c13 := UnitGradientCube(
		red23.Add(blue03).Add(green03),
		red33.Add(blue03).Add(green03),
		red33.Add(blue03).Add(green13),
		red23.Add(blue03).Add(green13),

		red23.Add(blue13).Add(green03),
		red33.Add(blue13).Add(green03),
		red33.Add(blue13).Add(green13),
		red23.Add(blue13).Add(green13),
	)

	c21 := UnitGradientCube(
		red03.Add(blue03).Add(green13),
		red13.Add(blue03).Add(green13),
		red13.Add(blue03).Add(green23),
		red03.Add(blue03).Add(green23),

		red03.Add(blue13).Add(green13),
		red13.Add(blue13).Add(green13),
		red13.Add(blue13).Add(green23),
		red03.Add(blue13).Add(green23),
	)

	c22 := UnitGradientCube(
		red13.Add(blue03).Add(green13),
		red23.Add(blue03).Add(green13),
		red23.Add(blue03).Add(green23),
		red13.Add(blue03).Add(green23),

		red13.Add(blue13).Add(green13),
		red23.Add(blue13).Add(green13),
		red23.Add(blue13).Add(green23),
		red13.Add(blue13).Add(green23),
	)

	c23 := UnitGradientCube(
		red23.Add(blue03).Add(green13),
		red33.Add(blue03).Add(green13),
		red33.Add(blue03).Add(green23),
		red23.Add(blue03).Add(green23),

		red23.Add(blue13).Add(green13),
		red33.Add(blue13).Add(green13),
		red33.Add(blue13).Add(green23),
		red23.Add(blue13).Add(green23),
	)

	c31 := UnitGradientCube(
		red03.Add(blue03).Add(green23),
		red13.Add(blue03).Add(green23),
		red13.Add(blue03).Add(green33),
		red03.Add(blue03).Add(green33),

		red03.Add(blue13).Add(green23),
		red13.Add(blue13).Add(green23),
		red13.Add(blue13).Add(green33),
		red03.Add(blue13).Add(green33),
	)

	c32 := UnitGradientCube(
		red13.Add(blue03).Add(green23),
		red23.Add(blue03).Add(green23),
		red23.Add(blue03).Add(green33),
		red13.Add(blue03).Add(green33),

		red13.Add(blue13).Add(green23),
		red23.Add(blue13).Add(green23),
		red23.Add(blue13).Add(green33),
		red13.Add(blue13).Add(green33),
	)

	c33 := UnitGradientCube(
		red23.Add(blue03).Add(green23),
		red33.Add(blue03).Add(green23),
		red33.Add(blue03).Add(green33),
		red23.Add(blue03).Add(green33),

		red23.Add(blue13).Add(green23),
		red33.Add(blue13).Add(green23),
		red33.Add(blue13).Add(green33),
		red23.Add(blue13).Add(green33),
	)

	// diagonalCube := c1.WithTransform(
	// 	geometry.MatrixProduct(
	// 		geometry.RotateMatrixX(-0.615),
	// 		geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
	// 	)) // cube with lower point at (0,0,0), upper at (0,sqrt(3) ,0)

	spacing := 2.0

	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			c11.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, -spacing + min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{-spacing + min(t*1.3, 1), 0, 0}),  // position within the group

				)
			}),
			c12.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, -spacing + min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),                         // position within the group
				)
			}),
			c13.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, -spacing + min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{spacing - min(t*1.3, 1), 0, 0}),   // position within the group
				)
			}),

			c21.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),                       // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{-spacing + min(t*1.3, 1), 0, 0}), // position within the group

				)
			}),
			c22.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),  // position within the group
				)
			}),
			c23.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),                      // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{spacing - min(t*1.3, 1), 0, 0}), // position within the group
				)
			}),
			c31.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, spacing - min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{-spacing + min(t*1.3, 1), 0, 0}), // position within the group

				)
			}),
			c32.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, spacing - min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, 0}),                        // position within the group
				)
			}),
			c33.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, spacing - min(t*1.3, 1), -5}), // position within the scene
					geometry.TranslationMatrix(geometry.Vector3D{spacing - min(t*1.3, 1), 0, 0}),  // position within the group
				)
			}),
			colorCube.WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}), // position within the scene
					// geometry.RotateMatrixY(t*maths.Rotation),                     // rotation around common center
					geometry.TranslationMatrix(geometry.Vector3D{2 * spacing, 0, 0}), // position within the group
					// geometry.RotateMatrixY(math.Sin(-2*t*maths.Rotation)),        // rotation around own axis
				)
			}),
		},
		Background: background,
	}
}

func DummySpinningCube(background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			UnitRGBCube().WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -2}),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
					geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
				)
			}),
		},
		Background: background,
	}
}

func DummyTextureSpinningCube(t colors.DynamicTexture, background DynamicBackground) DynamicScene {
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			UnitTextureCube(t, t, t, t, t, t).WithDynamicTransform(func(t float64) geometry.HomogeneusMatrix {
				return geometry.MatrixProduct(
					geometry.TranslationMatrix(geometry.Vector3D{0, 0, -1.5}),
					geometry.RotateMatrixY(maths.SigmoidSlowFastSlow(t)*maths.Rotation),
					geometry.RotateMatrixX(-0.615),    // arcsin of 1/sqrt(3) (angle between short and long diagonals in a cube)
					geometry.RotateMatrixZ(math.Pi/4), // arcsin(1/sqrt(2)), angle between edge and short diagonal
				)
			}),
		},
		Background: background,
	}
}

func CheckerboardSquare(background DynamicBackground) DynamicScene {
	texture := colors.StaticTexture(colors.Checkerboard{16})
	return CombinedDynamicScene{
		Objects: []objects.DynamicObject{
			objects.Parallelogram(
				geometry.Point{0, 0, 0},
				geometry.Point{2, 0, 0},
				geometry.Point{-2, 2, 0},
				texture).WithDynamicTransform(
				func(t float64) geometry.HomogeneusMatrix {
					return geometry.MatrixProduct(
						geometry.TranslationMatrix(geometry.Vector3D{0, 0, -5}),
						geometry.RotateMatrixX(t*maths.Rotation),
					)
				},
			),
		},
		Background: background,
	}
}
