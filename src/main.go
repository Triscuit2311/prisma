package main

import (
	"fmt"
	"html/template"
	"net/http"
	ct "prisma/colortheory"
)

// func main() {
// 	col := ct.NewColorProfileFromRGB(200, 17, 187)

// 	fmt.Println(col.RGB.String())
// }

func main() {

	// ubuntuColors := terminalPalette{
	// 	Black:        ct.RGB{R: 0, G: 0, B: 0},
	// 	LightRed:     ct.RGB{R: 187, G: 0, B: 0},
	// 	LightGreen:   ct.RGB{R: 0, G: 187, B: 0},
	// 	Yellow:       ct.RGB{R: 187, G: 187, B: 0},
	// 	LightBlue:    ct.RGB{R: 0, G: 0, B: 187},
	// 	LightMagenta: ct.RGB{R: 187, G: 0, B: 187},
	// 	LightCyan:    ct.RGB{R: 0, G: 187, B: 187},
	// 	HighWhite:    ct.RGB{R: 187, G: 187, B: 187},
	// 	Grey:         ct.RGB{R: 85, G: 85, B: 85},
	// 	Red:          ct.RGB{R: 255, G: 85, B: 85},
	// 	Green:        ct.RGB{R: 85, G: 255, B: 85},
	// 	Brown:        ct.RGB{R: 255, G: 255, B: 85},
	// 	Blue:         ct.RGB{R: 85, G: 85, B: 255},
	// 	Magenta:      ct.RGB{R: 255, G: 85, B: 255},
	// 	Cyan:         ct.RGB{R: 85, G: 255, B: 255},
	// 	White:        ct.RGB{R: 255, G: 255, B: 255},
	// }

	baseColor := ct.NewColorProfileFromRGB(200, 33, 217)

	triads := ct.GetHarmonics(&baseColor, 3)

	tetrads := ct.GetHarmonics(&baseColor, 4)

	harmonics := ct.GetHarmonics(&baseColor, 9)

	analogous := ct.GetAnalogous(&baseColor, 9, 50)

	monochromatics := ct.GetMonochromatic(&baseColor, 20, 180)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.ColorProfile{baseColor}, "Base Color", "Base "),
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

// type terminalPalette struct {
// 	Black,
// 	LightRed,
// 	LightGreen,
// 	Yellow,
// 	LightBlue,
// 	LightMagenta,
// 	LightCyan,
// 	HighWhite,
// 	Grey,
// 	Red,
// 	Green,
// 	Brown,
// 	Blue,
// 	Magenta,
// 	Cyan,
// 	White ct.RGB
// }

type displayColor struct {
	Color       ct.ColorProfile
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

func makeColorGroup(colSlice []ct.ColorProfile, groupTitle string, individualIdentifier string) displayColorGroup {
	colors := []displayColor{}

	for i, color := range colSlice {
		colors = append(colors,
			displayColor{
				Color:       color,
				Description: fmt.Sprintf("%s[%d]", individualIdentifier, i),
				HexStr:      color.RGB.AsHEXSTR(),
				HsvStr:      color.HSV.String(),
				HslStr:      color.HSL.String(),
				RgbStr:      color.RGB.String(),
			},
		)
	}

	return displayColorGroup{
		Colors:     colors,
		GroupTitle: groupTitle,
	}
}
