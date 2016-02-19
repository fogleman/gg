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
	points := Polygon(5, 512, 512, 400)
	indexes := []int{0, 2, 4, 1, 3, 0}
	dc := dd.NewContext(1024, 1024)
	dc.SetSourceRGB(1, 1, 1)
	dc.Paint()
	for _, index := range indexes {
		p := points[index]
		dc.LineTo(p.X, p.Y)
	}
	dc.SetSourceRGBA(1, 0, 0, 0.5)
	dc.SetFillRule(dd.FillRuleEvenOdd)
	dc.FillPreserve()
	dc.SetSourceRGB(0, 0, 0)
	dc.SetLineWidth(8)
	dc.Stroke()
	dc.WriteToPNG("out.png")
}
