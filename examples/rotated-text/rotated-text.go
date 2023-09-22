package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/wildberries-ru/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	const S = 400
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 40,
	})
	dc.SetFontFace(face)
	text := "Hello, world!"
	w, h := dc.MeasureString(text)
	dc.Rotate(gg.Radians(10))
	dc.DrawRectangle(100, 180, w, h)
	dc.Stroke()
	dc.DrawStringAnchored(text, 100, 180, 0.0, 0.0, draw.BiLinear)
	dc.SavePNG("out.png")
}
