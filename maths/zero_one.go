package maths

// the higher k, the faster the movement. It should be in the range 6 <= k <= 20
func GetSigmoidSlowFastSlow(k float64) func(float64) float64 {
	return func(t float64) float64 {
		t = 2*t*k - k
		return Sigmoid(t)
	}
}

func SigmoidSlowFastSlow(t float64) float64 {
	const k = 15
	t = 2*t*k - k
	return Sigmoid(t)
}

func BezierSlowFastSlow(t float64) float64 {
	return BezierZeroOne(t)
}
