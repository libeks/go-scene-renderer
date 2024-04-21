package geometry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test3DInverse(t *testing.T) {
	tests := []struct {
		name        string
		matrix      Matrix3D
		wantInverse Matrix3D
	}{
		{"text matrix",
			Matrix3D{
				1, 2, 0,
				1, 3, 0,
				0, 4, -5,
			},
			Matrix3D{
				3, -2, 0,
				-1, 1, 0,
				-.8, .8, -.2,
			},
		},
		{"text matrix 1",
			Matrix3D{
				1, 2, 3,
				1, 4, 3,
				7, 8, 9,
			},
			Matrix3D{
				-.5, -.25, .25,
				-0.5, 0.5, 0,
				.8333333333333333, -.25, -.08333333333333333,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, ok := tt.matrix.Inverse()
			if !ok {
				t.Errorf("Could not find inverse of %s", tt.matrix)
			}
			t.Logf("Inverse is %s", inv)
			if diff := cmp.Diff(tt.wantInverse, inv); diff != "" {
				t.Errorf("wrong inverse, diff: %s", diff)
			}
			got := inv.MatrixMult(tt.matrix)
			if diff := cmp.Diff(Identity3D, got, approxFloatOpt); diff != "" {
				t.Errorf("did not get identity, diff: %s", diff)
			}
		})
	}
}
