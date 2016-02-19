# Go Graphics

`gg` is a library for rendering 2D graphics in pure Go.

## Installation

    go get github.com/fogleman/gg

## Hello, Circle!

```go
package main

import "github.com/fogleman/gg"

func main() {
    // create a new 1000x1000 pixel image
    dc := gg.NewContext(1000, 1000)

    // draw a circle at the center with radius 400
    // nothing is drawn yet, we have to call Stroke or Fill
    dc.DrawCircle(500, 500, 400)

    // use black and 10px stroke width
    dc.SetRGB(0, 0, 0)
    dc.SetLineWidth(10)
    dc.Stroke()

    // write the result to file
    dc.SavePNG("out.png")
}
```

## Drawing Functions

Ever used a graphics library that didn't have functions for drawing rectangles
or circles? What a pain!

```go
DrawLine(x1, y1, x2, y2 float64)
DrawRectangle(x, y, w, h float64)
DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64)
DrawEllipse(x, y, rx, ry float64)
DrawArc(x, y, r, angle1, angle2 float64)
DrawCircle(x, y, r float64)
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
SetColor(c color.Color)
SetHexColor(x string)
SetRGBA255(r, g, b, a int)
SetRGB255(r, g, b int)
SetRGBA(r, g, b, a float64)
SetRGB(r, g, b float64)
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
Push()
Pop()
```
