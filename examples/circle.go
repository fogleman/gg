package main

import "github.com/fogleman/gg"

func main() {
	dc := gg.NewContext(1000, 1000)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.DrawCircle(500, 500, 400)
	dc.SetRGBA(0, 0, 0, 0.25)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0.5)
	dc.SetLineWidth(8)
	dc.Stroke()
	dc.WritePNG("out.png")
}
