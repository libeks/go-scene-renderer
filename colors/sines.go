package colors

import (
	"math"

	"github.com/libeks/go-scene-renderer/maths"
)

type SineWaveAnimation struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     Gradient
}

func (s SineWaveAnimation) GetFrameColor(x, y, t float64) Color {
	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio)
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
	return GrayscaleColor(valZeroOne)
}

type SineWaveWCrossAnimation struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	Gradient     Gradient
	TOffset      float64
	TScale       float64
}

func (s SineWaveWCrossAnimation) GetTextureColor(x, y float64, t float64) Color {
	// tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	// waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
	waveComponent := 0.0
	crossComponent := (math.Abs(x * y)) / ((t + s.TOffset) * s.TScale)
	valMinOneToOne := math.Cos(waveComponent + crossComponent)
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
	return s.Gradient.Interpolate(valZeroOne)
}

func AnimatedToDynamicTexture(ani AnimatedTexture) DynamicTexture {
	return dynamicTextureHelper{
		ani,
	}
}
