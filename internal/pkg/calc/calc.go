package calc

import (
	"math"

	"golang.org/x/exp/constraints"
)

func RoundFloat64(value float64, precision int) float64 {
	n := math.Pow(10, float64(precision))
	return math.Round(value*n) / n * 1.0
}

func Clamp[T constraints.Ordered](value, min, max T) T {
	if value > max {
		return max
	} else if value < min {
		return min
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

	return RoundFloat64(Clamp(scale, scaleMin, scaleMax), 2)
}
