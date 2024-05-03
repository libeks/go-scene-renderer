package geometry

import "fmt"

type Line struct {
	A Point
	B Point
}

func (l Line) String() string {
	return fmt.Sprintf("Line %s %s", l.A, l.B)
}

// return a line, cropped to being in front of the screen, and optionally a new point that was created
// as part of the cropping action
func (l Line) CropToFrontOfCamera(minDepth float64) *Line {
	if !l.A.IsInFrontOfCamera(minDepth) && !l.B.IsInFrontOfCamera(minDepth) {
		// both endpoints are behind camera, no line
		return nil
	}
	if l.A.IsInFrontOfCamera(minDepth) && l.B.IsInFrontOfCamera(minDepth) {
		// both lines are fine
		return &l
	}
	// one is behind the screen, the other is not
	start, end := l.A, l.B
	tIntersect := (l.A.Z - minDepth) / (l.A.Z - l.B.Z)
	midpoint := Point{
		l.A.X + tIntersect*(l.B.X-l.A.X),
		l.A.Y + tIntersect*(l.B.Y-l.A.Y),
		l.A.Z + tIntersect*(l.B.Z-l.A.Z),
	}
	if !l.A.IsInFrontOfCamera(minDepth) {
		start = midpoint
	} else {
		end = midpoint
	}
	return &Line{start, end}
}

// a line that could be written to the screen, there is no z-coordinate, only x and y
type RasterLine struct {
	A      Pixel
	ADepth float64
	B      Pixel
	BDepth float64
}

func newRasterLine(aP, bP Point) *RasterLine {
	a, aDepth := aP.ToPixel()
	b, bDepth := bP.ToPixel()
	if a == nil || b == nil {
		panic(fmt.Errorf("Point is unexpectedly behind camera %s %s", aP, bP))
	}
	if a.X == b.X && a.Y == b.Y {
		return nil
	}
	var txmin, txmax, tymin, tymax float64
	if a.X == b.X {
		if a.X > 1 || a.X < -1 {
			return nil
		}
		txmin, txmax = 0, 1
	} else {
		txleft := (-1.0 - a.X) / (b.X - a.X)
		txright := (1.0 - a.X) / (b.X - a.X)
		txmin, txmax = min(txleft, txright), max(txleft, txright)
	}
	if a.Y == b.Y {
		if a.Y > 1 || a.Y < -1 {
			return nil
		}
		tymin, tymax = 0, 1
	} else {
		tyleft := (-1 - a.Y) / (b.Y - a.Y)
		tyright := (1 - a.Y) / (b.Y - a.Y)
		tymin, tymax = min(tyleft, tyright), max(tyleft, tyright)
	}

	// clip pixels to window
	tmin := max(0, txmin, tymin)
	tmax := min(1, txmax, tymax)

	if tmin == 0 && tmax == 1 {
		return &RasterLine{
			A: *a, B: *b,
			ADepth: aDepth, BDepth: bDepth,
		}
	}
	// no part of line is inside the screen
	if tmin > tmax {
		return nil
	}
	return &RasterLine{
		A: Pixel{
			a.X + (b.X-a.X)*tmin,
			a.Y + (b.Y-a.Y)*tmin,
		},
		ADepth: aDepth + (bDepth-aDepth)*tmin,
		B: Pixel{
			a.X + (b.X-a.X)*tmax,
			a.Y + (b.Y-a.Y)*tmax,
		},
		BDepth: aDepth + (bDepth-aDepth)*tmax,
	}
}

func (l RasterLine) String() string {
	return fmt.Sprintf("Raster Line (A: %s, B: %s)", l.A, l.B)
}

func (l Line) CropToScreenView() *RasterLine {
	return newRasterLine(l.A, l.B)
}
