package gg

import (
	"image"
	"image/color"

	"github.com/golang/freetype/raster"
)

// m is the maximum color value returned by image.Color.RGBA.
const m = 1<<16 - 1

type RepeatOp int

const (
	RepeatBoth RepeatOp = iota
	RepeatX
	RepeatY
	RepeatNone
)

type Pattern interface {
	ColorAt(x, y int) color.Color
}

// Solid Pattern
type solidPattern struct {
	color color.Color
}

func (p *solidPattern) ColorAt(x, y int) color.Color {
	return p.color
}

func NewSolidPattern(color color.Color) Pattern {
	return &solidPattern{color: color}
}

// Surface Pattern
type surfacePattern struct {
	im image.Image
	op RepeatOp
}

func (p *surfacePattern) ColorAt(x, y int) color.Color {
	b := p.im.Bounds()
	switch p.op {
	case RepeatX:
		if y >= b.Dy() {
			return color.Transparent
		}
	case RepeatY:
		if x >= b.Dx() {
			return color.Transparent
		}
	case RepeatNone:
		if x >= b.Dx() || y >= b.Dy() {
			return color.Transparent
		}
	}
	x = x%b.Dx() + b.Min.X
	y = y%b.Dy() + b.Min.Y
	return p.im.At(x, y)
}

func NewSurfacePattern(im image.Image, op RepeatOp) Pattern {
	return &surfacePattern{im: im, op: op}
}

type patternPainterRGBA struct {
	im   *image.RGBA
	mask *image.Alpha
	p    Pattern
}

// Paint srcatisfies the Painter interface.
func (r *patternPainterRGBA) Paint(ss []raster.Span, done bool) {
	b := r.im.Bounds()
	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}
		y := s.Y - r.im.Rect.Min.Y
		x0 := s.X0 - r.im.Rect.Min.X
		// RGBAPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
		i0 := (s.Y-r.im.Rect.Min.Y)*r.im.Stride + (s.X0-r.im.Rect.Min.X)*4
		i1 := i0 + (s.X1-s.X0)*4
		for i, x := i0, x0; i < i1; i, x = i+4, x+1 {
			ma := s.Alpha
			if r.mask != nil {
				ma = ma * uint32(r.mask.AlphaAt(x, y).A) / 255
				if ma == 0 {
					continue
				}
			}
			c := r.p.ColorAt(x, y)
			cr, cg, cb, ca := c.RGBA()
			dr := uint32(r.im.Pix[i+0])
			dg := uint32(r.im.Pix[i+1])
			db := uint32(r.im.Pix[i+2])
			da := uint32(r.im.Pix[i+3])
			a := (m - (ca * ma / m)) * 0x101
			r.im.Pix[i+0] = uint8((dr*a + cr*ma) / m >> 8)
			r.im.Pix[i+1] = uint8((dg*a + cg*ma) / m >> 8)
			r.im.Pix[i+2] = uint8((db*a + cb*ma) / m >> 8)
			r.im.Pix[i+3] = uint8((da*a + ca*ma) / m >> 8)
		}
	}
}

func newPatternPainterRGBA(im *image.RGBA, mask *image.Alpha, p Pattern) *patternPainterRGBA {
	return &patternPainterRGBA{im, mask, p}
}

type patternPainterAlpha struct {
	im   *image.Alpha
	mask *image.Alpha
	p    Pattern
}

// Paint srcatisfies the Painter interface.
func (r *patternPainterAlpha) Paint(ss []raster.Span, done bool) {
	b := r.im.Bounds()
	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}
		y := s.Y - r.im.Rect.Min.Y
		x0 := s.X0 - r.im.Rect.Min.X
		// AlphaOverPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
		i0 := (s.Y-r.im.Rect.Min.Y)*r.im.Stride + (s.X0 - r.im.Rect.Min.X)
		i1 := i0 + (s.X1 - s.X0)

		for i, x := i0, x0; i < i1; i, x = i+1, x+1 {
			ma := s.Alpha
			if r.mask != nil {
				ma = ma * uint32(r.mask.AlphaAt(x, y).A) / 255
				if ma == 0 {
					continue
				}
			}

			srcy, _, _, srca := r.p.ColorAt(x, y).RGBA()
			dst := uint32(r.im.Pix[i])
			dst |= dst << 8

			srca = srca * ma >> 16              // combine source alpha with mask
			v := (srcy*srca + dst*(m-srca)) / m // blend source over destination
			r.im.Pix[i] = uint8(v >> 8)
		}
	}
}

func newPatternPainterAlpha(im *image.Alpha, mask *image.Alpha, p Pattern) *patternPainterAlpha {
	return &patternPainterAlpha{im, mask, p}
}
