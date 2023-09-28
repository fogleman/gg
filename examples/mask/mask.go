package main

import (
	"log"

	"github.com/wildberries-ru/gg"
	"golang.org/x/image/draw"
)

func main() {
	im, err := gg.LoadImage("examples/baboon.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(512, 512)
	dc.DrawRoundedRectangle(0, 0, 512, 512, 64)
	dc.Clip()
	dc.DrawImage(im, 0, 0, draw.BiLinear)
	_ = dc.SavePNG("out.png")
}
