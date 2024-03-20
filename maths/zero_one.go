package maths

func SigmoidSlowFastSlow(t float64) float64 {
	const k = 7
	t = 2*t*k - k
	return Sigmoid(t)
}

func BezierSlowFastSlow(t float64) float64 {
	return BezierZeroOne(t)
}
