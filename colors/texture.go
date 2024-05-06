package colors

import (
	"github.com/libeks/go-scene-renderer/sampler"
)

type Texture interface {
	// a,b range from (0,1), when used for triangles, only the lower triangluar values will be called
	GetTextureColor(b, c float64) Color
}

type TransparentTexture interface {
	GetTextureColor(b, c float64) *Color
}

type DynamicTexture interface {
	GetFrame(t float64) Texture
}

type DynamicTransparentTexture interface {
	GetFrame(t float64) TransparentTexture
}

// a helper for when a static texture is needed as a dynamic texture
type staticTexture struct {
	t Texture
}

type transparentTexture struct {
	texture      Texture
	transparency Transparency
}

func (t transparentTexture) GetTextureColor(b, c float64) *Color {
	if !t.transparency.GetAlpha(b, c) {
		return nil
	}
	color := t.texture.GetTextureColor(b, c)
	return &color
}

type dynamicTransparentTexture struct {
	texture      DynamicTexture
	transparency DynamicTransparency
}

func (t dynamicTransparentTexture) GetFrame(tt float64) TransparentTexture {
	return transparentTexture{
		texture:      t.texture.GetFrame(tt),
		transparency: t.transparency.GetFrame(tt),
	}
}

func GetDynamicTransparentTexture(texture DynamicTexture, transparency DynamicTransparency) dynamicTransparentTexture {
	return dynamicTransparentTexture{
		texture, transparency,
	}
}

type opaqueTexture struct {
	texture Texture
}

func (t opaqueTexture) GetTextureColor(b, c float64) *Color {
	color := t.texture.GetTextureColor(b, c)
	return &color
}

func OpaqueTexture(t Texture) TransparentTexture {
	return opaqueTexture{t}
}

type dynamicOpaqueTexture struct {
	texture DynamicTexture
}

func (t dynamicOpaqueTexture) GetFrame(tt float64) TransparentTexture {
	return opaqueTexture{t.texture.GetFrame(tt)}
}

func OpaqueDynamicTexture(t DynamicTexture) DynamicTransparentTexture {
	return dynamicOpaqueTexture{t}
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

func GenerateDynamicFromAnimatedTexture(f func(b float64, c float64, t float64) Color) genericDynamicSampler {
	return genericDynamicSampler{f}
}

type genericDynamicSampler struct {
	f func(b, c, t float64) Color
}

func (g genericDynamicSampler) GetFrame(t float64) Texture {
	return genericStaticSampler{g, t}
}

type genericStaticSampler struct {
	genericDynamicSampler
	t float64
}

func (g genericStaticSampler) GetTextureColor(b, c float64) Color {
	return g.genericDynamicSampler.f(b, c, g.t)
}

func RBGSamplerTexture(R, G, B sampler.Sampler) DynamicTexture {
	return GenerateDynamicFromAnimatedTexture(func(b, c, t float64) Color {
		return Color{
			R.GetFrameValue(b, c, t),
			G.GetFrameValue(b, c, t),
			B.GetFrameValue(b, c, t),
		}
	})
}

func RBGSamplerDynamicTexture(R, G, B sampler.DynamicSampler) rgbDynamicSampler {
	return rgbDynamicSampler{
		R,
		G,
		B,
	}
}

type rgbDynamicSampler struct {
	R sampler.DynamicSampler
	G sampler.DynamicSampler
	B sampler.DynamicSampler
}

func (s rgbDynamicSampler) GetFrame(t float64) Texture {
	return RGBFrame{
		R: s.R.GetFrame(t),
		G: s.G.GetFrame(t),
		B: s.B.GetFrame(t),
	}
}

type RGBFrame struct {
	R sampler.StaticSampler
	G sampler.StaticSampler
	B sampler.StaticSampler
}

func (f RGBFrame) GetTextureColor(b, c float64) Color {
	return Color{
		R: f.R.GetValue(b, c),
		G: f.G.GetValue(b, c),
		B: f.B.GetValue(b, c),
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
	t TransparentTexture
}

func (t rotatedTexture) GetTextureColor(b, c float64) *Color {
	return t.t.GetTextureColor(1-b, 1-c)
}

func RotateTexture180(texture TransparentTexture) TransparentTexture {
	return rotatedTexture{texture}
}

type rotatedDynamicTexture struct {
	t DynamicTransparentTexture
}

func (r rotatedDynamicTexture) GetFrame(t float64) TransparentTexture {
	return RotateTexture180(r.t.GetFrame(t))
}

func RotateDynamicTexture180(texture DynamicTransparentTexture) DynamicTransparentTexture {
	return rotatedDynamicTexture{texture}
}

type VerticalGradientTexture struct {
	Gradient
}

func (g VerticalGradientTexture) GetTextureColor(b, c float64) Color {
	return g.Gradient.Interpolate(c)
}
