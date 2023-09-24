package main

import (
	"fmt"
	"html/template"
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

	baseColor := ct.RGB{R: 200, G: 33, B: 217}

	triads := ct.GetHarmonics(baseColor.ToHSL(), 3)
	tetrads := ct.GetHarmonics(baseColor.ToHSL(), 4)

	harmonics := ct.GetHarmonics(baseColor.ToHSL(), 9)

	analogous := ct.GetAnalogous(baseColor.ToHSL(), 9, 50)

	monochromatics := ct.GetMonochromatic(baseColor.ToHSL(), 20, 180)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.HSL{baseColor.ToHSL()}, "Base Color", "Base "),
		makeColorGroup(triads, "Triads - 3 colors evenly distributed", "Triad"),
		makeColorGroup(tetrads, "Tetrads - 4 colors evenly distributed", "Tetrad"),
		makeColorGroup(harmonics, "Harmonic Set - 9 Steps", "9th_Harmonics"),
		makeColorGroup(analogous, "Analogous Set - 50 Degrees 9 shades", "9th_50_analogous"),
		makeColorGroup(monochromatics, "Monochromatic", "monochrome"),
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
	fmt.Println("Serving at: http://localhost:9000/")
	http.ListenAndServe(":9000", nil)

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
