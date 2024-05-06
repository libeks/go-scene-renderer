package textures

import "github.com/libeks/go-scene-renderer/sampler"

type Transparency interface {
	// true means the object is opaque, false means it is transparent
	GetAlpha(b, c float64) bool
}

type DynamicTransparency interface {
	GetFrame(t float64) Transparency
}

type AnimatedTransparency interface {
	GetAlpha(b, c, t float64) bool
}

// a helper for when a static texture is needed as a dynamic texture
type staticTransparency struct {
	t Transparency
}

func (tr staticTransparency) GetFrame(f float64) Transparency {
	return tr.t
}

func StaticTransparency(t Transparency) DynamicTransparency {
	return staticTransparency{t}
}

type constantTransparency struct {
	val bool
}

func (tr constantTransparency) GetAlpha(b, c float64) bool {
	return tr.val
}

func Opaque() Transparency {
	return constantTransparency{val: true}
}

func DynamicFromAnimatedTransparency(ani AnimatedTransparency) DynamicTransparency {
	return animatedTransparencyHelper{
		ani,
	}
}

type animatedTransparencyHelper struct {
	ani AnimatedTransparency
}

func (tr animatedTransparencyHelper) GetFrame(t float64) Transparency {
	return animatedStaticTransparencyHelper{
		ani: tr.ani,
		t:   t,
	}
}

type animatedStaticTransparencyHelper struct {
	ani AnimatedTransparency
	t   float64
}

func (tr animatedStaticTransparencyHelper) GetAlpha(b, c float64) bool {
	return tr.ani.GetAlpha(b, c, tr.t)
}

type samplerTransparency struct {
	sampler.Sampler
	threshold float64
}

func (tr samplerTransparency) GetAlpha(b, c, t float64) bool {
	return tr.Sampler.GetFrameValue(b, c, t) > tr.threshold
}

func SamplerTransparency(s sampler.Sampler, threshold float64) AnimatedTransparency {
	return samplerTransparency{
		Sampler:   s,
		threshold: threshold,
	}
}

type CircleCutout struct {
	Radius float64
}

func (tr CircleCutout) GetAlpha(b, c, t float64) bool {
	b, c = 2*b-1, 2*c-1
	return b*b+c*c > tr.Radius*tr.Radius
}

func InvertTransparency(ani AnimatedTransparency) AnimatedTransparency {
	return invertedAnimatedTransparency{AnimatedTransparency: ani}
}

type invertedAnimatedTransparency struct {
	AnimatedTransparency
}

func (tr invertedAnimatedTransparency) GetAlpha(b, c, t float64) bool {
	return !tr.AnimatedTransparency.GetAlpha(b, c, t)
}
