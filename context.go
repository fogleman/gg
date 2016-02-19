package gg

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
	lineCap   LineCap
	lineJoin  LineJoin
	fillRule  FillRule
	fontFace  font.Face
	matrix    Matrix
}

func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

func NewContextForRGBA(im *image.RGBA) *Context {
	return &Context{
		width:     im.Bounds().Size().X,
		height:    im.Bounds().Size().Y,
		im:        im,
		color:     color.Transparent,
		lineWidth: 1,
		fillRule:  FillRuleWinding,
		fontFace:  basicfont.Face7x13,
		matrix:    Identity(),
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

func (dc *Context) WritePNG(path string) error {
	return writePNG(path, dc.im)
}

func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}

func (dc *Context) SetLineCap(lineCap LineCap) {
	dc.lineCap = lineCap
}

func (dc *Context) SetLineCapRound() {
	dc.lineCap = LineCapRound
}

func (dc *Context) SetLineCapButt() {
	dc.lineCap = LineCapButt
}

func (dc *Context) SetLineCapSquare() {
	dc.lineCap = LineCapSquare
}

func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	dc.lineJoin = lineJoin
}

func (dc *Context) SetLineJoinRound() {
	dc.lineJoin = LineJoinRound
}

func (dc *Context) SetLineJoinBevel() {
	dc.lineJoin = LineJoinBevel
}

func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}

func (dc *Context) SetFillRuleWinding() {
	dc.fillRule = FillRuleWinding
}

func (dc *Context) SetFillRuleEvenOdd() {
	dc.fillRule = FillRuleEvenOdd
}

// Color Setters

func (dc *Context) SetColor(c color.Color) {
	dc.color = c
}

func (dc *Context) SetHexColor(x string) {
	r, g, b := parseHexColor(x)
	dc.SetRGB255(r, g, b)
}

func (dc *Context) SetRGBA255(r, g, b, a int) {
	dc.color = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.color = color.NRGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	}
}

func (dc *Context) SetRGB(r, g, b float64) {
	dc.SetRGBA(r, g, b, 1)
}

// Path Manipulation

func (dc *Context) MoveTo(x, y float64) {
	x, y = dc.TransformPoint(x, y)
	dc.start = fp(x, y)
	dc.path.Start(dc.start)
}

func (dc *Context) LineTo(x, y float64) {
	x, y = dc.TransformPoint(x, y)
	if len(dc.path) == 0 {
		dc.MoveTo(x, y)
	} else {
		dc.path.Add1(fp(x, y))
	}
}

func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
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

func (dc *Context) ClearPath() {
	dc.path.Clear()
}

// Path Drawing

func (dc *Context) StrokePreserve() {
	var capper raster.Capper
	switch dc.lineCap {
	case LineCapButt:
		capper = raster.ButtCapper
	case LineCapRound:
		capper = raster.RoundCapper
	case LineCapSquare:
		capper = raster.SquareCapper
	}
	var joiner raster.Joiner
	switch dc.lineJoin {
	case LineJoinBevel:
		joiner = raster.BevelJoiner
	case LineJoinRound:
		joiner = raster.RoundJoiner
	}
	painter := raster.NewRGBAPainter(dc.im)
	painter.SetColor(dc.color)
	r := raster.NewRasterizer(dc.width, dc.height)
	r.UseNonZeroWinding = true
	r.AddStroke(dc.path, fi(dc.lineWidth), capper, joiner)
	r.Rasterize(painter)
}

func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
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
	dc.ClearPath()
}

// Convenient Drawing Functions

func (dc *Context) Clear() {
	draw.Draw(dc.im, dc.im.Bounds(), image.NewUniform(dc.color), image.ZP, draw.Src)
}

func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

func (dc *Context) DrawRectangle(x, y, w, h float64) {
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.LineTo(x, y)
}

func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
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
	dc.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
}

func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

func (dc *Context) DrawCircle(x, y, r float64) {
	dc.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
}

// Text Functions

func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
}

func (dc *Context) LoadFontFace(path string, size float64) {
	dc.fontFace = loadFontFace(path, size)
}

func (dc *Context) DrawString(x, y float64, s string) {
	x, y = dc.TransformPoint(x, y)
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

// Transformation Matrix Operations

func (dc *Context) Identity() {
	dc.matrix = Identity()
}

func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

func (dc *Context) Rotate(angle float64) {
	dc.matrix = dc.matrix.Rotate(angle)
}

func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.matrix = dc.matrix.RotateAbout(angle, x, y)
}

func (dc *Context) Shear(x, y float64) {
	dc.matrix = dc.matrix.Shear(x, y)
}

func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	return dc.matrix.TransformPoint(x, y)
}
