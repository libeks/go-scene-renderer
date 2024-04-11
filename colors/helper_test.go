package colors

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBucketRemainder(t *testing.T) {
	tests := []struct {
		name    string
		x       float64
		d       float64
		wantDiv float64
		wantRem float64
	}{
		{"t1", .7, .3, 0.6, 1 / float64(3)},
		{"t2", .7, .15, 0.6, 2 / float64(3)},
		{"t3", .6, .4, 0.4, 0.5},
	}
	opt := cmp.Comparer(func(x, y float64) bool {
		delta := math.Abs(x - y)
		mean := math.Abs(x+y) / 2.0
		return delta/mean < 0.00001
	})

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			div, rem := bucketRemainder(tt.x, tt.d)
			if !cmp.Equal(div, tt.wantDiv, opt) || !cmp.Equal(rem, tt.wantRem, opt) {
				t.Errorf("wanted (%.9f, %.9f), got (%.9f, %.9f)", tt.wantDiv, tt.wantRem, div, rem)
			}
		})
	}
}
