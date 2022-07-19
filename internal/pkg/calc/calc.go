package calc

import (
	"math"
)

func RoundFloat64(value float64, precision int) float64 {
	n := math.Pow(10, float64(precision))
	return math.Round(value*n) / n * 1.0
}

func ClampFloat64(value float64, min_ float64, max_ float64) float64 {
	if value > max_ {
		return max_
	} else if value < min_ {
		return min_
	}

	return value
}

func StaticTextScale(text string) float64 {
	scaleMin := 1.0
	scaleMax := 1.5
	lengthMin := float64(len("getquad semi")) // 12
	lengthMax := 3.0 * lengthMin

	lengthFactor := (float64(len(text)) - lengthMin) / (lengthMax - lengthMin)
	scale := scaleMax - (lengthFactor * (scaleMax - scaleMin))

	return RoundFloat64(ClampFloat64(scale, scaleMin, scaleMax), 2)
}
