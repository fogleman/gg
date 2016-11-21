package gg

import (
	"image/color"
	"math"
	"sort"
)

type stop struct {
	pos   float64
	color color.Color
}

type stops []stop

// Len satisfies the Sort interface.
func (s stops) Len() int {
	return len(s)
}

// Less satisfies the Sort interface.
func (s stops) Less(i, j int) bool {
	return s[i].pos < s[j].pos
}

// Swap satisfies the Sort interface.
func (s stops) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Gradient interface {
	Pattern
	AddColorStop(offset float64, color color.Color)
}

// Linear Gradient
type linearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

func (g *linearGradient) ColorAt(x, y int) color.Color {
	if len(g.stops) == 0 {
		return color.Transparent
	}

	fx, fy := float64(x), float64(y)
	x0, y0, x1, y1 := g.x0, g.y0, g.x1, g.y1
	dx, dy := x1-x0, y1-y0

	// Horizontal
	if dy == 0 && dx != 0 {
		return getColor((fx-x0)/dx, g.stops)
	}

	// Vertical
	if dx == 0 && dy != 0 {
		return getColor((fy-y0)/dy, g.stops)
	}

	// Dot product
	s0 := dx*(fx-x0) + dy*(fy-y0)
	if s0 < 0 {
		return g.stops[0].color
	}
	// Calculate distance to (x0,y0) alone (x0,y0)->(x1,y1)
	mag := math.Hypot(dx, dy)
	u := ((fx-x0)*-dy + (fy-y0)*dx) / (mag * mag)
	x2, y2 := x0+u*-dy, y0+u*dx
	d := math.Hypot(fx-x2, fy-y2) / mag
	return getColor(d, g.stops)
}

func (g *linearGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

func NewLinearGradient(x0, y0, x1, y1 float64) Gradient {
	g := &linearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
	return g
}

func getColor(pos float64, stops stops) color.Color {
	if pos <= 0.0 || len(stops) == 1 {
		return stops[0].color
	}

	last := stops[len(stops)-1]

	if pos >= last.pos {
		return last.color
	}

	for i, stop := range stops[1:] {
		if pos < stop.pos {
			pos = (pos - stops[i].pos) / (stop.pos - stops[i].pos)
			return colorLerp(stops[i].color, stop.color, pos)
		}
	}

	return last.color
}

func colorLerp(c0, c1 color.Color, t float64) color.Color {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()

	return color.NRGBA{
		lerp(r0, r1, t),
		lerp(g0, g1, t),
		lerp(b0, b1, t),
		lerp(a0, a1, t),
	}
}

func lerp(a, b uint32, t float64) uint8 {
	return uint8(int32(float64(a)*(1.0-t)+float64(b)*t) >> 8)
}
