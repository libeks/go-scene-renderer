package objects

import (
	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/sampler"
)

// returns an object bounded by x in (-1,1) and z (-1,1) with y value varying based on Perlin noise source
type HeightMap struct {
	Gradient colors.Gradient
	Height   sampler.Sampler
	N        int
}

func (o HeightMap) getAt(x, y, t float64) float64 {
	return o.Height.GetFrameValue(x, y, t)
}

func (o HeightMap) Frame(t float64) StaticObject {
	triangles := []StaticTriangle{}
	zMult := 1.0
	for xd := range o.N {
		for yd := range o.N {
			dx, dy := 2/float64(o.N-1), 2/float64(o.N-1)
			x, y := (2*float64(xd)/float64(o.N-1))-1.0, (2*float64(yd)/float64(o.N-1))-1.0

			a, b, c, d := o.getAt(x, y, t), o.getAt(x, y+dy, t), o.getAt(x+dx, y, t), o.getAt(x+dx, y+dy, t)
			triangles = append(triangles,
				StaticTriangle{
					Triangle: Triangle{
						A: geometry.Pt(x, zMult*a, y),
						B: geometry.Pt(x, zMult*b, y+dy),
						C: geometry.Pt(x+dx, zMult*c, y),
					},
					Colorer: colors.TriangleGradientInterpolationTexture{
						Gradient: o.Gradient,

						A: a, B: b, C: c, D: d,
					},
				},
			)
			triangles = append(triangles,
				StaticTriangle{
					Triangle: Triangle{
						A: geometry.Pt(x+dx, zMult*d, y+dy),
						B: geometry.Pt(x, zMult*b, y+dy),
						C: geometry.Pt(x+dx, zMult*c, y),
					},
					Colorer: colors.TriangleGradientInterpolationTexture{
						Gradient: o.Gradient,

						A: d, B: b, C: c, D: a,
					},
				},
			)
		}
	}
	return StaticObject{
		triangles: triangles,
	}
}

// returns an object bounded by x in (-1,1) and z (-1,1) with y value varying based on Perlin noise source
type HeightMapCircle struct {
	Gradient colors.Gradient
	Height   sampler.Sampler
	N        int
}

func (o HeightMapCircle) getAt(x, y, t float64) float64 {
	return o.Height.GetFrameValue(x, y, t)
}

func (o HeightMapCircle) Frame(t float64) StaticObject {
	triangles := []StaticTriangle{}
	zMult := 1.0
	for xd := range o.N {
		for yd := range o.N {
			dx, dy := 2/float64(o.N-1), 2/float64(o.N-1)
			x, y := (2*float64(xd)/float64(o.N-1))-1.0, (2*float64(yd)/float64(o.N-1))-1.0

			if inCircle(x, y) && inCircle(x+dx, y) && inCircle(x, y+dy) {
				triangles = append(triangles,
					StaticTriangle{
						Triangle: Triangle{
							A: geometry.Pt(x, zMult*o.getAt(x, y, t), y),
							B: geometry.Pt(x, zMult*o.getAt(x, y+dy, t), y+dy),
							C: geometry.Pt(x+dx, zMult*o.getAt(x+dx, y, t), y),
						},
						Colorer: colors.TriangleGradientTexture(
							o.Gradient.Interpolate(o.getAt(x, y, t)),
							o.Gradient.Interpolate(o.getAt(x, y+dy, t)),
							o.Gradient.Interpolate(o.getAt(x+dx, y, t)),
						),
					},
				)
			}
			if inCircle(x+dx, y+dy) && inCircle(x+dx, y) && inCircle(x, y+dy) {
				triangles = append(triangles,
					StaticTriangle{
						Triangle: Triangle{
							A: geometry.Pt(x+dx, zMult*o.getAt(x+dx, y+dy, t), y+dy),
							B: geometry.Pt(x, zMult*o.getAt(x, y+dx, t), y+dy),
							C: geometry.Pt(x+dx, zMult*o.getAt(x+dx, y, t), y),
						},
						Colorer: colors.TriangleGradientTexture(
							o.Gradient.Interpolate(o.getAt(x+dx, y+dy, t)),
							o.Gradient.Interpolate(o.getAt(x, y+dy, t)),
							o.Gradient.Interpolate(o.getAt(x+dx, y, t)),
						),
					},
				)
			}
		}
	}
	return StaticObject{
		triangles: triangles,
	}
}

func inCircle(x, y float64) bool {
	return x*x+y*y < 1.0
}
