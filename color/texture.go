package color

type Texture interface {
	// a,b range from (0,1), when used for triangles, only the lower triangluar values will be called
	// TODO: Texture is surprisingly similar to Frame, maybe there's some generalizations?
	GetTextureColor(b, c float64) Color
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

type RotatedTexture struct {
	t Texture
}

func (t RotatedTexture) GetTextureColor(b, c float64) Color {
	return t.t.GetTextureColor(1-b, 1-c)
}

func RotateTexture180(texture Texture) Texture {
	return RotatedTexture{texture}
}
