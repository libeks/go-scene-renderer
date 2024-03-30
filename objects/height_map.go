package objects

import (
	// "fmt"
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
)

// returns an object bounded by x in (-1,1) and z (-1,1) with y value varying based on Perlin noise source
type HeightMap struct {
	PerlinNoise colors.PerlinNoise
	N           int
}

func (o HeightMap) Frame(t float64) StaticObject {

	gradient := colors.Grayscale
	triangles := []StaticTriangle{}
	t = t
	zMult := 1.0
	for xd := range o.N {
		for yd := range o.N {
			dx, dy := 2/float64(o.N-1), 2/float64(o.N-1)
			x, y := (2*float64(xd)/float64(o.N-1))-1.0, (2*float64(yd)/float64(o.N-1))-1.0
			triangles = append(triangles,
				StaticTriangle{
					Triangle: Triangle{
						A: geometry.Point{x, zMult * o.PerlinNoise.GetFrameValue(x, y, t), y},
						B: geometry.Point{x, zMult * o.PerlinNoise.GetFrameValue(x, y+dy, t), y + dy},
						C: geometry.Point{x + dx, zMult * o.PerlinNoise.GetFrameValue(x+dx, y, t), y},
					},
					Colorer: colors.TriangleGradientTexture(
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x, y, t)),
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x, y+dy, t)),
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x+dx, y, t)),
					),
				},
				StaticTriangle{
					Triangle: Triangle{
						A: geometry.Point{x + dx, zMult * o.PerlinNoise.GetFrameValue(x+dx, y+dy, t), y + dy},
						B: geometry.Point{x, zMult * o.PerlinNoise.GetFrameValue(x, y+dx, t), y + dy},
						C: geometry.Point{x + dx, zMult * o.PerlinNoise.GetFrameValue(x+dx, y, t), y},
					},
					Colorer: colors.TriangleGradientTexture(
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x+dx, y+dy, t)),
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x, y+dy, t)),
						gradient.Interpolate(o.PerlinNoise.GetFrameValue(x+dx, y, t)),
					),
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
	PerlinNoise colors.PerlinNoise
	N           int
}

func (o HeightMapCircle) getAt(x, y, t float64) float64 {
	// return o.PerlinNoise.GetFrameValue(x, y, t)
	return o.PerlinNoise.GetFrameValue(x, y, t) * max((1-2*maths.Sigmoid(20*(1.05*radius(x, y)-1))), 0)
}

func temp(x, y float64) float64 {
	return (1 - 2*maths.Sigmoid(20*(radius(x, y)-1)))
}

func (o HeightMapCircle) Frame(t float64) StaticObject {
	gradient := colors.Grayscale
	triangles := []StaticTriangle{}
	t = t
	zMult := 1.0

	// fmt.Printf("temp at 0,0 is %0.3f, at 1,0 is %0.3f, at 0.5, 0.5 is %.3f\n", temp(0, 0), temp(1, 0), temp(0.5, 0.5))
	// fmt.Printf("radius at 0,0 is %0.3f, at 1,0 is %0.3f, at 0.5, 0.5 is %.3f\n", radius(0, 0), radius(1, 0), radius(0.5, 0.5))
	for xd := range o.N {
		for yd := range o.N {
			dx, dy := 2/float64(o.N-1), 2/float64(o.N-1)
			x, y := (2*float64(xd)/float64(o.N-1))-1.0, (2*float64(yd)/float64(o.N-1))-1.0
			// fmt.Printf("At %.3f %.3f is %.3f\n", x, y, o.getAt(x, y, t))

			if inCircle(x, y) && inCircle(x+dx, y) && inCircle(x, y+dy) {
				triangles = append(triangles,
					StaticTriangle{
						Triangle: Triangle{
							A: geometry.Point{x, zMult * o.getAt(x, y, t), y},
							B: geometry.Point{x, zMult * o.getAt(x, y+dy, t), y + dy},
							C: geometry.Point{x + dx, zMult * o.getAt(x+dx, y, t), y},
						},
						Colorer: colors.TriangleGradientTexture(
							gradient.Interpolate(o.getAt(x, y, t)),
							gradient.Interpolate(o.getAt(x, y+dy, t)),
							gradient.Interpolate(o.getAt(x+dx, y, t)),
						),
					},
				)
			}
			if inCircle(x+dx, y+dy) && inCircle(x+dx, y) && inCircle(x, y+dy) {
				triangles = append(triangles,
					StaticTriangle{
						Triangle: Triangle{
							A: geometry.Point{x + dx, zMult * o.getAt(x+dx, y+dy, t), y + dy},
							B: geometry.Point{x, zMult * o.getAt(x, y+dx, t), y + dy},
							C: geometry.Point{x + dx, zMult * o.getAt(x+dx, y, t), y},
						},
						Colorer: colors.TriangleGradientTexture(
							gradient.Interpolate(o.getAt(x+dx, y+dy, t)),
							gradient.Interpolate(o.getAt(x, y+dy, t)),
							gradient.Interpolate(o.getAt(x+dx, y, t)),
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

func radius(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}
