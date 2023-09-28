package main

import (
	"math"

	"github.com/wildberries-ru/gg"
	"golang.org/x/image/draw"
)

func main() {
	im1, err := gg.LoadPNG("examples/baboon.png")
	if err != nil {
		panic(err)
	}

	im2, err := gg.LoadPNG("examples/gopher.png")
	if err != nil {
		panic(err)
	}

	s1 := im1.Bounds().Size()
	s2 := im2.Bounds().Size()

	width := int(math.Max(float64(s1.X), float64(s2.X)))
	height := s1.Y + s2.Y

	dc := gg.NewContext(width, height)
	dc.DrawImage(im1, 0, 0, draw.BiLinear)
	dc.DrawImage(im2, 0, s1.Y, draw.BiLinear)
	_ = dc.SavePNG("out.png")
}
