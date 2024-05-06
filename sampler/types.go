package sampler

// x,y values go from 0-1
type Sampler interface {
	GetFrameValue(x, y, t float64) float64
}

type StaticSampler interface {
	GetValue(x, y float64) float64
}

type DynamicSampler interface {
	GetFrame(t float64) StaticSampler
}
