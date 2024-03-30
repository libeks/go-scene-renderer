package colors

import (
	"math"

	"github.com/libeks/go-scene-renderer/maths"
)

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

// gets values from rotating around the OffsetX,Y by distance of Radius, parameterized by t,
// a total of Rotation rotations (negative to rotate in opposite direction)
type RotatingSampler struct {
	Sampler
	Rotations float64
	Radius    float64
	OffsetX   float64
	OffsetY   float64
	OffsetT   float64
}

func (s RotatingSampler) GetFrameValue(x, y, t float64) float64 {
	theta := t * s.Rotations * maths.Rotation
	xd, yd := s.Radius*math.Cos(theta), s.Radius*math.Sin(theta)
	return s.Sampler.GetFrameValue(x+xd+s.OffsetX, y+yd+s.OffsetY, s.OffsetT)
}

// Clamps down all values outside MaxRadius of the origin to 0, using a sigmoid with Decay (>6.0) as factor.
// See https://www.desmos.com/calculator/gqy2bw9yt1
type UnitCircleClamper struct {
	Sampler
	MaxRadius float64
	Decay     float64
}

func (s UnitCircleClamper) GetFrameValue(x, y, t float64) float64 {
	return s.Sampler.GetFrameValue(x, y, t) * max((1-2*maths.Sigmoid(s.Decay*(1/(s.MaxRadius)*maths.Radius(x, y)-1))), 0)
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
