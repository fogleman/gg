package main

import (
	"image/color"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(400, 400)

	grad := gg.NewConicGradient(200, 200, 90)
	grad.AddColorStop(0.0, color.Black)
	grad.AddColorStop(0.5, color.RGBA{255, 215, 0, 255})
	grad.AddColorStop(1.0, color.RGBA{255, 0, 0, 255})

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetStrokeStyle(grad)
	dc.SetLineWidth(15)
	dc.DrawCircle(200, 200, 180)
	dc.Stroke()

	dc.SetFillStyle(grad)
	dc.DrawCircle(200, 200, 155)
	dc.Fill()

	dc.SavePNG("out.png")
}
