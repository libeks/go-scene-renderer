package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
)

type Scene interface {
	// GetPixel operates on the pixel space of the image, starting at (0,0), ending in (width, height)
	// GetPixel(x, y int, t float64) color.Color

	// GetColor operates on a frame where x, y range from -1.0 -> 1.0, centered on (0.0, 0.0)
	GetColor(x, y float64, t float64) color.Color

	// GetColorPalette returns a list of colors that the scene will contain
	// This will help the GIF rendering, as we won't have to perform k-means algo
	// on all colors in every frame, saving time
	// This should ideally be no more than 256 colors
	// An empty slice means k-means will be performed
	GetColorPalette(t float64) []color.Color
}

type PictureFrame struct {
	Width  int
	Height int
}
