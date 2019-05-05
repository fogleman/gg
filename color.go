package gg

import "image/color"

// gray16Alpha represents a 16-bit greyscale color having an alpha channel.
type gray16Alpha struct {
	Y, A uint16
}

func (c gray16Alpha) RGBA() (r, g, b, a uint32) {
	cy := uint32(c.Y)
	ca := uint32(c.A)
	return cy, cy, cy, ca
}

// Gray16AlphaModel implements the color.Model interface.
var Gray16AlphaModel = color.ModelFunc(func(c color.Color) color.Color {
	if _, ok := c.(gray16Alpha); ok {
		return c // passthrough
	}
	// gray16Model() in $GOPATH/src/github.com/image/color/color.go
	r, g, b, a := c.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16
	return gray16Alpha{Y: uint16(y), A: uint16(a)}
})
