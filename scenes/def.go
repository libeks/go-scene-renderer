package scenes

import (
	"image/color"
)

type Scene interface {
	GetPixel(x, y int, t float64) color.Color
}

type PictureFrame struct {
	Width  int
	Height int
}
