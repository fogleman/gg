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
	width      int
	height     int
	im         *image.RGBA
	color      color.Color
	path       raster.Path
	start      fixed.Point26_6
	lineWidth  float64
	lineCap    LineCap
	lineJoin   LineJoin
	fillRule   FillRule
	fontFace   font.Face
	fontHeight float64
	matrix     Matrix
	stack      []*Context
}

func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

func NewContextForRGBA(im *image.RGBA) *Context {
	return &Context{
		width:      im.Bounds().Size().X,
		height:     im.Bounds().Size().Y,
		im:         im,
		color:      color.Transparent,
		lineWidth:  1,
		fillRule:   FillRuleWinding,
		fontFace:   basicfont.Face7x13,
		fontHeight: 13,
		matrix:     Identity(),
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

func (dc *Context) SavePNG(path string) error {
	return SavePNG(path, dc.im)
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

func (dc *Context) capper() raster.Capper {
	switch dc.lineCap {
	case LineCapButt:
		return raster.ButtCapper
	case LineCapRound:
		return raster.RoundCapper
	case LineCapSquare:
		return raster.SquareCapper
	}
	return nil
}

func (dc *Context) joiner() raster.Joiner {
	switch dc.lineJoin {
	case LineJoinBevel:
		return raster.BevelJoiner
	case LineJoinRound:
		return raster.RoundJoiner
	}
	return nil
}

func (dc *Context) StrokePreserve() {
	painter := raster.NewRGBAPainter(dc.im)
	painter.SetColor(dc.color)
	r := raster.NewRasterizer(dc.width, dc.height)
	r.UseNonZeroWinding = true
	r.AddStroke(dc.path, fi(dc.lineWidth), dc.capper(), dc.joiner())
	r.Rasterize(painter)
}

func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
}

func (dc *Context) FillPreserve() {
	painter := raster.NewRGBAPainter(dc.im)
	painter.SetColor(dc.color)
	r := raster.NewRasterizer(dc.width, dc.height)
	r.UseNonZeroWinding = dc.fillRule == FillRuleWinding
	r.AddPath(dc.path)
	r.Rasterize(painter)
}

func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// Convenient Drawing Functions

func (dc *Context) Clear() {
	src := image.NewUniform(dc.color)
	draw.Draw(dc.im, dc.im.Bounds(), src, image.ZP, draw.Src)
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

func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawArc(x2, y1, r, Radians(270), Radians(360))
	dc.LineTo(x3, y2)
	dc.DrawArc(x2, y2, r, Radians(0), Radians(90))
	dc.LineTo(x1, y3)
	dc.DrawArc(x1, y2, r, Radians(90), Radians(180))
	dc.LineTo(x0, y1)
	dc.DrawArc(x1, y1, r, Radians(180), Radians(270))
}

func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i < n; i++ {
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

func (dc *Context) DrawImage(im image.Image, x, y int) {
	dc.DrawImageAnchored(im, x, y, 0, 0)
}

func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	s := im.Bounds().Size()
	x -= int(ax * float64(s.X))
	y -= int(ay * float64(s.Y))
	p := image.Pt(x, y)
	r := image.Rectangle{p, p.Add(s)}
	draw.Draw(dc.im, r, im, image.ZP, draw.Over)
}

// Text Functions

func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
}

func (dc *Context) LoadFontFace(path string, points float64) {
	dc.fontFace = loadFontFace(path, points)
	dc.fontHeight = points * 72 / 96
}

func (dc *Context) DrawString(s string, x, y float64) {
	dc.DrawStringAnchored(s, x, y, 0, 0)
}

func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	w, h := dc.MeasureString(s)
	x -= ax * w
	y += ay * h
	x, y = dc.TransformPoint(x, y)
	d := &font.Drawer{
		Dst:  dc.im,
		Src:  image.NewUniform(dc.color),
		Face: dc.fontFace,
		Dot:  fp(x, y),
	}
	d.DrawString(s)
}

func (dc *Context) MeasureString(s string) (w, h float64) {
	d := &font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: dc.fontFace,
	}
	a := d.MeasureString(s)
	return float64(a >> 6), dc.fontHeight
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

func (dc *Context) ScaleAbout(sx, sy, x, y float64) {
	dc.matrix = dc.matrix.ScaleAbout(sx, sy, x, y)
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

func (dc *Context) InvertY() {
	dc.Translate(0, float64(dc.height))
	dc.Scale(1, -1)
}

// Stack

func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

func (dc *Context) Pop() {
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.path = before.path
}
