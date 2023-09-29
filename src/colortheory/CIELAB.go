package colortheory

import (
	"math"
)

type CIELAB struct {
	L, a, b float64
}

func DeltaE00(lab1, lab2 CIELAB) float64 {
	// Credits: The implementation is based on the CIEDE2000 algorithm,
	// detailed in the paper "The CIEDE2000 Color-Difference Formula:
	// Implementation Notes, Supplementary Test Data, and Mathematical Observations" by Gaurav Sharma, Wencheng Wu, Edul N. Dalal.

	// Weighting factors for the L, C, H components.
	kL, kC, kH := 1.0, 1.0, 1.0

	// Step 1: Calculate C1, C2, and the average of C1 and C2, CBar.
	// These are the chroma values in the CIELAB color space, representing the colorfulness of the color.
	C1 := math.Sqrt(lab1.a*lab1.a + lab1.b*lab1.b)
	C2 := math.Sqrt(lab2.a*lab2.a + lab2.b*lab2.b)
	CBar := (C1 + C2) / 2.0

	// Step 2: Compute G, a weighting value used to adjust the a components in CIELAB color space.
	G := 0.5 * (1 - math.Sqrt(math.Pow(CBar, 7)/(math.Pow(CBar, 7)+math.Pow(25, 7))))

	// Calculate modified a components, a1Prime and a2Prime.
	a1Prime := lab1.a * (1 + G)
	a2Prime := lab2.a * (1 + G)

	// Step 3: Compute the derived chroma and hue differences.
	C1Prime := math.Sqrt(a1Prime*a1Prime + lab1.b*lab1.b)
	C2Prime := math.Sqrt(a2Prime*a2Prime + lab2.b*lab2.b)

	deltaLPrime := lab2.L - lab1.L
	deltaCPrime := C2Prime - C1Prime

	// Calculate h1Prime and h2Prime, the modified hue components.
	h1Prime := math.Atan2(lab1.b, a1Prime)
	if h1Prime < 0 {
		h1Prime += 2 * math.Pi
	}
	h2Prime := math.Atan2(lab2.b, a2Prime)
	if h2Prime < 0 {
		h2Prime += 2 * math.Pi
	}

	// Compute the hue difference, deltahPrime.
	var deltahPrime float64
	if math.Abs(h1Prime-h2Prime) <= math.Pi {
		deltahPrime = h2Prime - h1Prime
	} else if h2Prime <= h1Prime {
		deltahPrime = h2Prime - h1Prime + 2*math.Pi
	} else {
		deltahPrime = h2Prime - h1Prime - 2*math.Pi
	}

	// Compute the delta H' value.
	deltaHPrime := 2 * math.Sqrt(C1Prime*C2Prime) * math.Sin(deltahPrime/2)

	// Step 4: Calculate average L, CPrime, and HPrime for later steps.
	LBar := (lab1.L + lab2.L) / 2.0
	CPrimeBar := (C1Prime + C2Prime) / 2.0
	HPrimeBar := (h1Prime + h2Prime) / 2.0
	if math.Abs(h1Prime-h2Prime) > math.Pi {
		HPrimeBar += math.Pi
	}

	// Step 5: Compute the parameters for the final formula.
	// These include SL, SC, SH weightings, T, a correction factor based on hue, and RT, a rotation function.
	T := 1 - 0.17*math.Cos(HPrimeBar-math.Pi/6) + 0.24*math.Cos(2*HPrimeBar) + 0.32*math.Cos(3*HPrimeBar+math.Pi/30) - 0.20*math.Cos(4*HPrimeBar-21*math.Pi/60)
	deltaTheta := 30 * math.Exp(-((HPrimeBar-30*math.Pi/180)*(HPrimeBar-30*math.Pi/180))/25)
	RC := 2 * math.Sqrt(math.Pow(CPrimeBar, 7)/(math.Pow(CPrimeBar, 7)+math.Pow(25, 7)))
	SL := 1 + (0.015*(LBar-50)*(LBar-50))/math.Sqrt(20+(LBar-50)*(LBar-50))
	SC := 1 + 0.045*CPrimeBar
	SH := 1 + 0.015*CPrimeBar*T
	RT := -RC * math.Sin(2*deltaTheta*math.Pi/180)

	// Step 6: Compute the CIEDE2000 color difference value.
	deltaE := math.Sqrt(math.Pow(deltaLPrime/(kL*SL), 2) + math.Pow(deltaCPrime/(kC*SC), 2) + math.Pow(deltaHPrime/(kH*SH), 2) + RT*(deltaCPrime/(kC*SC))*(deltaHPrime/(kH*SH)))

	return deltaE
}

func labToXYZ(lab CIELAB, whitePoint CIEXYZ) CIEXYZ {
	// Step 1: Compute fy, fx, and fz.
	var fy = (lab.L + 16.0) / 116.0
	var fx = fy + lab.a/500.0
	var fz = fy - lab.b/200.0

	// Step 2: Compute XYZ.
	// Applying cube transformation for each of X, Y, Z
	// If (fx^3 > 0.008856) then X = whitePoint.X * fx^3 else X = whitePoint.X * (fx - 16/116) * 3 * ((116)^2)
	// and similarly for Y and Z.
	var xyz CIEXYZ
	xyz.Y = whitePoint.Y * func(f float64) float64 {
		if pow := math.Pow(f, 3); pow > 0.008856 {
			return pow
		}
		return (f - 16.0/116.0) * (116.0 * 116.0) / 81.0
	}(fy)
	xyz.X = whitePoint.X * func(f float64) float64 {
		if pow := math.Pow(f, 3.0); pow > 0.008856 {
			return pow
		}
		return (f - 16.0/116.0) * (116.0 * 116.0) / 81.0
	}(fx)
	xyz.Z = whitePoint.Z * func(f float64) float64 {
		if pow := math.Pow(f, 3.0); pow > 0.008856 {
			return pow
		}
		return (f - 16.0/116.0) * (116 * 116) / 81
	}(fz)

	return xyz
}
