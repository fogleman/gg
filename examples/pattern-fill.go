package main

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

func main() {
	im, err := gg.LoadPNG("examples/lenna.png")
	if err != nil {
		panic(err)
	}
	var x0, y0, x1, y1 float64 = 5, 5, 15, 15
	a := math.Atan2(y1-y0, x1-x0)
	d := a * 180 / math.Pi
	fmt.Println(a, d)
	a = (2*math.Pi - a)
	d = a * 180 / math.Pi
	fmt.Println(a, d)
	m := gg.Translate(-x0, -y0).Rotate(a)
	tx, ty := m.TransformPoint(5, 5)
	fmt.Println(tx, ty)
	pattern := gg.NewSurfacePattern(im, gg.RepeatBoth)
	dc := gg.NewContext(600, 600)
	dc.MoveTo(20, 20)
	dc.LineTo(590, 20)
	dc.LineTo(590, 590)
	dc.LineTo(20, 590)
	dc.ClosePath()
	dc.SetFillStyle(pattern)
	dc.Fill()
	dc.SavePNG("out.png")
}
