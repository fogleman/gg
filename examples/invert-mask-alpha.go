package main

import "github.com/fogleman/gg"

func main() {
	W := 200.0
	H := 200.0
	dc := gg.NewAlphaContext(int(W), int(H))
	dc.DrawCircle(W/2, H/2, (W+H)/6)
	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, W, H)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	dc.SavePNG("out.png")
}
