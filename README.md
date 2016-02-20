# Go Graphics

`gg` is a library for rendering 2D graphics in pure Go.

## Installation

    go get github.com/fogleman/gg

## Hello, Circle!

Look how easy!

```go
package main

import "github.com/fogleman/gg"

func main() {
    dc := gg.NewContext(1000, 1000)
    dc.DrawCircle(500, 500, 400)
    dc.SetRGB(0, 0, 0)
    dc.Fill()
    dc.SavePNG("out.png")
}
```

## Creating Contexts

There are a few ways of creating a context.

```go
NewContext(width, height int) *Context
NewContextForImage(im image.Image) *Context
NewContextForRGBA(im *image.RGBA) *Context
```

## Drawing Functions

Ever used a graphics library that didn't have functions for drawing rectangles
or circles? What a pain!

```go
DrawLine(x1, y1, x2, y2 float64)
DrawRectangle(x, y, w, h float64)
DrawRoundedRectangle(x, y, w, h, r float64)
DrawCircle(x, y, r float64)
DrawArc(x, y, r, angle1, angle2 float64)
DrawEllipse(x, y, rx, ry float64)
DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64)
DrawImage(im image.Image, x, y int)
DrawString(x, y float64, s string)

MoveTo(x, y float64)
LineTo(x, y float64)
QuadraticTo(x1, y1, x2, y2 float64)
ClosePath()
ClearPath()

Clear()
Stroke()
Fill()
StrokePreserve()
FillPreserve()
```

## Color Functions

Colors can be set in several different ways for your convenience.

```go
SetRGB(r, g, b float64)
SetRGBA(r, g, b, a float64)
SetRGB255(r, g, b int)
SetRGBA255(r, g, b, a int)
SetColor(c color.Color)
SetHexColor(x string)
```

## Stroke & Fill Options

```go
SetLineWidth(lineWidth float64)
SetLineCap(lineCap LineCap)
SetLineJoin(lineJoin LineJoin)
SetFillRule(fillRule FillRule)
```

## Transformation Functions

```go
Identity()
Translate(x, y float64)
Scale(x, y float64)
Rotate(angle float64)
RotateAbout(angle, x, y float64)
Shear(x, y float64)
TransformPoint(x, y float64) (tx, ty float64)
InvertY()
Push()
Pop()
```

## Helper Functions

Sometimes you just don't want to write these yourself.

```go
Radians(degrees float64) float64
Degrees(radians float64) float64
LoadPNG(path string) (image.Image, error)
SavePNG(path string, im image.Image) error
```

## What's Missing?

If you need any of the features below, I recommend using `cairo` instead. Or even better, implement it and submit a pull request!

- Clipping Regions
- Gradients / Patterns
- Cubic Beziers
- Dashed Lines

## How Do it Do?

`gg` is mostly a wrapper around `github.com/golang/freetype/raster`. The goal is to provide some more functionality and a nicer API that will suffice for most use cases.

## Another Example

See the output of this example below.

```go
package main

import "github.com/fogleman/gg"

func main() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGBA(0, 0, 0, 0.1)
	for i := 0; i < 360; i += 15 {
		dc.Push()
		dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}
	dc.SavePNG("out.png")
}
```

![Ellipses](http://i.imgur.com/J9CBZef.png)
