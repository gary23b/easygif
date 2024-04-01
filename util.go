package easygif

import (
	"golang.org/x/exp/constraints"
)

// Linearly interpolate between any two numbers of the same type using a ratio.
// The ratio is capped between 0 and 1
func Lerp[T constraints.Integer | constraints.Float](a, b T, ratio float64) T {
	if ratio > 1 {
		ratio = 1
	} else if ratio < 0 {
		ratio = 0
	}
	aF := float64(a)
	bF := float64(b)
	y := aF + (bF-aF)*ratio
	return T(y)
}
