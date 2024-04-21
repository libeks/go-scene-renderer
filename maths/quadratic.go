package maths

import "math"

// return all the roots, if any, smallest first
func QuadraticRoots(a float64, b float64, c float64) []float64 {
	d := b*b - 4*a*c
	if d < 0 {
		return []float64{}
	}
	if d == 0 {
		return []float64{-b / (2 * a)}
	}
	x1 := (-b - math.Sqrt(d)) / (2 * a)
	x2 := (-b + math.Sqrt(d)) / (2 * a)
	return []float64{x1, x2}
}
