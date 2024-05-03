package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPointTowards(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		point Point
	}{
		{"1,0,0", Pt(1, 0, 0)},
		{"-1,0,0", Pt(-1, 0, 0)},
		{".5,.5,0", Point(V3(0.5, 0.5, 0).Unit())},
		{".25,0,.25", Point(V3(0.25, 0, 0.25).Unit())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatrix := PointTowards(tt.point)
			if diff := cmp.Diff(gotMatrix.Determinant(), 1.0, approxFloatOpt); diff != "" {
				t.Errorf("determinant is not 1, diff: %s", diff)
			}
			out := gotMatrix.MultVect(V3(0, 0, -1).ToHomogenous())
			if diff := cmp.Diff(tt.point.Vector().ToHomogenous(), out, approxFloatOpt); diff != "" {
				t.Errorf("failure, diff: %s", diff)
			}
		})
	}
}
