package geometry

import "math"

type tuple struct {
	a int
	b int
}

var binomialCache = make(map[tuple]int)

type BezierPath struct {
	Points []Point
}

type EulerDirection struct {
	ForwardVector Vector3D
	UpVector      Vector3D
	RightVector   Vector3D
}

func (d EulerDirection) GetRollPitchYaw() RollPitchYaw {
	yaw := math.Atan2(d.ForwardVector.Z, d.ForwardVector.X)
}

type Direction struct {
	// ForwardVector and UpVector should be orthogonal
	Origin      Point
	Orientation EulerDirection
}

type RollPitchYaw struct {
	// angles expressed in radians
	Roll  float64
	Pitch float64
	Yaw   float64
}

func (p BezierPath) GetDirection(t float64) Direction {
	upVector := Vector3D{0, 1, 0}
	forwardVector := p.direction(t)
	rightVector := forwardVector.CrossProduct(upVector).Unit()
	relativeUpVector := rightVector.CrossProduct(forwardVector).Unit()
	return Direction{
		Origin: p.bezierPoint(t),
		Orientation: EulerDirection{
			ForwardVector: forwardVector,
			UpVector:      relativeUpVector,
			RightVector:   rightVector,
		},
	}
}

func (p BezierPath) bezierPoint(t float64) Point {
	endPoint := Vector3D{0, 0, 0}
	nPoints := len(p.Points)
	for i, pt := range p.Points {
		pointComponent := pt.Vector().ScalarMultiply(float64(binomial(nPoints-1, i)) * tFactor(nPoints-1, i, t))
		endPoint = endPoint.AddVector(pointComponent)
	}
	return Point(endPoint)
}

func (p BezierPath) direction(t float64) Vector3D {
	retVector := Vector3D{0, 0, 0}
	nPoints := len(p.Points)
	for i := range nPoints - 1 {
		pointDifference := p.Points[i+1].Subtract(p.Points[i])
		pointComponent := pointDifference.ScalarMultiply(float64(binomial(nPoints-2, i)) * tFactor(nPoints-2, i, t))
		retVector = retVector.AddVector(pointComponent)
	}
	return retVector.ScalarMultiply(float64(nPoints)).Unit()
}

func tFactor(n, i int, t float64) float64 {
	res := 1.0
	for range i {
		res *= t
	}
	for range n - i {
		res *= (1 - t)
	}
	return res
}

func binomial(n, k int) int {
	a := tuple{n, k}
	if res, ok := binomialCache[a]; ok {
		return res
	}
	if k == 0 || n == k {
		return 1
	}
	res := binomial(n-1, k-1) + binomial(n-1, k)
	binomialCache[a] = res
	return res
}
