package main

import "github.com/fogleman/gg"

var line = ("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod" +
	"tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam," +
	"quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo" +
	"consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse" +
	"cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat" +
	"non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. ")

func main() {
	const W = 800
	const H = 400
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace("/Library/Fonts/Arial.ttf", 18)
	dc.DrawStringWrapped(line, 20, 20, W-40)
	dc.SavePNG("out.png")
}
