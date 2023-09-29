package colortheory

import (
	"math"
)

type CIEXYZ struct {
	X, Y, Z float64
}

// Reference white points
var (
	WhiteD50 = CIEXYZ{X: 96.4212, Y: 100.0000, Z: 82.5188}
	WhiteD55 = CIEXYZ{X: 95.6797, Y: 100.0000, Z: 92.1481}
	WhiteD65 = CIEXYZ{X: 95.0470, Y: 100.0000, Z: 108.8820}
	WhiteD75 = CIEXYZ{X: 94.9722, Y: 100.0000, Z: 122.6394}
	WhiteA   = CIEXYZ{X: 109.8504, Y: 100.0000, Z: 35.5856}
	WhiteB   = CIEXYZ{X: 99.0720, Y: 100.0000, Z: 85.2233}
	WhiteC   = CIEXYZ{X: 98.0720, Y: 100.0000, Z: 118.2320}
	WhiteE   = CIEXYZ{X: 100.0000, Y: 100.0000, Z: 100.0000}
	WhiteF2  = CIEXYZ{X: 99.1866, Y: 100.0000, Z: 67.3930}
	WhiteF7  = CIEXYZ{X: 95.0410, Y: 100.0000, Z: 108.7478}
	WhiteF11 = CIEXYZ{X: 100.9626, Y: 100.0000, Z: 64.3508}
)

func xyzToLab(xyz CIEXYZ, whitePoint CIEXYZ) CIELAB {
	// Step 1: Normalize XYZ values with the reference white point values.
	// The white point is a recognized set of XYZ values that are used to adjust
	// the color balance of the function. Different sources of light (D50, D65 etc.)
	// have different white points.
	x := xyz.X / whitePoint.X
	y := xyz.Y / whitePoint.Y
	z := xyz.Z / whitePoint.Z

	// Step 2: Apply a cube root transformation to the normalized values.
	// This step is intended to linearize the perception of color.
	// It compresses the highlights and accentuates the shadows, which aligns with human perception.
	// If a normalized value is below 0.008856, then it is scaled linearly, otherwise,
	// it's scaled with a cube root function. The constants in the linear scaling, 7.787 and 16/116,
	// are chosen to make the two parts of the function continuous.
	f := func(t float64) float64 {
		if t > 0.008856 {
			return math.Cbrt(t)
		}
		return (7.787 * t) + (16.0 / 116.0)
	}
	fx := f(x)
	fy := f(y)
	fz := f(z)

	// Step 3: Compute L*, a*, and b*.
	// L* represents lightness ranging from 0 to 100.
	// a* represents the color spectrum from green to red.
	// b* represents the color spectrum from blue to yellow.
	// The constants 116 and 500 and 200 are scaling constants to place the output
	// in the proper range. The '-16.0' in the L* calculation centers the scale around zero.
	var lab CIELAB
	lab.L = max(0, (116.0*fy)-16.0)
	lab.a = (fx - fy) * 500.0
	lab.b = (fy - fz) * 200.0

	return lab
}
