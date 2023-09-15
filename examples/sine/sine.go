package main

import (
	"math"

	"github.com/fogleman/gg"
)

func main() {
	const W = 1200
	const H = 60
	dc := gg.NewContext(W, H)
	// dc.SetHexColor("#FFFFFF")
	// dc.Clear()
	dc.ScaleAbout(0.95, 0.75, W/2, H/2)
	for i := 0; i < W; i++ {
		a := float64(i) * 2 * math.Pi / W * 8
		x := float64(i)
		y := (math.Sin(a) + 1) / 2 * H
		dc.LineTo(x, y)
	}
	dc.ClosePath()
	dc.SetHexColor("#3E606F")
	dc.FillPreserve()
	dc.SetHexColor("#19344180")
	dc.SetLineWidth(8)
	dc.Stroke()
	dc.SavePNG("out.png")
}
