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
