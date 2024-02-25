package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
)

type Dummy struct{}

func (d Dummy) GetPixel(x, y int, t float64) color.Color {
	return color.GrayscaleColor(t)
}

func (d Dummy) GetColor(x, y float64, t float64) color.Color {
	return color.GrayscaleColor(t)
}
