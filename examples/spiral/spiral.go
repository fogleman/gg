package main

import (
	"math"

	"github.com/fogleman/gg"
)

func main() {
	const S = 1024
	const N = 2048
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	for i := 0; i <= N; i++ {
		t := float64(i) / N
		d := t*S*0.4 + 10
		a := t * math.Pi * 2 * 20
		x := S/2 + math.Cos(a)*d
		y := S/2 + math.Sin(a)*d
		r := t * 8
		dc.DrawCircle(x, y, r)
	}
	dc.Fill()
	dc.SavePNG("out.png")
}
