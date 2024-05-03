package geometry

import (
	"fmt"
	"math"
)

type tuple struct {
	a int
	b int
}

var (
	OriginPosition = Direction{
		Origin: Point{0, 0, 0},
		Orientation: EulerDirection{
			V3(0, 0, -1), // negative since the camera points in negative z-direction
			V3(0, 1, 0),
			V3(1, 0, 0)},
	}

	binomialCache map[tuple]int
)

func init() {
	binomialCache = make(map[tuple]int)
}

type Path interface {
	// at every frame, get a camera position and direction it's pointing in
	GetDirection(float64) Direction
}

type EulerDirection struct {
	ForwardVector Vector3D
	UpVector      Vector3D
	RightVector   Vector3D
}

func (d EulerDirection) String() string {
	return fmt.Sprintf("EulerDirection: {Forward: %s, Up: %s, Right: %s}", d.ForwardVector, d.UpVector, d.RightVector)
}

func (d EulerDirection) ApplyMatrix(m Matrix3D) EulerDirection {
	return EulerDirection{
		m.MultVect(d.ForwardVector).Unit(),
		m.MultVect(d.UpVector).Unit(),
		m.MultVect(d.RightVector).Unit(),
	}
}

func (d EulerDirection) Inverse3DMatrix() Matrix3D {
	// return the matrix that transforms the x,y,z vectors into the given direction

	m := Matrix3D{
		d.RightVector.X,
		d.UpVector.X,
		-d.ForwardVector.X, // negative since the camera points in negative z-direction

		d.RightVector.Y,
		d.UpVector.Y,
		-d.ForwardVector.Y, // negative since the camera points in negative z-direction

		d.RightVector.Z,
		d.UpVector.Z,
		-d.ForwardVector.Z, // negative since the camera points in negative z-direction
	}
	if !comparator(m.Determinant(), 1) {
		fmt.Printf("ALERT! CAMERA TRANSFORMATION MATRIX HAS DETERMINANT %v\n", m.Determinant())
	}
	m, valid := m.Inverse()
	if !valid {
		return Identity3D
	}
	// m = m.ScalarMult(-1)
	return m
}

func (d EulerDirection) InverseHomoMatrix() HomogeneusMatrix {
	// return the matrix that transforms the x,y,z vectors into the given direction
	return d.Inverse3DMatrix().toHomogenous()
}

// func (d EulerDirection) GetRollPitchYaw() RollPitchYaw {
// 	fmt.Printf("Initially we have %s\n", d)
// 	// zero out the forward vector's X-coord by yawing
// 	yaw := math.Atan2(d.ForwardVector.X, d.ForwardVector.Z)
// 	fmt.Printf("Yaw at %.3f %.3f is %.3f\n", d.ForwardVector.X, d.ForwardVector.Z, yaw)
// 	newDirection := d.ApplyMatrix(RotateYaw3D(-yaw))
// 	// forward vector's Z-coord is now 0
// 	// now pitch the forward vector to have a zero Y-coord
// 	fmt.Printf("After un-yaw, we have %s\n", newDirection)
// 	pitch := math.Atan2(newDirection.ForwardVector.Y, newDirection.ForwardVector.Z)
// 	fmt.Printf("Pitch at %.3f %.3f is %.3f\n", newDirection.ForwardVector.Y, newDirection.ForwardVector.X, pitch)
// 	newDirection = d.ApplyMatrix(RotatePitch3D(-pitch))
// 	fmt.Printf("After un-pitching, we have %s\n", newDirection)
// 	// now Forward vector points along the Z axis (perpendicularly into the image plane)
// 	roll := math.Atan2(newDirection.RightVector.Y, newDirection.RightVector.X)
// 	fmt.Printf("Roll at %.3f %.3f is %.3f\n", newDirection.RightVector.Y, newDirection.RightVector.X, roll)
// 	newDirection = d.ApplyMatrix(RotateRoll3D(-roll))
// 	fmt.Printf("After a full unroll we have %s\n", newDirection)
// 	return RollPitchYaw{
// 		Roll:  roll,
// 		Pitch: pitch,
// 		Yaw:   yaw,
// 	}
// }

type Direction struct {
	// ForwardVector and UpVector should be orthogonal
	Origin      Point
	Orientation EulerDirection
}

func (d Direction) InverseHomoMatrix() HomogeneusMatrix {
	return d.Orientation.InverseHomoMatrix().MatrixMult(TranslationMatrix(Vector3D(d.Origin).ScalarMultiply(-1)))
}

type RollPitchYaw struct {
	// angles expressed in radians
	Roll  float64
	Pitch float64
	Yaw   float64
}

type BezierPath struct {
	Points []Point
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

func SamplePath(path BezierPath, start, end float64) sampledPath {
	return sampledPath{
		BezierPath: path,
		start:      start,
		end:        end,
	}
}

type sampledPath struct {
	BezierPath
	start float64
	end   float64
}

func (s sampledPath) GetDirection(t float64) Direction {
	delta := s.end - s.start
	newT := (t * delta) + s.start
	return s.BezierPath.GetDirection(newT)
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

func comparator(a, b float64) bool {
	// adapted from https://stackoverflow.com/a/33024979, algo used in Python 3.5
	relativeTolerance := 1e-9
	absTolerance := 1e-9
	return math.Abs(a-b) <= max(relativeTolerance*max(math.Abs(a), math.Abs(b)), absTolerance)
}
