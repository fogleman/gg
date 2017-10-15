package main

import (
	"fmt"

	"github.com/fogleman/gg"
)

func main() {

	// plot
	p := gg.NewPlotXY()
	p.LegNrow = 2
	p.LegAtBottom = false

	// x-values [-1, +1]
	npts := 21
	x := make([]float64, npts)
	for i := 0; i < npts; i++ {
		x[i] = -1 + 2*float64(i)/float64(npts-1)
	}

	// curve 1
	y1 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y1[i] = x[i]
	}
	c1 := p.AddCurve("y = x", x, y1)
	c1.M = "o"
	c1.C = "#000"
	c1.Mec = "#f00"

	// curve 2
	y2 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y2[i] = 2 * x[i]
	}
	c2 := p.AddCurve("y = 2x", x, y2)
	c2.M = "o"
	c2.Void = true

	// curve 3
	y3 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y3[i] = x[i] * x[i]
	}
	c3 := p.AddCurve("y = x²", x, y3)
	c3.M = "s"

	// curve 4
	y4 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y4[i] = 2 * x[i] * x[i]
	}
	c4 := p.AddCurve("y = 2x²", x, y4)
	c4.M = "s"
	c4.Void = true

	// curve 5
	y5 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y5[i] = -x[i]
	}
	c5 := p.AddCurve("y = -x", x, y5)
	c5.M = "+"

	// curve 6
	y6 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		y6[i] = -2 * x[i]
	}
	c6 := p.AddCurve("y = -2x", x, y6)
	c6.M = "x"

	// curve 7
	y7 := make([]float64, npts)
	c7 := p.AddCurve("gopher", y7, x)
	c7.M = "img:examples/gopher30.png"
	c7.Me = 5
	c7.Ls = "none"

	// render graph
	height := 400
	if p.LegAtBottom {
		height = 500
	}
	dc := gg.NewContext(500, height)
	p.Render(dc)

	// save
	dc.SavePNG("/tmp/figure.png")
	fmt.Printf("file </tmp/figure.png> written\n")
}
