package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	ct "prisma/colortheory"
)

type displayColor struct {
	Color       ct.HSL
	Description string
	HexStr      string
	HslStr      string
	HsvStr      string
	RgbStr      string
}

type displayColorGroup struct {
	Colors     []displayColor
	GroupTitle string
}

func clampRGBColVal(v int) uint8 {
	return uint8(min(max((v), 0), 255))
}

func main() {

	baseColor := ct.RGB{R: 200, G: 33, B: 217}

	triads := getHarmonics(baseColor.ToHSL(), 3)
	tetrads := getHarmonics(baseColor.ToHSL(), 4)

	harmonics := getHarmonics(baseColor.ToHSL(), 9)
	analogous := getAnalogous(baseColor.ToHSL(), 9, 50)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.HSL{baseColor.ToHSL()}, "Base Color", "Base "),
		makeColorGroup(triads, "Triads - 3 colors evenly distributed", "Triad"),
		makeColorGroup(tetrads, "Tetrads - 4 colors evenly distributed", "Tetrad"),
		makeColorGroup(harmonics, "Harmonic Set - 9 Steps", "9th_Harmonics"),
		makeColorGroup(analogous, "Analogous Set - 50 Degrees 9 shades", "9th_50_analogous"),
		makeColorGroup(getMonochromatic(baseColor.ToHSL(), 10, 180), "MonoChromatic Test 1 uniform sat/light", "mono1"),
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ColorGroups []displayColorGroup
		}{
			ColorGroups: allColorGroups,
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":9000", nil)
	fmt.Println("http://localhost:9000/")

}

func getHarmonics(hsl ct.HSL, count int) []ct.HSL {

	harmonics := []ct.HSL{hsl}

	percentInc := 1.0 / float64(count)

	for i := 1; i < count; i++ {
		col := hsl
		col.H = math.Mod(col.H+(percentInc*float64(i)), 1.0)
		harmonics = append(harmonics, col)
	}
	return harmonics
}

func getAnalogous(hsl ct.HSL, count, degreesSpread int) []ct.HSL {
	analogous := []ct.HSL{}

	percentInc := (float64(degreesSpread) / float64(count)) / 360.0

	start := hsl.H - (percentInc * (float64(count) / 2.0))

	for i := 0; i < count; i++ {
		col := ct.HSL{
			H: math.Mod(start+(percentInc*float64(i)), 1.0),
			S: hsl.S,
			L: hsl.L,
		}
		analogous = append(analogous, col)
	}
	return analogous
}

// Hue stays the same, saturation and lightness spread
func getMonochromatic(hsl ct.HSL, count, rangePercent int) []ct.HSL {
	colors := []ct.HSL{}
	rgb := hsl.ToRGB()

	scaledRange := int(float64(rangePercent) / 100.0 * 255.0)

	for i := 0 - (count / 2); i < count/2; i++ {
		col := ct.RGB{
			R: clampRGBColVal(int(rgb.R) + int(i*(scaledRange/count))),
			G: clampRGBColVal(int(rgb.G) + int(i*(scaledRange/count))),
			B: clampRGBColVal(int(rgb.B) + int(i*(scaledRange/count))),
		}
		colors = append(colors, col.ToHSL())
	}

	return colors
}

func makeColorGroup(colSlice []ct.HSL, groupTitle string, individualIdentifier string) displayColorGroup {
	colors := []displayColor{}

	for i, hsl := range colSlice {

		rgb := hsl.ToRGB()
		hsv := rgb.ToHSV()

		colors = append(colors,
			displayColor{
				Color:       hsl,
				Description: fmt.Sprintf("%s[%d]", individualIdentifier, i),
				HexStr:      rgb.AsHEXSTR(),
				HsvStr:      hsv.String(),
				HslStr:      hsl.String(),
				RgbStr:      rgb.String(),
			},
		)
	}

	return displayColorGroup{
		Colors:     colors,
		GroupTitle: groupTitle,
	}
}
