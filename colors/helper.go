package colors

import (
	"sync"

	"github.com/libeks/go-scene-renderer/sampler"
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

// returns x/d
func bucketRemainder(x, d float64) (float64, float64) {
	// return float64(int(x/d)) * d, math.Mod(x, d) * 1 / d
	f := float64(int(x / d))
	return f * d, (x - f*d) / d
}

type samplerColorer struct {
	sampler  sampler.Sampler
	gradient Gradient
}

func (s samplerColorer) GetFrameColor(x, y, t float64) Color {
	return s.gradient.Interpolate(s.sampler.GetFrameValue(x, y, t))
}

func GetAniTextureFromSampler(s sampler.Sampler, g Gradient) AnimatedTexture {
	return samplerColorer{
		sampler:  s,
		gradient: g,
	}
}

type v struct {
	x float64
	y float64
	t float64
}

func NewDynamicSubtexturer(s AnimatedTexture, n int, sampler sampler.Sampler) DynamicSubtexturer {
	cache := make(map[v]float64, 0)
	return DynamicSubtexturer{
		Subtexture:   s,
		N:            n,
		PointSampler: sampler,

		cache:   cache,
		RWMutex: &sync.RWMutex{},
	}
}

type DynamicSubtexturer struct {
	Subtexture   AnimatedTexture
	N            int // number of squares to tile
	PointSampler sampler.Sampler

	cache map[v]float64
	*sync.RWMutex
}

func (s DynamicSubtexturer) getCellValue(xMeta, yMeta, t float64) float64 {
	s.RLock()
	val, ok := s.cache[v{x: xMeta, y: yMeta, t: t}]
	s.RUnlock()
	if ok {
		return val
	} else {
		val := s.PointSampler.GetFrameValue(xMeta, yMeta, t)

		s.Lock()
		s.cache[v{x: xMeta, y: yMeta, t: t}] = val
		s.Unlock()
		// s.Unlock()
		return val
	}
}

func (s DynamicSubtexturer) GetFrameColor(x, y, t float64) Color {
	d := 1 / float64(s.N)
	xMeta, xValue := bucketRemainder(x, d)
	yMeta, yValue := bucketRemainder(y, d)
	tHere := s.getCellValue(xMeta, yMeta, t)

	// tHere := s.PointSampler.GetFrameValue(xMeta, yMeta, t)
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
	// t is most likely < 0
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
