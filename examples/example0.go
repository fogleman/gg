package main

import "github.com/fogleman/gg"

func main() {
	dc := gg.NewContext(256, 256)
	dc.SetSourceRGBA(1, 0, 0, 0.3)
	dc.Paint()
	dc.MoveTo(20, 20)
	dc.LineTo(236, 236)
	dc.LineTo(236, 128)
	dc.LineTo(20, 128)
	dc.QuadraticTo(0, 64, 120, 20)
	dc.SetSourceRGBA(1, 0, 0, 0.8)
	dc.FillPreserve()
	dc.SetSourceRGB(0, 0, 0)
	dc.SetLineWidth(8)
	dc.Stroke()
	dc.WriteToPNG("out.png")
}
