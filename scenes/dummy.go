package scenes

import "image/color"

type Dummy struct{}

func (d Dummy) GetPixel(x, y int, t float32) color.Color {
	return color.Gray{
		Y: uint8(256 * t),
	}
}
