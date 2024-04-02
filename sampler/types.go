package sampler

type Sampler interface {
	GetFrameValue(x, y, t float64) float64
}
