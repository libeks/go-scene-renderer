package sampler

import (
	// "fmt"
	"math"

	"github.com/libeks/go-scene-renderer/maths"
)

// gets values from rotating around the OffsetX,Y by distance of Radius, parameterized by t,
// a total of Rotation rotations (negative to rotate in opposite direction)
type RotatingSampler struct {
	DynamicSampler
	Rotations float64
	Radius    float64
	OffsetX   float64
	OffsetY   float64
	OffsetT   float64
}

func (s RotatingSampler) GetFrame(t float64) StaticSampler {
	theta := t * s.Rotations * maths.Rotation
	xd, yd := s.Radius*math.Cos(theta), s.Radius*math.Sin(theta)
	// return s.Sampler.GetFrameValue(x+xd+s.OffsetX, y+yd+s.OffsetY, s.OffsetT)
	return shiftedStatic{s.DynamicSampler.GetFrame(t), xd + s.OffsetX, yd + s.OffsetY}
}

func (s RotatingSampler) GetFrameValue(x, y, t float64) float64 {
	return s.GetFrame(t).GetValue(x, y)
}

// Clamps down all values outside MaxRadius of the origin to 0, using a sigmoid with Decay (>6.0) as factor.
// See https://www.desmos.com/calculator/gqy2bw9yt1
type UnitCircleClamper struct {
	DynamicSampler
	MaxRadius float64
	Decay     float64
}

func (s UnitCircleClamper) GetFrame(t float64) StaticSampler {
	return UnitCircleClamperStatic{
		StaticSampler: s.GetFrame(t),
		MaxRadius:     s.MaxRadius,
		Decay:         s.Decay,
	}
}

type UnitCircleClamperStatic struct {
	StaticSampler
	MaxRadius float64
	Decay     float64
}

func (s UnitCircleClamperStatic) GetValue(x, y float64) float64 {
	return s.StaticSampler.GetValue(x, y) * max((1-2*maths.Sigmoid(s.Decay*(1/(s.MaxRadius)*maths.Radius(x, y)-1))), 0)
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

func RotatingStatic(sampler StaticSampler, nRotations float64) rotatingStatic {
	return rotatingStatic{
		StaticSampler: sampler,
		NRotations:    nRotations,
	}
}

type rotatingStatic struct {
	StaticSampler
	NRotations float64
}

func (r rotatingStatic) GetFrameValue(x, y, t float64) float64 {
	return r.GetFrame(t).GetValue(x, y)
}

func (r rotatingStatic) GetFrame(t float64) StaticSampler {
	angle := r.NRotations * maths.Rotation * t
	return RotatedStatic(r.StaticSampler, angle)
}

func RotatedStatic(sampler StaticSampler, angle float64) StaticSampler {
	return rotatedStatic{
		StaticSampler: sampler,
		cos:           math.Cos(angle),
		sin:           math.Sin(angle),
	}
}

type rotatedStatic struct {
	StaticSampler
	cos float64
	sin float64
}

func (s rotatedStatic) GetValue(x, y float64) float64 {
	x, y = x*2-1, y*2-1
	x, y = s.cos*x+s.sin*y, -s.sin*x+s.cos*y
	return s.StaticSampler.GetValue(x/2+0.5, y/2+0.5)
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

func Shifted(s Sampler, x, y float64) Sampler {
	return shifted{
		Sampler: s,
		xOffset: x,
		yOffset: y,
	}
}

type shifted struct {
	Sampler
	xOffset float64
	yOffset float64
}

func (s shifted) GetFrameValue(x, y float64, t float64) float64 {
	return s.Sampler.GetFrameValue(x+s.xOffset, y+s.yOffset, t)
}

func (s shifted) GetFrame(t float64) StaticSampler {
	return shiftedStatic{
		StaticSampler: s.GetFrame(t),
		xOffset:       s.xOffset,
		yOffset:       s.yOffset,
	}
}

type shiftedStatic struct {
	StaticSampler
	xOffset float64
	yOffset float64
}

func (s shiftedStatic) GetValue(x, y float64) float64 {
	return s.StaticSampler.GetValue(x+s.xOffset, y+s.yOffset)
}

func TimeShifted(s Sampler, t float64) Sampler {
	return timeShifted{
		Sampler:    s,
		timeOffset: t,
	}
}

type timeShifted struct {
	Sampler
	timeOffset float64
}

func (s timeShifted) GetFrameValue(x, y float64, t float64) float64 {
	return s.Sampler.GetFrameValue(x, y, t+s.timeOffset)
}

func TimeShiftedDynamic(s DynamicSampler, t float64) DynamicSampler {
	return timeShiftedDynamic{
		DynamicSampler: s,
		timeOffset:     t,
	}
}

type timeShiftedDynamic struct {
	DynamicSampler
	timeOffset float64
}

func (s timeShiftedDynamic) GetFrame(t float64) StaticSampler {
	return s.DynamicSampler.GetFrame(t + s.timeOffset)
}
