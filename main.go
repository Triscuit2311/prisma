package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	ct "prisma/colortheory"
)

type displayColor struct {
	Color       ct.RGB
	Description string
	Hex         string
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

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.HSL{baseColor.ToHSL()}, "Base Color", "Base "),
		makeColorGroup(triads, "Triads - 3 colors evenly distributed", "Triad"),
		makeColorGroup(tetrads, "Tetrads - 4 colors evenly distributed", "Tetrad"),
		makeColorGroup(harmonics, "Harmonic Set - 8 Steps", "8th_Harmonics"),
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))

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

func colorEffectsTests(rgb ct.RGB) {

	hsl := rgb.ToHSL()
	hsl.L += 0.2 //lighten
	fmt.Println("Lightened 20%: ", hsl.ToRGB().AsHEXSTR())

	hsl = rgb.ToHSL()
	hsl.L -= 0.2 //darken
	fmt.Println("Lightened 20%: ", hsl.ToRGB().AsHEXSTR())

	hsl = rgb.ToHSL()
	hsl.S += 0.2 //saturate
	fmt.Println("Saturated 20%: ", hsl.ToRGB().AsHEXSTR())

	hsl = rgb.ToHSL()
	hsl.S -= 0.2 //desaturate
	fmt.Println("Desaturated 20%: ", hsl.ToRGB().AsHEXSTR())
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

//func getAnalogous()

func makeColorGroup(colSlice []ct.HSL, groupTitle string, individualIdentifier string) displayColorGroup {
	colors := []displayColor{}

	for i, hsl := range colSlice {
		rgb := hsl.ToRGB()
		hsv := rgb.ToHSV()
		colors = append(colors,
			displayColor{
				rgb,
				fmt.Sprintf("%s[%d] | %s | %s | %s\n",
					individualIdentifier, i,
					rgb.String(), hsl.String(), hsv.String()),
				rgb.AsHEXSTR(),
			},
		)
	}

	return displayColorGroup{
		Colors:     colors,
		GroupTitle: groupTitle,
	}
}
