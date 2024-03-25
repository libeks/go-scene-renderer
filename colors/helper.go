package colors

import "math"

type AnimatedTexture interface {
	GetFrameColor(x, y, f float64) Color
}

type dynamicTextureHelper struct {
	ani AnimatedTexture
}

func (d dynamicTextureHelper) GetFrame(t float64) Texture {
	return dynamicTextureFrameHelper{
		ani: d.ani,
		t:   t,
	}
}

type dynamicTextureFrameHelper struct {
	t   float64
	ani AnimatedTexture
}

func (f dynamicTextureFrameHelper) GetTextureColor(x, y float64) Color {
	return f.ani.GetFrameColor(x, y, f.t)
}

func DynamicFromAnimatedTexture(ani AnimatedTexture) DynamicTexture {
	return dynamicTextureHelper{
		ani: ani,
	}
}

type Sampler interface {
	GetFrameValue(x, y, t float64) float64
}

type DynamicSubtexturer struct {
	Subtexture   AnimatedTexture
	N            int // number of squares to tile
	PointSampler Sampler
}

// returns x/d
func bucketRemainder(x, d float64) (float64, float64) {
	return float64(int(x/d)) * d, math.Mod(x, d) * 1 / d
}

func (s DynamicSubtexturer) GetFrameColor(x, y, t float64) Color {
	d := 1 / float64(s.N)
	xMeta, xValue := bucketRemainder(x, d)
	yMeta, yValue := bucketRemainder(y, d)
	tHere := s.PointSampler.GetFrameValue(xMeta, yMeta, t)
	return s.Subtexture.GetFrameColor(xValue, yValue, tHere)
}

type TextureValueMapping struct {
	Above float64
	Texture
}

// StaticMapper displays the static Texture in the list, the first one whose Above value is below t
type StaticMapper struct {
	Mapping []TextureValueMapping // ordered in decreasing order of Above
}

func (m StaticMapper) GetFrameColor(x, y, t float64) Color {
	for _, mapping := range m.Mapping {
		if t >= mapping.Above {
			return mapping.Texture.GetTextureColor(x, y)
		}
	}
	return Red // shouldn't ever happen if the last Mapping starts at 0.0
}

func GetSpecialMapper(on, off Color, thickness float64) StaticMapper {
	return StaticMapper{
		Mapping: []TextureValueMapping{
			{0.9, Square{on, off, 1.0}},
			{0.8, Square{on, off, max(0.7, 2*thickness)}},
			{0.7, Cross{on, off, thickness}},
			{0.5, HorizontalLine{on, off, thickness}},
			{0.4, VerticalLine{on, off, thickness}},
			{0.1, Circle{on, off, thickness}},
			{0.0, Uniform{off}},
		},
	}
}
