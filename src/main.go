package main

import (
	"fmt"
	"html/template"
	"net/http"
	ct "prisma/colortheory"
)

func main() {
	baseColor := ct.NewColorProfileFromRGB(200, 33, 217)

	// Gruvbox color scheme
	gruvboxColors := NewTerminalPalette(
		ct.NewColorProfileFromRGB(40, 40, 40),    //Black
		ct.NewColorProfileFromRGB(204, 36, 29),   //LightRed
		ct.NewColorProfileFromRGB(152, 151, 26),  //LightGreen
		ct.NewColorProfileFromRGB(215, 153, 33),  //Yellow
		ct.NewColorProfileFromRGB(69, 133, 136),  //LightBlue
		ct.NewColorProfileFromRGB(177, 98, 134),  //LightMagenta
		ct.NewColorProfileFromRGB(104, 157, 106), //LightCyan
		ct.NewColorProfileFromRGB(168, 153, 132), //HighWhite
		ct.NewColorProfileFromRGB(146, 131, 116), //Grey
		ct.NewColorProfileFromRGB(251, 73, 52),   //Red
		ct.NewColorProfileFromRGB(184, 187, 38),  //Green
		ct.NewColorProfileFromRGB(250, 189, 47),  //Brown
		ct.NewColorProfileFromRGB(131, 165, 152), //Blue
		ct.NewColorProfileFromRGB(211, 134, 155), //Magenta
		ct.NewColorProfileFromRGB(142, 192, 124), //Cyan
		ct.NewColorProfileFromRGB(235, 219, 178), //White
		"Gruvbox - ",
	)

	// Solarized Dark color scheme
	solarizedDarkColors := NewTerminalPalette(
		ct.NewColorProfileFromRGB(7, 54, 66),     //Black
		ct.NewColorProfileFromRGB(220, 50, 47),   //LightRed
		ct.NewColorProfileFromRGB(133, 153, 0),   //LightGreen
		ct.NewColorProfileFromRGB(181, 137, 0),   //Yellow
		ct.NewColorProfileFromRGB(38, 139, 210),  //LightBlue
		ct.NewColorProfileFromRGB(211, 54, 130),  //LightMagenta
		ct.NewColorProfileFromRGB(42, 161, 152),  //LightCyan
		ct.NewColorProfileFromRGB(238, 232, 213), //HighWhite
		ct.NewColorProfileFromRGB(0, 43, 54),     //Grey
		ct.NewColorProfileFromRGB(203, 75, 22),   //Red
		ct.NewColorProfileFromRGB(88, 110, 117),  //Green
		ct.NewColorProfileFromRGB(101, 123, 131), //Brown
		ct.NewColorProfileFromRGB(131, 148, 150), //Blue
		ct.NewColorProfileFromRGB(108, 113, 196), //Magenta
		ct.NewColorProfileFromRGB(147, 161, 161), //Cyan
		ct.NewColorProfileFromRGB(253, 246, 227), //White
		"(Solarized Dark) ",
	)

	// Closest absolute
	//fmt.Println(ct.GetClosestColor(&red, ct.GetHarmonics(&baseColor, 20)).RGB.AsHEXSTR())

	diff := ct.TotalDeviance(&gruvboxColors.Magenta.RGB, &baseColor.RGB)

	// Closest SD
	fmt.Println(ct.GetClosestColorRelative(diff, 0.5, &gruvboxColors.Red, ct.GetHarmonics(&baseColor, 20)).RGB.AsHEXSTR())

	triads := ct.GetHarmonics(&baseColor, 3)

	tetrads := ct.GetHarmonics(&baseColor, 4)

	harmonics := ct.GetHarmonics(&baseColor, 20)

	analogous := ct.GetAnalogous(&baseColor, 9, 50)

	monochromatics := ct.GetMonochromatic(&baseColor, 20, 180)

	allColorGroups := []displayColorGroup{
		makeColorGroup([]ct.ColorProfile{baseColor}, "Base Color", "Base "),
		makeColorGroup(gruvboxColors.asSlice(), "Gruvbox", "gruv"),
		makeColorGroup(solarizedDarkColors.asSlice(), "Solarized Dark", "solar"),
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
	LightRed,
	LightGreen,
	Yellow,
	LightBlue,
	LightMagenta,
	LightCyan,
	HighWhite,
	Grey,
	Red,
	Green,
	Brown,
	Blue,
	Magenta,
	Cyan,
	White ct.ColorProfile
}

func NewTerminalPalette(
	Black, LightRed, LightGreen, Yellow, LightBlue,
	LightMagenta, LightCyan, HighWhite, Grey, Red, Green, Brown,
	Blue, Magenta, Cyan, White ct.ColorProfile, prefix string) terminalPalette {

	Black.SetName(prefix + "Black")
	LightRed.SetName(prefix + "LightRed")
	LightGreen.SetName(prefix + "LightGreen")
	Yellow.SetName(prefix + "Yellow")
	LightBlue.SetName(prefix + "LightBlue")
	LightMagenta.SetName(prefix + "LightMagenta")
	LightCyan.SetName(prefix + "LightCyan")
	HighWhite.SetName(prefix + "HighWhite")
	Grey.SetName(prefix + "Grey")
	Red.SetName(prefix + "Red")
	Green.SetName(prefix + "Green")
	Brown.SetName(prefix + "Brown")
	Blue.SetName(prefix + "Blue")
	Magenta.SetName(prefix + "Magenta")
	Cyan.SetName(prefix + "Cyan")
	White.SetName(prefix + "White")

	return terminalPalette{
		Black,
		LightRed,
		LightGreen,
		Yellow,
		LightBlue,
		LightMagenta,
		LightCyan,
		HighWhite,
		Grey,
		Red,
		Green,
		Brown,
		Blue,
		Magenta,
		Cyan,
		White,
	}
}

func (palette *terminalPalette) asSlice() []ct.ColorProfile {
	return []ct.ColorProfile{
		palette.Black,
		palette.LightRed,
		palette.LightGreen,
		palette.Yellow,
		palette.LightBlue,
		palette.LightMagenta,
		palette.LightCyan,
		palette.HighWhite,
		palette.Grey,
		palette.Red,
		palette.Green,
		palette.Brown,
		palette.Blue,
		palette.Magenta,
		palette.Cyan,
		palette.White,
	}
}

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
		description := fmt.Sprintf("%s[%d]", individualIdentifier, i)
		if color.Name != "" {
			description = fmt.Sprintf("%s ", color.Name)
		}
		colors = append(colors,
			displayColor{
				Color:       color,
				Description: description,
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
