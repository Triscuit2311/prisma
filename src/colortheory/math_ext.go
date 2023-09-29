package colortheory

import "math"

// clampFloat constrains a float64 value v between min and max.
func clampFloat(v, min, max float64) float64 {
	return math.Min(math.Max((v), min), max)
}

// clampRGBColVal constrains an int value v between 0 and 255 and returns it as uint8.
func clampRGBColVal(v int) uint8 {
	return uint8(min(max((v), 0), 255))
}

// positiveMod computes the positive modulus of a and b.
func positiveMod(a, b int) int {
	return (a%b + b) % b
}
