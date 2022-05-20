package main

import (
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Point struct {
	X, Y float64
}

func Polygon(n int) []Point {
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = Point{math.Cos(a), math.Sin(a)}
	}
	return result
}

func main() {
	const W = 1200
	const H = 120
	const S = 100
	dc := gg.NewContext(W, H)
	dc.SetHexColor("#FFFFFF")
	dc.Clear()
	n := 5
	points := Polygon(n)
	for x := S / 2; x < W; x += S {
		dc.Push()
		s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), H/2)
		dc.Rotate(rand.Float64() * 2 * math.Pi)
		dc.Scale(s, s)
		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor("#FFCC00")
		dc.StrokePreserve()
		dc.SetHexColor("#FFE43A")
		dc.Fill()
		dc.Pop()
	}
	dc.SavePNG("out.png")
}
