package calc

import (
	"math"
	"sort"
)

func RoundFloat64(value float64, precision int) float64 {
	n := math.Pow(10, float64(precision))
	return math.Round(value*n) / n * 1.0
}

func ClampFloat64(value float64, min_ float64, max_ float64) float64 {
	valueList := []float64{min_, max_, value}
	sort.Float64s(valueList)
	return valueList[1]
}
