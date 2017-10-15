package gg

import (
	"image"
	"strings"
)

// PlotArgs holds the arguments for drawing line graphs (charts)
type PlotArgs struct {

	// curve name
	L string // label

	// lines
	C  string  // line: color
	A  float64 // line: alpha (0, 1]. A<1e-14 => A=1.0
	Ls string  // line: style
	Lw float64 // line: width

	// markers
	M    string  // marker: type, e.g. "o", "s", "+", "x", "img:filename.png"
	Ms   int     // marker: size
	Me   int     // marker: mark-every
	Mec  string  // marker: edge color
	Mew  float64 // marker: edge width
	Void bool    // marker: void marker (draw edge only)

	// internal
	markerImg  image.Image // marker image
	markerName string      // filename corresponding to loaded image
}

// DrawMarker draws marker
func (o *PlotArgs) DrawMarker(dc *Context, x, y int) {

	// skip if marker type is empty
	if o.M == "" {
		return
	}

	// markersize and half-markersize
	s := o.markerSize()
	h := s / 2

	// draw marker
	switch o.M {

	// circle
	case "o":
		if !o.Void {
			o.Circle(dc, true, false, x, y, h)
		}
		o.Circle(dc, true, true, x, y, h)

	// square
	case "s":
		if !o.Void {
			o.Rect(dc, true, false, x-h, y-h, s, s)
		}
		o.Rect(dc, true, true, x-h, y-h, s, s)

	// cross
	case "+":
		o.Line(dc, true, true, x-h, y, x+h, y)
		o.Line(dc, true, true, x, y-h, x, y+h)

	// x
	case "x":
		o.Line(dc, true, true, x-h, y-h, x+h, y+h)
		o.Line(dc, true, true, x-h, y+h, x+h, y-h)

	// use image as marker
	default:
		if o.checkMarkerImage() {
			dc.DrawImageAnchored(o.markerImg, x, y, 0.5, 0.5)
		}
	}
}

// Activate activates properties of lines/shapes
func (o *PlotArgs) Activate(dc *Context, marker, edge bool) {

	// color
	clr := o.C
	if marker && edge && o.Mec != "" {
		clr = o.Mec
	}
	r, g, b, _ := parseHexColor(clr)

	// alpha
	alpha := o.A
	if alpha < 1e-14 || marker {
		alpha = 1.0
	}

	// set color
	dc.SetRGBA255(r, g, b, int(alpha*255))
}

// Circle draws Circle
func (o *PlotArgs) Circle(dc *Context, marker, edge bool, x, y, r int) {
	o.Activate(dc, marker, edge)
	dc.DrawCircle(float64(x), float64(y), float64(r))
	if edge {
		dc.Stroke()
	} else {
		dc.Fill()
	}
}

// Rect draws rectangle
func (o *PlotArgs) Rect(dc *Context, marker, edge bool, x, y, w, h int) {
	o.Activate(dc, marker, edge)
	dc.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	if edge {
		dc.Stroke()
	} else {
		dc.Fill()
	}
}

// Line draws Line
func (o *PlotArgs) Line(dc *Context, marker, edge bool, x1, y1, x2, y2 int) {
	o.Activate(dc, marker, edge)
	dc.DrawLine(float64(x1), float64(y1), float64(x2), float64(y2))
	if edge {
		dc.Stroke()
	} else {
		dc.Fill()
	}
}

// checkMarkerImage loads marker image if not loaded already
func (o *PlotArgs) checkMarkerImage() (useImage bool) {
	useImage = strings.HasPrefix(o.M, "img:")
	if useImage {
		fn := strings.TrimPrefix(o.M, "img:")
		if o.markerImg == nil || o.markerName != fn { // load only if not loaded already
			var err error
			o.markerImg, err = LoadPNG(fn)
			if err != nil {
				panic(err)
			}
			o.markerName = fn
		}
	}
	return
}

// markerSize returns the size of marker
func (o *PlotArgs) markerSize() int {
	if o.checkMarkerImage() {
		return o.markerImg.Bounds().Dy()
	}
	if o.Ms == 0 {
		return 8 // default value
	}
	return o.Ms
}

// colors /////////////////////////////////////////////////////////////////////////////////////////

// GetColor returns a color from a default palette
//  use palette < 0 for automatic color
func GetColor(i, palette int) string {
	if palette < 0 || palette >= len(palettes) {
		return ""
	}
	p := palettes[palette]
	return p[i%len(p)]
}

// palettes holds color palettes
var palettes = [][]string{
	{"#003fff", "#35b052", "#e8000b", "#8a2be2", "#ffc400", "#00d7ff"},
	{"blue", "green", "magenta", "orange", "red", "cyan", "black", "#de9700", "#89009d", "#7ad473", "#737ad4", "#d473ce", "#7e6322", "#462222", "#98ac9d", "#37a3e8", "yellow"},
	{"#4c72b0", "#55a868", "#c44e52", "#8172b2", "#ccb974", "#64b5cd"},
	{"#9b59b6", "#3498db", "#95a5a6", "#e74c3c", "#34495e", "#2ecc71"},
	{"#e41a1c", "#377eb8", "#4daf4a", "#984ea3", "#ff7f00", "#ffff33"},
	{"#7fc97f", "#beaed4", "#fdc086", "#ffff99", "#386cb0", "#f0027f", "#bf5b17"},
	{"#001c7f", "#017517", "#8c0900", "#7600a1", "#b8860b", "#006374"},
	{"#0072b2", "#009e73", "#d55e00", "#cc79a7", "#f0e442", "#56b4e9"},
	{"#4878cf", "#6acc65", "#d65f5f", "#b47cc7", "#c4ad66", "#77bedb"},
	{"#92c6ff", "#97f0aa", "#ff9f9a", "#d0bbff", "#fffea3", "#b0e0e6"},
}
