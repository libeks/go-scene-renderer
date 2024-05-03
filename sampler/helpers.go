package sampler

import (
	// "fmt"
	"math"

	"github.com/libeks/go-scene-renderer/maths"
)

// gets values from rotating around the OffsetX,Y by distance of Radius, parameterized by t,
// a total of Rotation rotations (negative to rotate in opposite direction)
type RotatingSampler struct {
	Sampler
	Rotations float64
	Radius    float64
	OffsetX   float64
	OffsetY   float64
	OffsetT   float64
}

func (s RotatingSampler) GetFrameValue(x, y, t float64) float64 {
	theta := t * s.Rotations * maths.Rotation
	xd, yd := s.Radius*math.Cos(theta), s.Radius*math.Sin(theta)
	return s.Sampler.GetFrameValue(x+xd+s.OffsetX, y+yd+s.OffsetY, s.OffsetT)
}

// Clamps down all values outside MaxRadius of the origin to 0, using a sigmoid with Decay (>6.0) as factor.
// See https://www.desmos.com/calculator/gqy2bw9yt1
type UnitCircleClamper struct {
	Sampler
	MaxRadius float64
	Decay     float64
}

func (s UnitCircleClamper) GetFrameValue(x, y, t float64) float64 {
	return s.Sampler.GetFrameValue(x, y, t) * max((1-2*maths.Sigmoid(s.Decay*(1/(s.MaxRadius)*maths.Radius(x, y)-1))), 0)
}

type Constant struct {
	Val float64
}

func (s Constant) GetFrameValue(x, y, t float64) float64 {
	return s.Val
}

type Sigmoid struct {
	Sampler
	Ratio float64
}

func (s Sigmoid) GetFrameValue(x, y, t float64) float64 {
	return maths.Sigmoid(s.Sampler.GetFrameValue(x, y, t) * s.Ratio)
}

type MinusPlusToZeroOne struct {
	Sampler
}

func (s MinusPlusToZeroOne) GetFrameValue(x, y, t float64) float64 {
	return s.Sampler.GetFrameValue(x, y, t)/2 + 0.5
}

type Scalar struct {
	Sampler
	Factor float64
}

func (s Scalar) GetFrameValue(x, y, t float64) float64 {
	return s.Sampler.GetFrameValue(x, y, t) * s.Factor
}

type Sinus struct {
	Sampler
}

func (s Sinus) GetFrameValue(x, y, t float64) float64 {
	return math.Sin(s.Sampler.GetFrameValue(x, y, t))
}

type SineWave struct {
	Factor float64
}

func (s SineWave) GetFrameValue(x, y, t float64) float64 {
	return math.Sin(s.Factor * (x + y))
}

type SineOtherDirectionWave struct {
	Factor float64
}

func (s SineOtherDirectionWave) GetFrameValue(x, y, t float64) float64 {
	return math.Sin(s.Factor * (x - y))
}

type SineWavy struct {
}

func (s SineWavy) GetFrameValue(x, y, t float64) float64 {
	return 2 * math.Sin(x*45+math.Sin(16*y*math.Pi))
}

type Rotated struct {
	Sampler
	Angle float64 // in radians
}

// assume it is evaluated on (0,1)
func (s Rotated) GetFrameValue(x, y, t float64) float64 {
	x, y = x*2-1, y*2-1
	x, y = math.Cos(s.Angle)*x+math.Sin(s.Angle)*y, -math.Sin(s.Angle)*x+math.Cos(s.Angle)*y
	return s.Sampler.GetFrameValue(x, y, t)
}

type Wiggle struct {
	Sampler
	NWiggles float64
	Angle    float64 // max angle
}

func (s Wiggle) GetFrameValue(x, y, t float64) float64 {
	return Rotated{s.Sampler, math.Sin(2*t*s.NWiggles*math.Pi) * s.Angle}.GetFrameValue(x, y, t)
}

type SineWaveAnimation struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
}

func (s SineWaveAnimation) GetFrameValue(x, y, t float64) float64 {
	tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	valMinOneToOne := math.Sin(float64(t)/tRatio + float64(x+y)/s.XYRatio)
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
	return valZeroOne
}

type SineWaveWCrossAnimation struct {
	XYRatio      float64
	SigmoidRatio float64
	SinCycles    int
	TOffset      float64
	TScale       float64
}

func (s SineWaveWCrossAnimation) GetFrameValue(x, y float64, t float64) float64 {
	// tRatio := 1 / (2 * math.Pi * float64(s.SinCycles))
	// waveComponent := float64(t)/tRatio + float64(x+y)/s.XYRatio
	waveComponent := 0.5
	crossComponent := (math.Abs(x * y)) / ((t + s.TOffset) * s.TScale)
	valMinOneToOne := math.Cos(waveComponent + crossComponent)
	valZeroOne := maths.Sigmoid(valMinOneToOne * s.SigmoidRatio)
	return valZeroOne
}
