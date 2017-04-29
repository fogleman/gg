package main

import "github.com/fogleman/gg"

func main() {
	const S = 4096 * 2
	const T = 16 * 2
	const F = 28
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("Xolonium-Regular.ttf", F); err != nil {
		panic(err)
	}
	for r := 0; r < 256; r++ {
		for c := 0; c < 256; c++ {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2
			dc.DrawStringAnchored(string(rune(i)), x, y, 0.5, 0.5)
		}
	}
	dc.SavePNG("out.png")
}
