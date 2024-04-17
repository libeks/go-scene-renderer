package geometry

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var approxFloatOpt = cmp.Comparer(func(x, y float64) bool {
	if x-y == 0.0 {
		return true
	}
	delta := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	return delta/mean < 0.1
})

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
			Origin:        Point{0, 0, 0},
			ForwardVector: Vector3D{1, 0, 0},
		}},
		{"0-to-1@0.5", twoPoints, 0.5, Direction{
			Origin:        Point{0.5, 0, 0},
			ForwardVector: Vector3D{1, 0, 0},
		}},
		{"0-to-1@1.0", twoPoints, 1.0, Direction{
			Origin:        Point{1, 0, 0},
			ForwardVector: Vector3D{1, 0, 0},
		}},
		{"0-to-2@0.0", threePoints, 0.0, Direction{
			Origin:        Point{0, 0, 0},
			ForwardVector: Vector3D{1, 0, 0},
		}},
		{"0-to-2@0.5", threePoints, 0.5, Direction{
			Origin:        Point{0.75, 0.25, 0},
			ForwardVector: Vector3D{0.7071067811865476, 0.7071067811865476, 0},
		}},
		{"0-to-2@1.0", threePoints, 1.0, Direction{
			Origin:        Point{1, 1, 0},
			ForwardVector: Vector3D{0, 1, 0},
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
