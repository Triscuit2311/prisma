package colortheory

import (
	"fmt"
	"math"
)

// Math Exts

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

// absInt calculates the absolute difference between two integers a and b.
func absInt(a, b int) int {
	n := a - b
	if n < 0 {
		n = -n
	}
	return n
}

const (
	one_sixth  float64 = 1.0 / 6.0
	one_half   float64 = 1.0 / 2.0
	two_thirds float64 = 2.0 / 3.0
	one_third  float64 = 1.0 / 3.0
)

// cHSL represents a color in Hue, Saturation, and Lightness format. (float64 [0-1])
type cHSL struct {
	H, S, L float64
}

// cHSV represents a color in Hue, Saturation, and Value format. (float64 [0-1])
type cHSV struct {
	H, S, V float64
}

// cRGB represents a color in Red, Green, and Blue format. (uint8 [0-255])
type cRGB struct {
	R, G, B uint8
}

// ColorProfile represents a color with different color representations and an Alpha value.
type ColorProfile struct {
	HSL   cHSL
	HSV   cHSV
	RGB   cRGB
	Alpha float64
}

// NewColorProfileFromRGB creates a new ColorProfile from Red, Green, and Blue values.
func NewColorProfileFromRGB(r, g, b uint8) ColorProfile {
	rgb := cRGB{r, g, b}
	return ColorProfile{
		RGB:   rgb,
		HSL:   rgb.ToHSL(),
		HSV:   rgb.ToHSV(),
		Alpha: 1.0,
	}
}

// NewColorProfileFromHSL creates a new ColorProfile from Hue, Saturation, and Lightness values.
func NewColorProfileFromHSL(h, s, l float64) ColorProfile {
	hsl := cHSL{h, s, l}
	rgb := hsl.ToRGB()
	return ColorProfile{
		RGB:   rgb,
		HSL:   hsl,
		HSV:   rgb.ToHSV(),
		Alpha: 1.0,
	}
}

func newColorProfileFromFullHSL(hsl cHSL) ColorProfile {
	rgb := hsl.ToRGB()
	return ColorProfile{
		RGB:   rgb,
		HSL:   hsl,
		HSV:   rgb.ToHSV(),
		Alpha: 1.0,
	}
}

func newColorProfileFromFullRGB(rgb cRGB) ColorProfile {
	return ColorProfile{
		RGB:   rgb,
		HSL:   rgb.ToHSL(),
		HSV:   rgb.ToHSV(),
		Alpha: 1.0,
	}
}

// Pretty format [0-360°,0-100%,0-100%]
func (hsl *cHSL) String() string {
	return fmt.Sprintf("cHSL(%d°, %d%%, %d%%)",
		positiveMod(int(hsl.H*360), 360),
		positiveMod(int(hsl.S*100), 100),
		positiveMod(int(hsl.L*100), 100))
}

// Pretty format [0-360°,0-100%,0-100%]
func (hsv *cHSV) String() string {
	return fmt.Sprintf("cHSV(%d°, %d%%, %d%%)",
		positiveMod(int(hsv.H*360), 360),
		positiveMod(int(hsv.S*100), 100),
		positiveMod(int(hsv.V*100), 100))
}

// Pretty format [0-255,0-255,0-255]
func (rgb *cRGB) String() string {
	return fmt.Sprintf("cRGB(%d, %d, %d)", rgb.R, rgb.G, rgb.B)
}

// ToHSL converts RGB color representation to HSL.
func (rgb *cRGB) ToHSL() cHSL {

	// Get cRGB as % of max
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
		return cHSL{
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

	return cHSL{H, S, L}
}

// ToHSV converts RGB color representation to HSV.
func (rgb *cRGB) ToHSV() cHSV {

	// Get cRGB as % of max
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
	return cHSV{H, S, V}
}

// AsHEXSTR returns the HEX string representation of an RGB color.
func (rgb cRGB) AsHEXSTR() string {
	return fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B)
}

// AsArray returns the RGB color as an array of uint8.
func (rgb cRGB) AsArray() [3]uint8 {
	return [3]uint8{rgb.R, rgb.G, rgb.B}
}

// RGBfromArray creates an RGB color from an array of uint8.
func RGBfromArray(arr [3]uint8) cRGB {
	return cRGB{arr[0], arr[1], arr[2]}
}

