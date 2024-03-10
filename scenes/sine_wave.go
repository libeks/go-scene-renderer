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
	TOffset      float64
	TScale       float64
}

func (s SineWaveWCross) GetFrameColor(x, y float64, t float64) color.Color {
	// tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	// waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
	waveComponent := 0.0
	crossComponent := (math.Abs(x * y)) / ((t + s.TOffset) * s.TScale)
	valMinOneToOne := math.Cos(waveComponent + crossComponent)
	valZeroOne := sigmoid(valMinOneToOne * s.SigmoidRatio)
	return s.Gradient.Interpolate(valZeroOne)
}

func (s SineWaveWCross) GetColorPalette(t float64) []color.Color {
	return color.GetGradientColorPalette(s.Gradient)
}

func (s SineWaveWCross) GetFrame(t float64) Frame {
	return FrameHelper{t, s}
}

type FrameHelper struct {
	t     float64
	scene AnimatedScene
}

func (f FrameHelper) GetColor(x, y float64) color.Color {
	return f.scene.GetFrameColor(x, y, f.t)
}

func minPlusOneToZO(v float64) float64 {
	return v/2 + 0.5
}

func sigmoid(v float64) float64 {
	// takes from (-inf, +int) to (0.0, 1.0), with an S-like shape centered on 0.0.
	return 1 / (1 + math.Exp(-v))
}
