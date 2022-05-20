package main

import (
	"image/color"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(400, 200)

	grad := gg.NewRadialGradient(100, 100, 10, 100, 120, 80)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, 200, 200)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawCircle(100, 100, 10)
	dc.Stroke()
	dc.DrawCircle(100, 120, 80)
	dc.Stroke()

	dc.SavePNG("out.png")
}
