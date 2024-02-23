package scenes

import (
	"image/color"
	"math"
)

type SineWave struct {
	TRatio       float64
	XYRatio      float64
	SigmoidRatio float64
}

func (s SineWave) GetPixel(x, y int, t float64) color.Color {
	x = x - 300
	y = y - 300
	valMinOneToOne := math.Sin(float64(t)/s.TRatio + float64(x+y)/s.XYRatio)
	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
	return color.Gray{
		Y: uint8(256 * (valZeroOne)),
	}
}

type SineWaveWBump struct {
	Frame        PictureFrame
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
}

func (s SineWaveWBump) GetPixel(x, y int, t float64) color.Color {
	x = x - s.Frame.Width/2
	y = y - s.Frame.Height/2
	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio + float64(1/(math.Abs(float64(x*y))+1)/(1/(t*10000))))
	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
	return grayscaleColor(gammaCorrect(valZeroOne))
}

func grayscaleColor(v float64) color.Color {
	return color.Gray{
		Y: uint8(256 * v),
	}
}

func minPlusOneToZO(v float64) float64 {
	return v/2 + 0.5
}

func sigmoid(v float64) float64 {
	return 1 / (1 + math.Exp(-v))
}

func gammaCorrect(v float64) float64 {
	return math.Pow(v, 1.8)
}

// func (d SineWave) GetPixel(x, y int, t float32) color.Color {
// 	return color.Gray{
// 		Y: uint8(256 * (math.Sin(float64(t/10+float32(x)/200)) + math.Cos(float64(t/10+float32(y)/200)))),
// 	}
// }
