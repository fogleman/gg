package main

import (
	"math"

	"github.com/fogleman/dd"
)

type Point struct {
	X, Y float64
}

func Polygon(n int, x, y, r float64) []Point {
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = Point{x + r*math.Cos(a), y + r*math.Sin(a)}
	}
	return result
}

func main() {
	n := 5
	points := Polygon(n, 512, 512, 400)
	dc := dd.NewContext(1024, 1024)
	dc.SetSourceRGB(1, 1, 1)
	dc.Paint()
	for i := 0; i < n+1; i++ {
		index := (i * 2) % n
		p := points[index]
		dc.LineTo(p.X, p.Y)
	}
	dc.SetSourceRGBA(0, 0.5, 0, 1)
	dc.SetFillRule(dd.FillRuleEvenOdd)
	dc.FillPreserve()
	dc.SetSourceRGBA(0, 1, 0, 0.5)
	dc.SetLineWidth(16)
	dc.Stroke()
	dc.WriteToPNG("out.png")
}
