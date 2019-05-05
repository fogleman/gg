package main

import "github.com/fogleman/gg"

func main() {
	dc := gg.NewAlphaContext(200, 200)
	dc.DrawCircle(100, 100, 80)
	dc.SetRGB(.75, .75, .75)
	dc.Fill()
	dc.SavePNG("out.png")
}
