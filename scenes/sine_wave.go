package scenes

import (
	"math"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
)

type SineWave struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     color.Gradient
}

func (s SineWave) GetFrame(t float64) Frame {
	return FrameHelper{
		t,
		s,
	}
}

func (s SineWave) GetFrameColor(x, y, t float64) color.Color {
	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio)
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
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
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
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

func (f FrameHelper) Flatten() []*objects.Triangle {
	return nil
}
