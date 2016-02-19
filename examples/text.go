package main

import (
	"fmt"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(1000, 1000)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace("/Library/Fonts/Arial.ttf", 96)
	s := "Hello, world!"
	w, h := dc.MeasureString(s)
	fmt.Println(w, h)
	dc.DrawString(500-w/2, 500+h/2, s)
	dc.SavePNG("out.png")
}
