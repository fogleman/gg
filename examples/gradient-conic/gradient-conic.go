// +build ignore

package main

import (
	"image/color"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(400, 400)

	grad1 := gg.NewConicGradient(200, 200, 0)
	grad1.AddColorStop(0.0, color.Black)
	grad1.AddColorStop(0.5, color.RGBA{255, 215, 0, 255})
	grad1.AddColorStop(1.0, color.RGBA{255, 0, 0, 255})

	grad2 := gg.NewConicGradient(200, 200, 90)
	grad2.AddColorStop(0.00, color.RGBA{255, 0, 0, 255})
	grad2.AddColorStop(0.16, color.RGBA{255, 255, 0, 255})
	grad2.AddColorStop(0.33, color.RGBA{0, 255, 0, 255})
	grad2.AddColorStop(0.50, color.RGBA{0, 255, 255, 255})
	grad2.AddColorStop(0.66, color.RGBA{0, 0, 255, 255})
	grad2.AddColorStop(0.83, color.RGBA{255, 0, 255, 255})
	grad2.AddColorStop(1.00, color.RGBA{255, 0, 0, 255})

	dc.SetStrokeStyle(grad1)
	dc.SetLineWidth(20)
	dc.DrawCircle(200, 200, 180)
	dc.Stroke()

	dc.SetFillStyle(grad2)
	dc.DrawCircle(200, 200, 150)
	dc.Fill()

	dc.SavePNG("gradient-conic.png")
}
