package textures

import (
	"testing"
)

func TestGrid(t *testing.T) {
	tests := []struct {
		name      string
		x         int
		y         int
		n         int
		wantIndex int
	}{
		{"t1", 1, 1, 2, 3},
		{"t2", 18, 4, 20, 364},
		{"t2", 0, 5, 10, 5},
		{"t2", 5, 0, 10, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := GetCoord(tt.x, tt.y, tt.n)
			if idx != tt.wantIndex {
				t.Errorf("wanted %d, got %d", tt.wantIndex, idx)
			}
			x, y := IndexToCoord(idx, tt.n)
			if x != tt.x || y != tt.y {
				t.Errorf("wanted (%d,%d), got (%d,%d)", tt.x, tt.y, x, y)
			}
		})
	}
}
