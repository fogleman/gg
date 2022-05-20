package main

import (
	"math/rand"

	"github.com/fogleman/gg"
)

func CreatePoints(n int) []gg.Point {
	points := make([]gg.Point, n)
	for i := 0; i < n; i++ {
		x := 0.5 + rand.NormFloat64()*0.1
		y := x + rand.NormFloat64()*0.1
		points[i] = gg.Point{x, y}
	}
	return points
}

func main() {
	const S = 1024
	const P = 64
	dc := gg.NewContext(S, S)
	dc.InvertY()
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	points := CreatePoints(1000)
	dc.Translate(P, P)
	dc.Scale(S-P*2, S-P*2)
	// draw minor grid
	for i := 1; i <= 10; i++ {
		x := float64(i) / 10
		dc.MoveTo(x, 0)
		dc.LineTo(x, 1)
		dc.MoveTo(0, x)
		dc.LineTo(1, x)
	}
	dc.SetRGBA(0, 0, 0, 0.25)
	dc.SetLineWidth(1)
	dc.Stroke()
	// draw axes
	dc.MoveTo(0, 0)
	dc.LineTo(1, 0)
	dc.MoveTo(0, 0)
	dc.LineTo(0, 1)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()
	// draw points
	dc.SetRGBA(0, 0, 1, 0.5)
	for _, p := range points {
		dc.DrawCircle(p.X, p.Y, 3.0/S)
		dc.Fill()
	}
	// draw text
	dc.Identity()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("/Library/Fonts/Arial Bold.ttf", 24); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Chart Title", S/2, P/2, 0.5, 0.5)
	if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 18); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("X Axis Title", S/2, S-P/2, 0.5, 0.5)
	dc.SavePNG("out.png")
}