// ToRGB converts HSL color representation to RGB.
func (hsl *cHSL) ToRGB() cRGB {

	var hueTocRGB = func(lightness, chroma, hue float64) float64 {
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
		return cRGB{lightness, lightness, lightness}
	}

	chroma := 0.0

	if hsl.L < 0.5 {
		chroma = hsl.L * (1 + hsl.S)
	} else {
		chroma = hsl.L + hsl.S - hsl.L*hsl.S
	}
	lightness := 2*hsl.L - chroma

	r := hueTocRGB(lightness, chroma, hsl.H+one_third)
	g := hueTocRGB(lightness, chroma, hsl.H)
	b := hueTocRGB(lightness, chroma, hsl.H-one_third)

	return cRGB{uint8(r * 255), uint8(g * 255), uint8(b * 255)}
}

// ToRGB converts HSV color representation to RGB.
func (hsv *cHSV) ToRGB() cRGB {

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

	return cRGB{uint8(R * 255), uint8(G * 255), uint8(B * 255)}
}

// Lighten lightens the HSL color by a specific percent.
func (hsl *cHSL) Lighten(percent int) {
	hsl.L = clampFloat(hsl.L+(0.01*float64(percent)), 0.0, 1.0)
}

// Darken darkens the HSL color by a specific percent.
func (hsl *cHSL) Darken(percent int) {
	hsl.L = clampFloat(hsl.L-(0.01*float64(percent)), 0.0, 1.0)
}

// Saturate increases the saturation of the HSL color by a specific percent.
func (hsl *cHSL) Saturate(percent int) {
	hsl.S = clampFloat(hsl.S+(0.01*float64(percent)), 0.0, 1.0)
}

// Desaturate decreases the saturation of the HSL color by a specific percent.
func (hsl *cHSL) Desaturate(percent int) {
	hsl.S = clampFloat(hsl.S-(0.01*float64(percent)), 0.0, 1.0)
}

// totalDeviance calculates the total deviance between two RGB colors and returns it as a percent.
func totalDeviance(a *cRGB, b *cRGB) float64 {
	rD := float64(absInt(int(a.R), int(b.R)))
	gD := float64(absInt(int(a.B), int(b.B)))
	bD := float64(absInt(int(a.G), int(b.G)))

	return (rD + gD + bD) / 765.0
}

// GetClosestColor finds the closest color from a list to the given color.
func GetClosestColor(col *ColorProfile, list []ColorProfile) {

}

// GetHarmonics generates harmonic colors from a base color.
func GetHarmonics(color *ColorProfile, count int) []ColorProfile {

	harmonics := []ColorProfile{*color}

	percentInc := 1.0 / float64(count)

	for i := 1; i < count; i++ {
		col := color.HSL
		col.H = math.Mod(col.H+(percentInc*float64(i)), 1.0)
		harmonics = append(harmonics, newColorProfileFromFullHSL(col))
	}
	return harmonics
}

// GetAnalogous generates analogous colors from a base color.
func GetAnalogous(color *ColorProfile, count, degreesSpread int) []ColorProfile {

	analogous := []ColorProfile{}

	percentInc := (float64(degreesSpread) / float64(count)) / 360.0

	start := color.HSL.H - (percentInc * (float64(count) / 2.0))

	for i := 0; i < count; i++ {
		col := cHSL{
			H: math.Mod(start+(percentInc*float64(i)), 1.0),
			S: color.HSL.S,
			L: color.HSL.L,
		}
		analogous = append(analogous, newColorProfileFromFullHSL(col))
	}
	return analogous
}

// GetMonochromatic generates monochromatic colors from a base color.
func GetMonochromatic(color *ColorProfile, count, rangePercent int) []ColorProfile {
	colors := []ColorProfile{}
	rgb := color.RGB

	scaledRange := int(float64(rangePercent) / 100.0 * 255.0)

	for i := 0 - (count / 2); i < count/2; i++ {
		col := cRGB{
			R: clampRGBColVal(int(rgb.R) + int(i*(scaledRange/count))),
			G: clampRGBColVal(int(rgb.G) + int(i*(scaledRange/count))),
			B: clampRGBColVal(int(rgb.B) + int(i*(scaledRange/count))),
		}
		colors = append(colors, newColorProfileFromFullRGB(col))
	}

	return colors
}
