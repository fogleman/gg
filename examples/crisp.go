package main

import (
	"github.com/fogleman/gg"
)

func main() {
	const W = 1000
	const H = 1000
	const Minor = 10
	const Major = 100

	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// minor grid
	for x := Minor; x < W; x += Minor {
		fx := float64(x) + 0.5
		dc.DrawLine(fx, 0, fx, H)
	}
	for y := Minor; y < H; y += Minor {
		fy := float64(y) + 0.5
		dc.DrawLine(0, fy, W, fy)
	}
	dc.SetLineWidth(1)
	dc.SetRGBA(0, 0, 0, 0.25)
	dc.Stroke()

	// major grid
	for x := Major; x < W; x += Major {
		fx := float64(x) + 0.5
		dc.DrawLine(fx, 0, fx, H)
	}
	for y := Major; y < H; y += Major {
		fy := float64(y) + 0.5
		dc.DrawLine(0, fy, W, fy)
	}
	dc.SetLineWidth(1)
	dc.SetRGBA(0, 0, 0, 0.5)
	dc.Stroke()

	dc.SavePNG("out.png")
}
