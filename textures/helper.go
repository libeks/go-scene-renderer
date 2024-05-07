package textures

import (
	"fmt"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/sampler"
)

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

func (f dynamicTextureFrameHelper) GetTextureColor(x, y float64) colors.Color {
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
	gradient colors.Gradient
}

func (s samplerColorer) GetFrameColor(x, y, t float64) colors.Color {
	return s.gradient.Interpolate(s.sampler.GetFrameValue(x, y, t))
}

func GetAniTextureFromSampler(s sampler.Sampler, g colors.Gradient) AnimatedTexture {
	return samplerColorer{
		sampler:  s,
		gradient: g,
	}
}

type TextureValueMapping struct {
	Above float64
	Texture
}

// StaticMapper displays the static Texture in the list, the first one whose Above value is below t
type StaticMapper struct {
	Mapping []TextureValueMapping // ordered in decreasing order of Above
}

func (m StaticMapper) GetFrameColor(x, y, t float64) colors.Color {
	for _, mapping := range m.Mapping {
		if t >= mapping.Above {
			return mapping.Texture.GetTextureColor(x, y)
		}
	}
	// t is most likely < 0
	return colors.Red // shouldn't ever happen if the last Mapping starts at 0.0
}

// implements AnimatedTexture
type BinaryAnimatedSamplerWithColors struct {
	sampler.Sampler
	On  colors.Color
	Off colors.Color
}

func (s BinaryAnimatedSamplerWithColors) GetFrameColor(x, y, t float64) colors.Color {
	val := s.Sampler.GetFrameValue(x, y, t)
	if val == 1 {
		return s.On
	}
	if val == 0 {
		return s.Off
	}
	panic(fmt.Errorf("BinarySampler returned value %.3f", val))
}

// implements AnimatedTexture
type BinaryDynamicSamplerWithColors struct {
	sampler.DynamicSampler
	On  colors.Color
	Off colors.Color
}

func (s BinaryDynamicSamplerWithColors) GetFrameColor(x, y, t float64) colors.Color {
	return s.GetFrame(t).GetTextureColor(x, y)
}

func (s BinaryDynamicSamplerWithColors) GetFrame(t float64) Texture {
	frame := s.DynamicSampler.GetFrame(t)
	return BinarySamplerWithColors{
		StaticSampler: frame,
		On:            s.On,
		Off:           s.Off,
	}
}

type BinarySamplerWithColors struct {
	sampler.StaticSampler
	On  colors.Color
	Off colors.Color
}

// Implements Texture
func (s BinarySamplerWithColors) GetTextureColor(x, y float64) colors.Color {
	val := s.StaticSampler.GetValue(x, y)
	if val == 1 {
		return s.On
	}
	if val == 0 {
		return s.Off
	}
	panic(fmt.Errorf("BinarySampler returned value %.3f", val))
}
