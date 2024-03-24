package colors

import (
	"math"

	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
)

type RotatingLine struct {
	Gradient
	Thickness float64 // in texture coordinates
}

func (t RotatingLine) GetFrameColor(x, y, f float64) Color {
	v := geometry.Vector2D{x*2 - 1, y*2 - 1} // do math in the square -1 to 1
	orth := geometry.Vector2D{1, 0}          // orthogonal vector is rotated 90 degrees from vertical
	m := geometry.RotateMatrix2D(f * maths.Rotation)
	orth = m.MultVect(orth)
	distance := math.Abs(v.DotProduct(orth))
	if distance < t.Thickness && v.Mag() < 1.0 {
		return t.Gradient.Interpolate(0)
	}
	return t.Gradient.Interpolate(1)
}

type DynamicSubtexturer struct {
	Subtexture   AnimatedTexture
	N            int // number of squares to tile
	PointSampler Sampler
}

// returns x/d
func bucketRemainder(x, d float64) (float64, float64) {
	return float64(int(x/d)) * d, math.Mod(x, d) * 1 / d
}

func (s DynamicSubtexturer) GetFrameColor(x, y, t float64) Color {
	d := 1 / float64(s.N)
	xMeta, xValue := bucketRemainder(x, d)
	yMeta, yValue := bucketRemainder(y, d)
	tHere := s.PointSampler.GetFrameValue(xMeta, yMeta, t)
	return s.Subtexture.GetFrameColor(xValue, yValue, tHere)
}
