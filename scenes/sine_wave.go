package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/color"
)

type SineWave struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     color.Gradient
}

// func (s SineWave) GetPixel(x, y int, t float64) color.Color {
// 	x = x - 300
// 	y = y - 300
// 	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
// 	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio)
// 	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
// 	return color.GrayscaleColor(valZeroOne)
// }

func (s SineWave) GetColor(x, y, t float64) color.Color {
	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio)
	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
	return color.GrayscaleColor(valZeroOne)
}

type SineWaveWCross struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     color.Gradient
}

// func (s SineWaveWCross) GetPixel(x, y int, t float64) color.Color {
// 	x = x - s.Frame.Width/2
// 	y = y - s.Frame.Height/2
// 	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
// 	waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
// 	crossComponent := float64(1 / (math.Abs(float64(x*y/6000)) + 1) / (1 / (t * 10000)))
// 	valMinOneToOne := math.Sin(waveComponent + crossComponent)
// 	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
// 	return s.Gradient.Interpolate(valZeroOne)
// }

func (s SineWaveWCross) GetColor(x, y float64, t float64) color.Color {
	// tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	// waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
	waveComponent := 0.0
	crossComponent := (math.Abs(x * y)) / (t * 0.03)
	valMinOneToOne := math.Cos(waveComponent + crossComponent)
	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
	return s.Gradient.Interpolate(valZeroOne)
}

func (s SineWaveWCross) GetColorPalette(t float64) []color.Color {
	return color.GetGradientColorPalette(s.Gradient)
}

type SineWaveWBump struct {
	Frame        PictureFrame
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     color.Gradient
}

// func (s SineWaveWBump) GetPixel(x, y int, t float64) color.Color {
// 	x = x - s.Frame.Width/2
// 	y = y - s.Frame.Height/2
// 	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
// 	waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
// 	bumpComponent := float64(float64((math.Pow(float64(x), float64(2))) * (math.Pow(float64(y), float64(2)))))
// 	valMinOneToOne := math.Sin(waveComponent + 100/bumpComponent)
// 	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
// 	return s.Gradient.Interpolate(valZeroOne)
// }

func minPlusOneToZO(v float64) float64 {
	return v/2 + 0.5
}

func sigmoid(v float64) float64 {
	// takes from (-inf, +int) to (0.0, 1.0), with an S-like shape centered on 0.0.
	return 1 / (1 + math.Exp(-v))
}
