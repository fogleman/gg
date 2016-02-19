package dd

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinRound LineJoin = iota
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type Context struct {
	width     int
	height    int
	im        *image.RGBA
	color     color.Color
	path      raster.Path
	start     fixed.Point26_6
	lineWidth float64
	capper    raster.Capper
	joiner    raster.Joiner
	fillRule  FillRule
	fontFace  font.Face
}

func NewContext(width, height int) *Context {
	im := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Context{
		width:     width,
		height:    height,
		im:        im,
		color:     color.Transparent,
		lineWidth: 1,
		fillRule:  FillRuleWinding,
		fontFace:  basicfont.Face7x13,
	}
}

func (dc *Context) Image() image.Image {
	return dc.im
}

func (dc *Context) Width() int {
	return dc.width
}

func (dc *Context) Height() int {
	return dc.height
}

func (dc *Context) WriteToPNG(path string) error {
	return writeToPNG(path, dc.im)
}

func (dc *Context) Paint() {
	draw.Draw(dc.im, dc.im.Bounds(), image.NewUniform(dc.color), image.ZP, draw.Src)
}

func (dc *Context) SetSourceRGBA(r, g, b, a float64) {
	dc.color = color.NRGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	}
}

func (dc *Context) SetSourceRGB(r, g, b float64) {
	dc.SetSourceRGBA(r, g, b, 1)
}

func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}

func (dc *Context) SetLineCap(lineCap LineCap) {
	switch lineCap {
	case LineCapButt:
		dc.capper = raster.ButtCapper
	case LineCapRound:
		dc.capper = raster.RoundCapper
	case LineCapSquare:
		dc.capper = raster.SquareCapper
	}
}

func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	switch lineJoin {
	case LineJoinBevel:
		dc.joiner = raster.BevelJoiner
	case LineJoinRound:
		dc.joiner = raster.RoundJoiner
	}
}

func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}

func (dc *Context) MoveTo(x, y float64) {
	dc.start = fp(x, y)
	dc.path.Start(dc.start)
}

func (dc *Context) LineTo(x, y float64) {
	if len(dc.path) == 0 {
		dc.MoveTo(x, y)
	} else {
		dc.path.Add1(fp(x, y))
	}
}

func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if len(dc.path) == 0 {
		dc.MoveTo(x1, y1)
	} else {
		dc.path.Add2(fp(x1, y1), fp(x2, y2))
	}
}

func (dc *Context) ClosePath() {
	if len(dc.path) > 0 {
		dc.path.Add1(dc.start)
	}
}

func (dc *Context) NewPath() {
	dc.path.Clear()
}

func (dc *Context) StrokePreserve() {
	painter := raster.NewRGBAPainter(dc.im)
	painter.SetColor(dc.color)
	r := raster.NewRasterizer(dc.width, dc.height)
	r.UseNonZeroWinding = true
	r.AddStroke(dc.path, fi(dc.lineWidth), dc.capper, dc.joiner)
	r.Rasterize(painter)
}

func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.NewPath()
}

func (dc *Context) FillPreserve() {
	// make sure the path is closed
	path := make(raster.Path, len(dc.path))
	copy(path, dc.path)
	path.Add1(dc.start)
	painter := raster.NewRGBAPainter(dc.im)
	painter.SetColor(dc.color)
	r := raster.NewRasterizer(dc.width, dc.height)
	r.UseNonZeroWinding = dc.fillRule == FillRuleWinding
	r.AddPath(path)
	r.Rasterize(painter)
}

func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.NewPath()
}

// Convenient Drawing Functions

func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

func (dc *Context) DrawEllipseArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i <= n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos(a1+(a2-a1)/2)
		y1 := y + ry*math.Sin(a1+(a2-a1)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			dc.MoveTo(x0, y0)
		}
		dc.QuadraticTo(cx, cy, x2, y2)
	}
}

func (dc *Context) DrawEllipse(x, y, rx, ry float64) {
	dc.DrawEllipseArc(x, y, rx, ry, 0, 2*math.Pi)
}

func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipseArc(x, y, r, r, angle1, angle2)
}

func (dc *Context) DrawCircle(x, y, r float64) {
	dc.DrawEllipseArc(x, y, r, r, 0, 2*math.Pi)
}

// Text Functions

func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
}

func (dc *Context) LoadFontFace(path string, size float64) {
	dc.fontFace = loadFontFace(path, size)
}

func (dc *Context) DrawString(x, y float64, s string) {
	d := &font.Drawer{
		Dst:  dc.im,
		Src:  image.NewUniform(dc.color),
		Face: dc.fontFace,
		Dot:  fp(x, y),
	}
	d.DrawString(s)
}

func (dc *Context) MeasureString(s string) float64 {
	d := &font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: dc.fontFace,
	}
	a := d.MeasureString(s)
	return float64(a >> 6)
}
