package renderer

import (
	"image"

	"github.com/libeks/go-scene-renderer/colors"
)

// func ClipRasterLine(a,b RasterPixel, width, height int) Rasterline {

// }

// func NewRasterLine(a, b RasterPixel) *RasterLine {
// 	if a.X == b.X && a.Y == b.Y {
// 		return nil
// 	}
// 	var txmin, txmax, tymin, tymax float64
// 	if a.X == b.X {
// 		if a.X > 1 || a.X < -1 {
// 			return nil
// 		}
// 		txmin, txmax = 0, 1
// 	} else {
// 		txleft := (-1.0 - a.X) / (b.X - a.X)
// 		txright := (1.0 - a.X) / (b.X - a.X)
// 		txmin, txmax = min(txleft, txright), max(txleft, txright)
// 	}
// 	if a.Y == b.Y {
// 		if a.Y > 1 || a.Y < -1 {
// 			return nil
// 		}
// 		txmin, txmax = 0, 1
// 	} else {
// 		txleft := (-1 - a.X) / (b.X - a.X)
// 		txright := (1 - a.X) / (b.X - a.X)
// 		txmin, txmax = min(txleft, txright), max(txleft, txright)
// 	}

// 	// clip pixels to window
// 	tmin := max(0, txmin, tymin)
// 	tmax := min(1, txmax, tymax)
// 	if tmin == 0 && tmax == 1 {
// 		return &RasterLine{
// 			a, b,
// 		}
// 	}
// 	// no part of line is inside the screen
// 	if tmin > tmax {
// 		return nil
// 	}
// 	return &RasterLine{
// 		A: RasterPixel{
// 			a.X + (b.X-a.X)*tmin,
// 			a.Y + (b.Y-a.Y)*tmin,
// 		},
// 		B: RasterPixel{
// 			a.X + (b.X-a.X)*tmax,
// 			a.Y + (b.Y-a.Y)*tmax,
// 		},
// 	}
// }

func NewRasterLine(a, b RasterPixel) *RasterLine {
	return &RasterLine{
		a, b,
	}
}

type RasterLine struct {
	A RasterPixel
	B RasterPixel
}

type RasterPixel struct {
	X int
	Y int
}

func NewImage(ip ImagePreset) *Image {
	return &Image{
		im: image.NewRGBA(
			image.Rect(
				0, 0, ip.width, ip.height,
			),
		),
		ip: ip,
	}
}

type Image struct {
	im *image.RGBA
	ip ImagePreset
}

// insert pixels with flipped y- coord, so y would be -1 at the bottom, +1 at the top of the image
func (i *Image) Set(x, y int, c colors.Color) {
	i.im.Set(x, i.ip.height-y-1, c)
}

func (i *Image) GetImage() image.Image {
	return i.im
}

// adapted from https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (i *Image) RenderLine(line *RasterLine, gradient colors.Gradient) {
	if line == nil {
		return
	}
	x0, y0, x1, y1 := line.A.X, line.A.Y, line.B.X, line.B.Y
	dx := abs(x1 - x0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	dy := -abs(y1 - y0)
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	error := dx + dy

	xprogress := float64(0)
	for {
		var c colors.Color
		if dx == 0 {
			c = gradient.Interpolate(0.0)
		} else {
			c = gradient.Interpolate(xprogress / float64(dx))
		}

		i.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * error
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			error = error + dy
			x0 += sx
			xprogress += 1
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			error = error + dx
			y0 += sy
		}
	}
}

func (i *Image) Fill(c colors.Color) {
	for x := 0; x < i.ip.width; x++ {
		for y := 0; y < i.ip.height; y++ {
			i.Set(x, y, c)
		}
	}
}
