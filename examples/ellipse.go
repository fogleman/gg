package main

import "github.com/fogleman/gg"

func main() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGBA(0, 0, 0, 0.1)
	for i := 0; i < 360; i += 15 {
		dc.Push()
		dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}
	if im, err := gg.LoadPNG("examples/gopher.png"); err == nil {
		w := im.Bounds().Size().X
		h := im.Bounds().Size().Y
		dc.DrawImage(im, S/2-w/2, S/2-h/2)
	}
	dc.SavePNG("out.png")
}
