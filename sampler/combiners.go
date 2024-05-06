package sampler

import "slices"

func MaxCombiner(s ...DynamicSampler) DynamicSampler {
	return dynamicCombiner{samplers: s}
}

// implements both Sampler and DynamicSampler
type dynamicCombiner struct {
	samplers []DynamicSampler
}

func (s dynamicCombiner) GetFrame(t float64) StaticSampler {
	samplers := make([]StaticSampler, len(s.samplers))
	for i, sampler := range s.samplers {
		samplers[i] = sampler.GetFrame(t)
	}
	return staticCombiner{
		samplers,
	}
}

func (s dynamicCombiner) GetFrameValue(b, c, t float64) float64 {
	return s.GetFrame(t).GetValue(b, c)
}

func MaxStaticCombiner(s ...StaticSampler) StaticSampler {
	return staticCombiner{
		samplers: s,
	}
}

type staticCombiner struct {
	samplers []StaticSampler
}

func (s staticCombiner) GetValue(b, c float64) float64 {
	vals := make([]float64, len(s.samplers))
	for i, sampler := range s.samplers {
		vals[i] = sampler.GetValue(b, c)
	}
	return slices.Max(vals)
}
