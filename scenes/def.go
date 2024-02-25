package scenes

import (
	"github.com/libeks/go-scene-renderer/color"
)

type Scene interface {
	// GetColor operates on a frame where x, y range from -1.0 -> 1.0, centered on (0.0, 0.0)
	GetColor(x, y float64, t float64) color.Color

	// GetColorPalette returns a list of colors that the scene will contain
	// This will help the GIF rendering, as we won't have to perform k-means algo
	// on all colors in every frame, saving time
	// This should ideally be no more than 256 colors
	// An empty slice means k-means will be performed
	GetColorPalette(t float64) []color.Color
}

type Object interface {
	// returns the color of the object at a ray
	// emanating from the camera at (0,0,0), pointed in the direction
	// (x,y, -1), with perspective
	// if there is no intersection, return nil
	GetColor(x, y float64) *color.Color

	// // Returns true if object intersects the ray from (0,0,0) in the direction
	// // (x, y, -1)
	// DoesIntersect(x, y float64) bool
}
