package main

import "github.com/fogleman/gg"

func main() {
	p := gg.NewPlotXY()
	p.Title = "my graph"
	c := p.AddCurve("data", []float64{1, 2, 3, 4}, []float64{1, 4, 9, 16})
	c.M = "o"
	dc := gg.NewContext(400, 300)
	p.Render(dc)
	dc.SavePNG("/tmp/figure-simple.png")
}
