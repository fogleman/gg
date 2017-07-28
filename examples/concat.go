package main

import (
	"math"

	"github.com/fogleman/gg"
)

func main() {
	im1, err := gg.LoadPNG("examples/lenna.png")
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
	dc.DrawImage(im1, 0, 0)
	dc.DrawImage(im2, 0, s1.Y)
	dc.SavePNG("out.png")
}
