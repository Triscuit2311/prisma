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

func main() {

	baseColor := ct.RGB{200, 33, 217}

	triads := getHarmonics(baseColor.ToHSL(), 3)
	tetrads := getHarmonics(baseColor.ToHSL(), 4)

	harmonics := getHarmonics(baseColor.ToHSL(), 8)
	analogous := getAnalogous(baseColor.ToHSL(), 9, 100)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.HSL{baseColor.ToHSL()}, "Base Color", "Base "),
		makeColorGroup(triads, "Triads - 3 colors evenly distributed", "Triad"),
		makeColorGroup(tetrads, "Tetrads - 4 colors evenly distributed", "Tetrad"),
		makeColorGroup(harmonics, "Harmonic Set - 8 Steps", "8th_Harmonics"),
		makeColorGroup(analogous, "Analogous Set - 20 Degrees 8 shades", "8th_20_analogous"),
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			PageTitle   string
			ColorGroups []displayColorGroup
		}{
			PageTitle:   "Your Page Title",
			ColorGroups: allColorGroups,
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":9000", nil)
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
