package main

import (
	"fmt"
	"html/template"
	"net/http"
	ct "prisma/colortheory"
)

func main() {
	baseColor := ct.NewColorProfileFromRGB(200, 33, 217)

	// ubuntuColors := NewTerminalPalette(
	// 	ct.NewColorProfileFromRGB(0, 0, 0),       //Black
	// 	ct.NewColorProfileFromRGB(187, 0, 0),     //Red
	// 	ct.NewColorProfileFromRGB(0, 187, 0),     //Green
	// 	ct.NewColorProfileFromRGB(187, 187, 0),   //Yellow
	// 	ct.NewColorProfileFromRGB(0, 0, 187),     //Blue
	// 	ct.NewColorProfileFromRGB(187, 0, 187),   //Magenta
	// 	ct.NewColorProfileFromRGB(0, 187, 187),   //Cyan
	// 	ct.NewColorProfileFromRGB(187, 187, 187), //White
	// 	ct.NewColorProfileFromRGB(85, 85, 85),    //Grey
	// 	ct.NewColorProfileFromRGB(255, 85, 85),   //LightRed
	// 	ct.NewColorProfileFromRGB(85, 255, 85),   //LightGreen
	// 	ct.NewColorProfileFromRGB(255, 255, 85),  //Brown
	// 	ct.NewColorProfileFromRGB(85, 85, 255),   //LightBlue
	// 	ct.NewColorProfileFromRGB(255, 85, 255),  //LightMagenta
	// 	ct.NewColorProfileFromRGB(85, 255, 255),  //LightCyan
	// 	ct.NewColorProfileFromRGB(255, 255, 255), //White
	// 	"(Ubuntu) ",
	// )

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

	myPalette := GenerateTerminalPalette(
		baseColor, ayuMirageColors.Magenta, 0.4,
		ayuMirageColors, "MyPalette", 10, 90, 340, 40, 20,
	)

	triads := ct.GetHarmonics(&baseColor, 3)

	tetrads := ct.GetHarmonics(&baseColor, 4)

	harmonics := ct.GetHarmonics(&baseColor, 20)

	analogous := ct.GetAnalogous(&baseColor, 9, 50)

	monochromatics := ct.GetMonochromatic(&baseColor, 20, 180)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.ColorProfile{baseColor}, "Base Color", "Base "),
		makeColorGroup(ayuMirageColors.asSlice(), "Ayu", "ubu"),
		makeColorGroup(myPalette.asSlice(), "My Palette", "mine"),
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

func GenerateTerminalPalette(
	color, referenceColor ct.ColorProfile, weight float64,
	referencePalette terminalPalette, name string,
	monoLightOffsetPercent, monoSatOffsetPercent,
	monoDegrees, analogousDegrees, colorDensity int) terminalPalette {

	diff := ct.TotalDeviance(&color.RGB, &referenceColor.RGB)
	harmonics := ct.GetHarmonics(&color, colorDensity)
	monochromatics := ct.GetMonochromatic(&color, colorDensity, monoDegrees)
	analogousSet := ct.GetAnalogous(&color, colorDensity, analogousDegrees)

	colorLibrary := append(harmonics, monochromatics...)
	colorLibrary = append(colorLibrary, analogousSet...)

	return NewTerminalPalette(
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Black, colorLibrary).Darkened(monoLightOffsetPercent).Desaturated(monoSatOffsetPercent), //Black
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Red, colorLibrary), //Red
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Green, colorLibrary), //Green
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Yellow, colorLibrary), //Yellow
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Blue, colorLibrary), //Blue
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Magenta, colorLibrary), //Magenta
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Cyan, colorLibrary), //Cyan
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.White, colorLibrary).Darkened(monoLightOffsetPercent).Desaturated(monoSatOffsetPercent), //White
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Grey, colorLibrary).Lightened(monoLightOffsetPercent).Desaturated(monoSatOffsetPercent), //Grey
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.LightRed, colorLibrary), //LightRed
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.LightGreen, colorLibrary), //LightGreen
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.Brown, colorLibrary), //Brown
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.LightBlue, colorLibrary), //LightBlue
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.LightMagenta, colorLibrary), //LightMagenta
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.LightCyan, colorLibrary), //LightCyan
		ct.GetClosestColorRelative(diff, weight,
			&referencePalette.HighWhite, colorLibrary).Lightened(monoLightOffsetPercent).Desaturated(monoSatOffsetPercent), //HighWhite
		name+" - ",
	)

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
