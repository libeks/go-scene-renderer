package geometry

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var approxFloatOpt = cmp.Comparer(comparator)

func TestBinomial(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		k          int
		wantResult int
	}{
		{"(1,1)", 1, 1, 1},
		{"(4,2)", 4, 2, 6},
		{"(8,4)", 8, 4, 70},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := binomial(tt.n, tt.k)
			if got != tt.wantResult {
				t.Errorf("wanted %d, got %d", tt.wantResult, got)
			}
		})
	}
}

func TestTFactor(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		i          int
		t          float64
		wantResult float64
	}{
		{"1,0@0.0", 1, 0, 0.0, 1.0},
		{"1,0@1.0", 1, 0, 1.0, 0.0},
		{"1,1@0.0", 1, 1, 0.0, 0.0},
		{"1,1@1.0", 1, 1, 1.0, 1.0},
		{"1,1@0.5", 1, 1, 0.5, 0.5},
		{"2,1@0.5", 2, 1, 0.5, 0.25},
		{"2,0@0.5", 2, 0, 0.5, 0.25},
		{"2,2@0.5", 2, 2, 0.5, 0.25},
		{"2,0@0.1", 2, 0, 0.1, 0.81},
		{"2,2@0.1", 2, 2, 0.1, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tFactor(tt.n, tt.i, tt.t)
			if diff := cmp.Diff(tt.wantResult, got, approxFloatOpt); diff != "" {
				t.Errorf("failure, diff: %s", diff)
			}
		})
	}
}

func TestBezier(t *testing.T) {
	p0 := Point{0, 0, 0}
	p1 := Point{1, 0, 0}
	p2 := Point{1, 1, 0}
	twoPoints := []Point{p0, p1}
	threePoints := []Point{p0, p1, p2}
	tests := []struct {
		name       string
		points     []Point
		t          float64
		wantResult Direction
	}{
		{"0-to-1@0.0", twoPoints, 0.0, Direction{
			Origin: Point{0, 0, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{1, 0, 0},
				UpVector: Vector3D{
					0,
					1,
					0,
				},
				RightVector: Vector3D{
					0, 0, 1,
				},
			},
		}},
		{"0-to-1@0.5", twoPoints, 0.5, Direction{
			Origin: Point{0.5, 0, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{1, 0, 0},
				UpVector: Vector3D{
					0,
					1,
					0,
				},
				RightVector: Vector3D{
					0, 0, 1,
				},
			},
		}},
		{"0-to-1@1.0", twoPoints, 1.0, Direction{
			Origin: Point{1, 0, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{1, 0, 0},
				UpVector: Vector3D{
					0,
					1,
					0,
				},
				RightVector: Vector3D{
					0, 0, 1,
				},
			},
		}},
		{"0-to-2@0.0", threePoints, 0.0, Direction{
			Origin: Point{0, 0, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{1, 0, 0},
				UpVector: Vector3D{
					0,
					1,
					0,
				},
				RightVector: Vector3D{
					0, 0, 1,
				},
			},
		}},
		{"0-to-2@0.5", threePoints, 0.5, Direction{
			Origin: Point{0.75, 0.25, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{0.7071067811865476, 0.7071067811865476, 0},
				UpVector: Vector3D{
					-0.7071067811865476,
					0.7071067811865476,
					0,
				},
				RightVector: Vector3D{
					0, 0, 1,
				},
			},
		}},
		{"0-to-2@1.0", threePoints, 1.0, Direction{
			Origin: Point{1, 1, 0},
			Orientation: EulerDirection{
				ForwardVector: Vector3D{0, 1, 0},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BezierPath{tt.points}.GetDirection(tt.t)
			if diff := cmp.Diff(tt.wantResult, got); diff != "" {
				t.Errorf("failure, diff: %s", diff)
			}
		})
	}
}

func TestInverseDirectionMatrix(t *testing.T) {
	t.Parallel()
	standardOr := OriginPosition.Orientation
	yawLeftOr := OriginPosition.Orientation.ApplyMatrix(
		RotateYaw3D(math.Pi / 2),
	)
	yawRightOr := OriginPosition.Orientation.ApplyMatrix(
		RotateYaw3D(-math.Pi / 2),
	)
	pitchUpOr := OriginPosition.Orientation.ApplyMatrix(
		RotatePitch3D(-math.Pi / 2),
	)
	// pitchDownOr := OriginPosition.Orientation.ApplyMatrix(
	// 	RotatePitch3D(-math.Pi / 2),
	// )
	// fmt.Printf("left %s\n", yawLeftOr)
	sq3 := math.Sqrt(3)
	sq2 := math.Sqrt(2)
	towards111 := standardOr.ApplyMatrix(
		// geometry.RotateMatrixX(-0.615),
		// geometry.RotateMatrixZ(math.Pi/4)
		// RotateYaw3D(math.Asin(1 / sq3)).MatrixMult(RotatePitch3D(-math.Asin(1 / sq2))),
		RotateYaw3D(-math.Asin(1 / sq2)).MatrixMult(RotatePitch3D(-math.Asin(1 / sq3))),
	)
	fmt.Printf("Towards 111: %s\n", towards111)
	// inv, _ := towards111.Inverse3DMatrix().Inverse()
	// fmt.Printf("Towards 111 ^-1: %s\n", inv)
	tests := []struct {
		name        string
		orientation EulerDirection
		in          Vector3D
		want        Vector3D
	}{
		{"origin_x", standardOr, V3(1, 0, 0), V3(1, 0, 0)},
		{"origin_y", standardOr, V3(0, 1, 0), V3(0, 1, 0)},
		{"origin_z", standardOr, V3(0, 0, -1), V3(0, 0, -1)},
		{"111_z", towards111, V3(1/sq3, 1/sq3, -1/sq3), V3(0, 0, -1)},
		{"left_x", yawLeftOr, V3(-1, 0, 0), V3(0, 0, -1)},
		{"left_x2", yawLeftOr, V3(1, 0, 0), V3(0, 0, 1)},
		{"left_z", yawLeftOr, V3(0, 0, 1), V3(-1, 0, 0)},
		{"left_y", yawLeftOr, V3(0, 1, 0), V3(0, 1, 0)},
		{"right_x", yawRightOr, V3(1, 0, 0), V3(0, 0, -1)},
		{"right_z", yawRightOr, V3(0, 0, 1), V3(1, 0, 0)},
		{"right_y", yawRightOr, V3(0, 1, 0), V3(0, 1, 0)},
		{"up_x", pitchUpOr, V3(1, 0, 0), V3(1, 0, 0)},
		{"up_z", pitchUpOr, V3(0, 0, 1), V3(0, 1, 0)},
		{"up_y", pitchUpOr, V3(0, 1, 0), V3(0, 0, -1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.orientation.Inverse3DMatrix()
			got := m.MultVect(tt.in)
			fmt.Printf("m %s, got %s, in %s\n", m, got, tt.in)
			if diff := cmp.Diff(tt.want, got, approxFloatOpt); diff != "" {
				t.Errorf("failure, diff: %s", diff)
			}
		})
	}
}
