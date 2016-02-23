package gg

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}
	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

func fp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
}

func fi(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

func unfix(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}

func loadFontFace(path string, points float64) font.Face {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
}

func flattenPath(p raster.Path) [][]Point {
	var result [][]Point
	var path []Point
	var cx, cy float64
	for i := 0; i < len(p); {
		switch p[i] {
		case 0:
			if len(path) > 0 {
				result = append(result, path)
				path = nil
			}
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 1:
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 2:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			points := QuadraticBezier(cx, cy, x1, y1, x2, y2)
			path = append(path, points...)
			cx, cy = x2, y2
			i += 6
		case 3:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			x3 := unfix(p[i+5])
			y3 := unfix(p[i+6])
			points := CubicBezier(cx, cy, x1, y1, x2, y2, x3, y3)
			path = append(path, points...)
			cx, cy = x3, y3
			i += 8
		default:
			panic("bad path")
		}
	}
	if len(path) > 0 {
		result = append(result, path)
	}
	return result
}

func dashPath(paths [][]Point, dashes []float64) [][]Point {
	var result [][]Point
	if len(dashes) == 0 {
		return paths
	}
	if len(dashes) == 1 {
		dashes = append(dashes, dashes[0])
	}
	for _, path := range paths {
		if len(path) < 2 {
			continue
		}
		previous := path[0]
		pathIndex := 1
		dashIndex := 0
		segmentLength := 0.0
		var segment []Point
		segment = append(segment, previous)
		for pathIndex < len(path) {
			dashLength := dashes[dashIndex]
			point := path[pathIndex]
			d := previous.Distance(point)
			maxd := dashLength - segmentLength
			if d > maxd {
				t := maxd / d
				p := previous.Interpolate(point, t)
				segment = append(segment, p)
				if dashIndex%2 == 0 && len(segment) > 1 {
					result = append(result, segment)
				}
				segment = nil
				segment = append(segment, p)
				segmentLength = 0
				previous = p
				dashIndex = (dashIndex + 1) % len(dashes)
			} else {
				segment = append(segment, point)
				previous = point
				segmentLength += d
				pathIndex++
			}
		}
		if dashIndex%2 == 0 && len(segment) > 1 {
			result = append(result, segment)
		}
	}
	return result
}

func rasterPath(paths [][]Point) raster.Path {
	var result raster.Path
	for _, path := range paths {
		var previous fixed.Point26_6
		for i, point := range path {
			f := point.Fixed()
			if i == 0 {
				result.Start(f)
			} else {
				dx := f.X - previous.X
				dy := f.Y - previous.Y
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx+dy > 4 {
					// TODO: this is a hack for cases where two points are
					// too close - causes rendering issues with joins / caps
					result.Add1(f)
				}
			}
			previous = f
		}
	}
	return result
}

func dashed(path raster.Path, dashes []float64) raster.Path {
	return rasterPath(dashPath(flattenPath(path), dashes))
}
