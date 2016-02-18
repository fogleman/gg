package dd

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/golang/freetype/raster"
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
}

func NewContext(width, height int) *Context {
	im := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Context{
		width:     width,
		height:    height,
		im:        im,
		color:     color.Transparent,
		lineWidth: 1,
	}
}

func (c *Context) Image() image.Image {
	return c.im
}

func (c *Context) Width() int {
	return c.width
}

func (c *Context) Height() int {
	return c.height
}

func (c *Context) WriteToPNG(path string) error {
	return writeToPNG(path, c.im)
}

func (c *Context) Paint() {
	draw.Draw(c.im, c.im.Bounds(), image.NewUniform(c.color), image.ZP, draw.Src)
}

func (c *Context) SetSourceRGBA(r, g, b, a float64) {
	c.color = color.RGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	}
}

func (c *Context) SetSourceRGB(r, g, b float64) {
	c.SetSourceRGBA(r, g, b, 1)
}

func (c *Context) SetLineWidth(lineWidth float64) {
	c.lineWidth = lineWidth
}

func (c *Context) SetLineCap(lineCap LineCap) {
	switch lineCap {
	case LineCapButt:
		c.capper = raster.ButtCapper
	case LineCapRound:
		c.capper = raster.RoundCapper
	case LineCapSquare:
		c.capper = raster.SquareCapper
	}
}

func (c *Context) SetLineJoin(lineJoin LineJoin) {
	switch lineJoin {
	case LineJoinBevel:
		c.joiner = raster.BevelJoiner
	case LineJoinRound:
		c.joiner = raster.RoundJoiner
	}
}

func (c *Context) MoveTo(x, y float64) {
	c.start = fp(x, y)
	c.path.Start(c.start)
}

func (c *Context) LineTo(x, y float64) {
	c.path.Add1(fp(x, y))
}

func (c *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	c.path.Add2(fp(x1, y1), fp(x2, y2))
}

func (c *Context) ClosePath() {
	c.path.Add1(c.start)
}

func (c *Context) NewPath() {
	c.path.Clear()
}

func (c *Context) StrokePreserve() {
	painter := raster.NewRGBAPainter(c.im)
	painter.SetColor(c.color)
	r := raster.NewRasterizer(c.width, c.height)
	r.UseNonZeroWinding = true
	r.AddStroke(c.path, fi(c.lineWidth), c.capper, c.joiner)
	r.Rasterize(painter)
}

func (c *Context) Stroke() {
	c.StrokePreserve()
	c.NewPath()
}

func (c *Context) FillPreserve() {
	// make sure the path is closed
	path := make(raster.Path, len(c.path))
	copy(path, c.path)
	path.Add1(c.start)
	painter := raster.NewRGBAPainter(c.im)
	painter.SetColor(c.color)
	r := raster.NewRasterizer(c.width, c.height)
	r.AddPath(path)
	r.Rasterize(painter)
}

func (c *Context) Fill() {
	c.FillPreserve()
	c.NewPath()
}
