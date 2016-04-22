package main

import "github.com/fogleman/gg"

func main() {
	const S = 1000
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Translate(S/2, S/2)
	dc.Scale(40, 40)

	var x0, y0, x1, y1, x2, y2, x3, y3, x4, y4 float64
	x0, y0 = -10, 0
	x1, y1 = -5, -10
	x2, y2 = 0, 0
	x3, y3 = 5, 10
	x4, y4 = 10, 0

	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.LineTo(x4, y4)
	dc.SetHexColor("FF2D00")
	dc.SetLineWidth(8)
	dc.Stroke()

	dc.MoveTo(x0, y0)
	dc.QuadraticTo(x1, y1, x2, y2)
	dc.QuadraticTo(x3, y3, x4, y4)
	dc.SetHexColor("3E606F")
	dc.SetLineWidth(16)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.Stroke()

	dc.DrawCircle(x0, y0, 0.5)
	dc.DrawCircle(x1, y1, 0.5)
	dc.DrawCircle(x2, y2, 0.5)
	dc.DrawCircle(x3, y3, 0.5)
	dc.DrawCircle(x4, y4, 0.5)
	dc.SetRGB(1, 1, 1)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()

	dc.LoadFontFace("/Library/Fonts/Arial.ttf", 200)
	dc.DrawStringAnchored("g", -5, 5, 0.5, 0.5)
	dc.DrawStringAnchored("G", 5, -5, 0.5, 0.5)

	dc.SavePNG("out.png")
}
