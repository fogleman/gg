package main

import "github.com/fogleman/gg"

func main() {
	const W = 400
	const H = 200
	im, err := gg.LoadPNG("examples/gopher.png")
	if err != nil {
		panic(err)
	}
	dc := gg.NewContext(W, H)
	// draw outline
	dc.SetHexColor("#ff0000")
	dc.SetLineWidth(1)
	dc.DrawRectangle(0, 0, float64(W), float64(H))
	dc.Stroke()
	// draw image with current matrix applied
	dc.SetHexColor("#0000ff")
	dc.SetLineWidth(2)
	dc.Rotate(gg.Radians(10))
	dc.DrawRectangle(100, 0, float64(im.Bounds().Dx()), float64(im.Bounds().Dy())/2)
	dc.StrokePreserve()
	dc.Clip()
	dc.DrawImage(im, 100, 0)
	dc.SavePNG("out.png")
}
