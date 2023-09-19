package main

import (
	"fmt"
	"math"
)

type HSL struct {
	H, S, L int
}

func (hsl *HSL) String() string {
	return fmt.Sprintf("HSL(%d, %d%%, %d%%)", hsl.H, hsl.S, hsl.L)
}

func (hsl *HSL) HSLtoRGB() RGB {
	// TODO: Implement
	return RGB{0, 0, 0}
}

type HSV struct {
	H, S, V int
}

func (hsl *HSV) String() string {
	return fmt.Sprintf("HSV(%d, %d%%, %d%%)", hsl.H, hsl.S, hsl.V)
}

func (hsv *HSV) HSVtoRGB() RGB {
	// TODO: Implement
	return RGB{0, 0, 0}
}

type RGB struct {
	R, G, B uint8
}

func (rgb *RGB) RGBtoHSL() HSL {
	nR := float64(rgb.R) / 255.0
	nG := float64(rgb.G) / 255.0
	nB := float64(rgb.B) / 255.0

	minV := math.Min(nR, math.Min(nG, nB))
	maxV := math.Max(nR, math.Max(nG, nB))

	// Covers all 3 channels being the same (gray)
	if maxV == minV {
		// No Hue/Sat
		return HSL{
			0,
			0,
			int((maxV + minV) * 50),
		}
	}

	H, S, L := 0.0, 0.0, (maxV+minV)/2

	// Calculate L
	if L < 0.5 {
		S = (maxV - minV) / (maxV + minV)
	} else {
		S = (maxV - minV) / (2.0 - maxV - minV)
	}

	// Calculate H
	if nR == maxV {
		H = 60 * ((nG - nB) / (maxV - minV))

		if nG < nB {
			H += 360
		}
	} else if nG == maxV {
		H = 60*((nB-nR)/(maxV-minV)) + 120
	} else {
		H = 60*((nR-nG)/(maxV-minV)) + 240
	}

	if H < 0 {
		H += 360
	}

	return HSL{int(H), int(S * 100), int(L * 100)}
}

func (rgb *RGB) RGBtoHSV() HSV {

	nR := float64(rgb.R) / 255.0
	nG := float64(rgb.G) / 255.0
	nB := float64(rgb.B) / 255.0

	minV := math.Min(nR, math.Min(nG, nB))
	maxV := math.Max(nR, math.Max(nG, nB))

	// V is MaxVal normalized
	H, S, V := 0.0, 0.0, maxV

	// Calculate S
	if maxV != 0.0 {
		S = (maxV - minV) / maxV
	}

	// Calculate H
	if nR == maxV {
		H = 60 * ((nG - nB) / (maxV - minV))

		if nG < nB {
			H += 360
		}
	} else if nG == maxV {
		H = 60*((nB-nR)/(maxV-minV)) + 120
	} else {
		H = 60*((nR-nG)/(maxV-minV)) + 240
	}

	if H < 0 {
		H += 360
	}

	return HSV{int(H), int(S * 100), int(V * 100)}
}

func (rgb *RGB) RGBtoHEXSTR() string {
	return fmt.Sprintf("#%x%x%x", rgb.R, rgb.G, rgb.B)
}

func main() {

	c := RGB{200, 33, 217}

	hsl := c.RGBtoHSL()
	hsv := c.RGBtoHSV()

	fmt.Println(hsl)
	fmt.Println(hsv)
	fmt.Println(c.RGBtoHEXSTR())
}
