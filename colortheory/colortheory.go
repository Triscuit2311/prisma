package colortheory

import (
	"fmt"
	"math"
)

// Contains Hue, Saturation, Lightness (float64 [0-1])
type HSL struct {
	H, S, L float64
}

// Contains Hue, Saturation, Lightness (float64 [0-1])
type HSV struct {
	H, S, V float64
}

// Contains Red, Green, Blue (uint8 [0-255])
type RGB struct {
	R, G, B uint8
}

func positiveMod(a, b int) int {
	return (a%b + b) % b
}

// Pretty format [0-360°,0-100%,0-100%]
func (hsl *HSL) String() string {
	return fmt.Sprintf("HSL(%d°, %d%%, %d%%)",
		positiveMod(int(hsl.H*360), 360),
		positiveMod(int(hsl.S*100), 100),
		positiveMod(int(hsl.L*100), 100))
}

// Pretty format [0-360°,0-100%,0-100%]
func (hsv *HSV) String() string {
	return fmt.Sprintf("HSV(%d°, %d%%, %d%%)",
		positiveMod(int(hsv.H*360), 360),
		positiveMod(int(hsv.S*100), 100),
		positiveMod(int(hsv.V*100), 100))
}

// Pretty format [0-255,0-255,0-255]
func (rgb *RGB) String() string {
	return fmt.Sprintf("RGB(%d, %d, %d)", rgb.R, rgb.G, rgb.B)
}

// RGB -> HSL
func (rgb *RGB) ToHSL() HSL {

	// Get RGB as % of max
	nR := float64(rgb.R) / 255.0
	nG := float64(rgb.G) / 255.0
	nB := float64(rgb.B) / 255.0

	minV := math.Min(nR, math.Min(nG, nB))
	maxV := math.Max(nR, math.Max(nG, nB))
	delta := maxV - minV

	H, S, L := 0.0, 0.0, 0.0

	// Calculate L
	L = (maxV + minV) / 2.0

	// Grayscale, B/W
	if maxV == minV {
		// No Hue/Sat
		return HSL{
			0,
			0,
			L,
		}
	}

	// Calculate S
	if L < 0.5 {
		S = delta / (maxV + minV)
	} else {
		S = delta / (2.0 - maxV - minV)
	}

	// Calculate H
	switch {
	case rgb.R >= rgb.B && rgb.R >= rgb.G:
		if nG < nB {
			delta += 6
		}
		H = (nG - nB) / delta
	case rgb.G >= rgb.R && rgb.G >= rgb.B:
		H = (nB-nR)/delta + 2
	default:
		H = (nR-nG)/delta + 4
	}

	H /= 6

	return HSL{H, S, L}
}

// RGB -> HSV
func (rgb *RGB) ToHSV() HSV {

	// Get RGB as % of max
	nR := float64(rgb.R) / 255.0
	nG := float64(rgb.G) / 255.0
	nB := float64(rgb.B) / 255.0

	minV := math.Min(nR, math.Min(nG, nB))
	maxV := math.Max(nR, math.Max(nG, nB))
	delta := maxV - minV

	H, S, V := maxV, maxV, maxV

	if maxV == 0.0 {
		S = 0
	} else {
		S = delta / maxV
	}

	if maxV == minV {
		H = 0
	} else {
		switch {
		case rgb.R >= rgb.B && rgb.R >= rgb.G:
			if nG < nB {
				delta += 6
			}
			H = (nG - nB) / delta
		case rgb.G >= rgb.R && rgb.G >= rgb.B:
			H = (nB-nR)/delta + 2
		default:
			H = (nR-nG)/delta + 4
		}
		H /= 6
	}
	return HSV{H, S, V}
}

// RGB -> HEX
func (rgb RGB) AsHEXSTR() string {
	return fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B)
}
func (rgb RGB) AsArray() [3]uint8 {
	return [3]uint8{rgb.R, rgb.G, rgb.B}
}
func RGBfromArray(arr [3]uint8) RGB {
	return RGB{arr[0], arr[1], arr[2]}
}

// HSL -> RGB
func (hsl *HSL) ToRGB() RGB {

	const (
		one_sixth  float64 = 1.0 / 6.0
		one_half   float64 = 1.0 / 2.0
		two_thirds float64 = 2.0 / 3.0
		one_third  float64 = 1.0 / 3.0
	)

	var hueToRGB = func(lightness, chroma, hue float64) float64 {
		if hue < 0 {
			hue += 1
		}
		if hue > 1 {
			hue -= 1
		}
		if hue < one_sixth {
			return lightness + (chroma-lightness)*6.0*hue
		}
		if hue < one_half {
			return chroma
		}
		if hue < two_thirds {
			return lightness + (chroma-lightness)*(two_thirds-hue)*6.0
		}
		return lightness
	}

	// Grayscale, B/W
	if hsl.S == 0 {
		lightness := uint8(hsl.L * 255.0)
		return RGB{lightness, lightness, lightness}
	}

	chroma := 0.0

	if hsl.L < 0.5 {
		chroma = hsl.L * (1 + hsl.S)
	} else {
		chroma = hsl.L + hsl.S - hsl.L*hsl.S
	}
	lightness := 2*hsl.L - chroma

	r := hueToRGB(lightness, chroma, hsl.H+one_third)
	g := hueToRGB(lightness, chroma, hsl.H)
	b := hueToRGB(lightness, chroma, hsl.H-one_third)

	return RGB{uint8(r * 255), uint8(g * 255), uint8(b * 255)}
}

// HSV -> RGB
func (hsv *HSV) ToRGB() RGB {

	R, G, B := 0.0, 0.0, 0.0

	sectorNum := math.Floor(hsv.H * 6.0)
	fractionalSector := hsv.H*6.0 - sectorNum
	decSaturation := hsv.V * (1.0 - hsv.S)
	partialDec := hsv.V * (1.0 - fractionalSector*hsv.S)
	largeDec := hsv.V * (1.0 - (1.0-fractionalSector)*hsv.S)

	switch int(sectorNum) % 6 {
	case 0:
		R, G, B = hsv.V, largeDec, decSaturation
	case 1:
		R, G, B = partialDec, hsv.V, decSaturation
	case 2:
		R, G, B = decSaturation, hsv.V, largeDec
	case 3:
		R, G, B = decSaturation, partialDec, hsv.V
	case 4:
		R, G, B = largeDec, decSaturation, hsv.V
	case 5:
		R, G, B = hsv.V, decSaturation, partialDec

	}

	return RGB{uint8(R * 255), uint8(G * 255), uint8(B * 255)}
}
