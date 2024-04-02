package colors

import (
// "fmt"
)

type Texture interface {
	// a,b range from (0,1), when used for triangles, only the lower triangluar values will be called
	GetTextureColor(b, c float64) Color
}

type DynamicTexture interface {
	GetFrame(t float64) Texture
}

// a helper for when a static texture is needed as a dynamic texture
type staticTexture struct {
	t Texture
}

func (t staticTexture) GetFrame(f float64) Texture {
	return t.t
}

func StaticTexture(t Texture) DynamicTexture {
	return staticTexture{t}
}

func TriangleGradientTexture(A, B, C Color) Texture {
	return triangleGradientTexture{
		A, B, C,
	}
}

// returns a color gradient for the lower triangle of a unit square
type triangleGradientTexture struct {
	ColorA Color
	ColorB Color
	ColorC Color
}

// given coordinates from the A point towards B and C (each in the range of (0,1))
// return what color it should be
func (t triangleGradientTexture) GetTextureColor(b, c float64) Color {
	if c == 1 {
		return t.ColorB
	}
	abGradient := SimpleGradient{t.ColorA, t.ColorB}
	abColor := abGradient.Interpolate(b / (1 - c))
	triangleGradient := SimpleGradient{abColor, t.ColorC}
	cColor := triangleGradient.Interpolate(c)
	return cColor
}

type TriangleGradientInterpolationTexture struct {
	Gradient
	A float64
	B float64
	C float64
	D float64
}

func (g TriangleGradientInterpolationTexture) GetTextureColor(b, c float64) Color {
	if b == 1 {
		return g.Gradient.Interpolate(g.B)
	}

	AtoB := (b*(g.B) + (1-b)*g.A)
	CtoD := (b*(g.D) + (1-b)*g.C)
	t := (c*(CtoD) + (1-c)*AtoB)

	cColor := g.Gradient.Interpolate(t)
	return cColor
}

func SquareGradientTexture(A, B, C, D Color) Texture {
	return squareGradientTexture{
		triangleGradientTexture{
			A, B, C,
		},
		triangleGradientTexture{
			D, C, B,
		},
	}
}

// returns a color gradient for the whole unit square, based on colors of the four corners, laid out
// as A,B,D,C (A and D are opposites)
type squareGradientTexture struct {
	lower triangleGradientTexture
	upper triangleGradientTexture
}

func (g squareGradientTexture) GetTextureColor(b, c float64) Color {
	if b+c < 1.0 {
		return g.lower.GetTextureColor(b, c)
	}

	// besure to flip coordinates from the other end
	return g.upper.GetTextureColor(1-b, 1-c)
}

type rotatedTexture struct {
	t Texture
}

func (t rotatedTexture) GetTextureColor(b, c float64) Color {
	return t.t.GetTextureColor(1-b, 1-c)
}

func RotateTexture180(texture Texture) Texture {
	return rotatedTexture{texture}
}

type rotatedDynamicTexture struct {
	t DynamicTexture
}

func (r rotatedDynamicTexture) GetFrame(t float64) Texture {
	return RotateTexture180(r.t.GetFrame(t))
}

func RotateDynamicTexture180(texture DynamicTexture) DynamicTexture {
	return rotatedDynamicTexture{texture}
}
