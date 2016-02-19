package main

import "github.com/fogleman/dd"

func main() {
	dc := dd.NewContext(1000, 1000)
	dc.SetSourceRGB(1, 1, 1)
	dc.Paint()
	dc.DrawCircle(500, 500, 400)
	dc.SetSourceRGBA(0, 0, 0, 0.25)
	dc.FillPreserve()
	dc.SetSourceRGB(0, 0, 0.5)
	dc.SetLineWidth(8)
	dc.Stroke()
	dc.WriteToPNG("out.png")
}
