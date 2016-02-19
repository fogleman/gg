package main

import "github.com/fogleman/gg"

func main() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	for i := 0; i < 360; i += 15 {
		dc.Identity()
		dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.SetRGBA(0, 0, 0, 0.1)
		dc.Fill()
	}
	dc.WritePNG("out.png")
}
