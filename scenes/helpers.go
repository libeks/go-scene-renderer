package scenes

import "github.com/libeks/go-scene-renderer/colors"

type textureBackground struct {
	t colors.Texture
}

func (b textureBackground) GetColor(x, y float64) colors.Color {
	return b.t.GetTextureColor(2*x-1, 2*y-1)
}

type dynamicTextureBackground struct {
	t colors.DynamicTexture
}

func (b dynamicTextureBackground) GetFrame(t float64) Background {
	return textureBackground{
		b.t.GetFrame(t),
	}
}

func BackgroundFromTexture(t colors.DynamicTexture) DynamicBackground {
	return dynamicTextureBackground{
		t,
	}
}
