package main

import (
	"fmt"
	"html/template"
	"net/http"
	ct "prisma/colortheory"
)

func main() {
	ct.TestLAB()
	return

	baseColor := ct.NewColorProfileFromRGB(200, 33, 217)

	// Ayu Mirage color scheme
	ayuMirageColors := NewTerminalPalette(
		ct.NewColorProfileFromRGB(20, 21, 25),    //Black
		ct.NewColorProfileFromRGB(255, 85, 85),   //Red
		ct.NewColorProfileFromRGB(136, 192, 208), //Green
		ct.NewColorProfileFromRGB(250, 189, 47),  //Yellow
		ct.NewColorProfileFromRGB(108, 135, 255), //Blue
		ct.NewColorProfileFromRGB(220, 140, 255), //Magenta
		ct.NewColorProfileFromRGB(167, 212, 159), //Cyan
		ct.NewColorProfileFromRGB(173, 178, 183), //White
		ct.NewColorProfileFromRGB(99, 100, 101),  //Grey
		ct.NewColorProfileFromRGB(255, 85, 85),   //LightRed
		ct.NewColorProfileFromRGB(162, 210, 142), //LightGreen
		ct.NewColorProfileFromRGB(219, 186, 115), //Brown
		ct.NewColorProfileFromRGB(79, 121, 200),  //LightBlue
		ct.NewColorProfileFromRGB(197, 120, 228), //LightMagenta
		ct.NewColorProfileFromRGB(141, 210, 138), //LightCyan
		ct.NewColorProfileFromRGB(255, 255, 255), //White
		"Ayu Mirage - ",
	)

	triads := ct.GetHarmonics(&baseColor, 3)

	tetrads := ct.GetHarmonics(&baseColor, 4)

	harmonics := ct.GetHarmonics(&baseColor, 20)

	analogous := ct.GetAnalogous(&baseColor, 9, 50)

	monochromatics := ct.GetMonochromatic(&baseColor, 20, 180)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.ColorProfile{baseColor}, "Base Color", "Base "),
		makeColorGroup(ayuMirageColors.asSlice(), "Ayu", "ubu"),
		makeColorGroup(triads, "Triads - 3 colors evenly distributed", "Triad"),
		makeColorGroup(tetrads, "Tetrads - 4 colors evenly distributed", "Tetrad"),
		makeColorGroup(harmonics, "Harmonic Set", "harmonics"),
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

type terminalPalette struct {
	Black,
	Red,
	Green,
	Yellow,
	Blue,
	Magenta,
	Cyan,
	White,
	Grey,
	LightRed,
	LightGreen,
	Brown,
	LightBlue,
	LightMagenta,
	LightCyan,
	HighWhite ct.ColorProfile
}

func NewTerminalPalette(
	Black, Red, Green, Yellow, Blue,
	Magenta, Cyan, White, Grey, LightRed, LightGreen, Brown,
	LightBlue, LightMagenta, LightCyan, HighWhite ct.ColorProfile, prefix string) terminalPalette {

	Black.SetName(prefix + "Black")
	Red.SetName(prefix + "Red")
	Green.SetName(prefix + "Green")
	Yellow.SetName(prefix + "Yellow")
	Blue.SetName(prefix + "Blue")
	Magenta.SetName(prefix + "Magenta")
	Cyan.SetName(prefix + "Cyan")
	White.SetName(prefix + "White")
	Grey.SetName(prefix + "Grey")
	LightRed.SetName(prefix + "LightRed")
	LightGreen.SetName(prefix + "LightGreen")
	Brown.SetName(prefix + "Brown")
	LightBlue.SetName(prefix + "LightBlue")
	LightMagenta.SetName(prefix + "LightMagenta")
	LightCyan.SetName(prefix + "LightCyan")
	HighWhite.SetName(prefix + "HighWhite")

	return terminalPalette{
		Black,
		Red,
		Green,
		Yellow,
		Blue,
		Magenta,
		Cyan,
		White,
		Grey,
		LightRed,
		LightGreen,
		Brown,
		LightBlue,
		LightMagenta,
		LightCyan,
		HighWhite,
	}
}

func (palette *terminalPalette) asSlice() []ct.ColorProfile {
	return []ct.ColorProfile{
		palette.Black,
		palette.Red,
		palette.Green,
		palette.Yellow,
		palette.Blue,
		palette.Magenta,
		palette.Cyan,
		palette.White,
		palette.Grey,
		palette.LightRed,
		palette.LightGreen,
		palette.Brown,
		palette.LightBlue,
		palette.LightMagenta,
		palette.LightCyan,
		palette.HighWhite,
	}
}

type displayColor struct {
	Color       ct.ColorProfile
	Description string
	HexStr      string
	HslStr      string
	RgbStr      string
}

type displayColorGroup struct {
	Colors     []displayColor
	GroupTitle string
}

func makeColorGroup(colSlice []ct.ColorProfile, groupTitle string, individualIdentifier string) displayColorGroup {
	colors := []displayColor{}

	for i, color := range colSlice {
		description := fmt.Sprintf("%s[%d]", individualIdentifier, i)
		if color.Name != "" {
			description = fmt.Sprintf("%s ", color.Name)
		}
		colors = append(colors,
			displayColor{
				Color:       color,
				Description: description,
				HexStr:      color.RGB.AsHEXSTR(),
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
