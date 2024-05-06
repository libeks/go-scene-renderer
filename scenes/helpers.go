package scenes

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/textures"
)

type textureBackground struct {
	t textures.Texture
}

func (b textureBackground) GetColor(x, y float64) colors.Color {
	// convert from (-1,1) to (0,1)
	return b.t.GetTextureColor(x/2+0.5, y/2+0.5)
}

type dynamicTextureBackground struct {
	t textures.DynamicTexture
}

func (b dynamicTextureBackground) GetFrame(t float64) Background {
	return textureBackground{
		b.t.GetFrame(t),
	}
}

func BackgroundFromTexture(t textures.DynamicTexture) DynamicBackground {
	return dynamicTextureBackground{
		t,
	}
}
