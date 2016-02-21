package main

import (
	"math/rand"

	"github.com/fogleman/gg"
)

func random() float64 {
	return rand.Float64()*2 - 1
}

func point() (x, y float64) {
	return random(), random()
}

func randomQuadratic(dc *gg.Context) {
	x0, y0 := point()
	x1, y1 := point()
	x2, y2 := point()
	dc.MoveTo(x0, y0)
	dc.QuadraticTo(x1, y1, x2, y2)
}

func randomCubic(dc *gg.Context) {
	x0, y0 := point()
	x1, y1 := point()
	x2, y2 := point()
	x3, y3 := point()
	dc.MoveTo(x0, y0)
	dc.CubicTo(x1, y1, x2, y2, x3, y3)
}

func main() {
	const (
		S = 256
		W = 8
		H = 8
	)
	dc := gg.NewContext(S*W, S*H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			x := float64(i)*S + S/2
			y := float64(j)*S + S/2
			dc.Push()
			dc.Translate(x, y)
			dc.Scale(S/2, S/2)
			randomCubic(dc)
			// randomQuadratic(dc)
			dc.Pop()
		}
	}
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(16)
	dc.StrokePreserve()
	dc.SavePNG("out.png")
}
