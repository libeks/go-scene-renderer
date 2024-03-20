package maths

import "math"

var (
	Rotation = (2 * math.Pi)
)

func Sigmoid(v float64) float64 {
	// takes from (-inf, +int) to (0.0, 1.0), with an S-like shape centered on 0.0.
	return 1 / (1 + math.Exp(-v))
}
