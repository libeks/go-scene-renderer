package sampler

import (
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
