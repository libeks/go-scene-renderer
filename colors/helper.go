package colors

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
