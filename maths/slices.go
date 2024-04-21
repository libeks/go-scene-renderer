package maths

func KeepOnlyPositives(ts []float64) []float64 {
	newList := []float64{}
	for _, t := range ts {
		if t >= 0 {
			newList = append(newList, t)
		}
	}
	return newList
}
