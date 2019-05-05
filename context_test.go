package gg

import (
	"crypto/md5"
	"flag"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"testing"
)

var save bool

func init() {
	flag.BoolVar(&save, "save", true, "save PNG output for each test case")
	flag.Parse()
}

func hash(dc *Context) string {
	var pixels []uint8
	switch dc.pm {
	case painterModeRGBA:
		pixels = dc.im.(*image.RGBA).Pix
	case painterModeAlpha:
		pixels = dc.im.(*image.Alpha).Pix
	}
	return fmt.Sprintf("%x", md5.Sum(pixels))
}

func checkHash(t *testing.T, dc *Context, expected string) {
	actual := hash(dc)
	if actual != expected {
		t.Fatalf("expected hash: %s != actual hash: %s", expected, actual)
	}
}

func saveImage(dc *Context, name string) error {
	if save {
		return SavePNG(fmt.Sprintf("%s-%s.png", name, dc.pm), dc.Image())
	}
	return nil
}

func TestBlank(t *testing.T) {
	tests := []struct {
		dc   *Context
		hash string
	}{
		{NewContext(100, 100), "4e0a293a5b638f0aba2c4fe2c3418d0e"},
		{NewAlphaContext(100, 100), "b85d6fb9ef4260dcf1ce0a1b0bff80d3"},
	}
	for _, tc := range tests {
		saveImage(tc.dc, "TestBlank")
		checkHash(t, tc.dc, tc.hash)
	}
}

func TestGrid(t *testing.T) {
	tests := []struct {
		dc   *Context
		hash string
	}{
		{NewContext(100, 100), "78606adda71d8abfbd8bb271087e4d69"},
		{NewAlphaContext(100, 100), "7bc189656cceb04f285f3cbe49637b17"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetRGB(1, 1, 1)
		dc.Clear()
		for i := 10; i < 100; i += 10 {
			x := float64(i) + 0.5
			dc.DrawLine(x, 0, x, 100)
			dc.DrawLine(0, x, 100, x)
		}
		dc.SetRGB(0, 0, 0)
		dc.Stroke()
		saveImage(dc, "TestGrid")
		checkHash(t, dc, tc.hash)
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.NRGBA{127, 127, 127, 255}, "036bd220e2529955cc48425dd72bb686"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "d416070054d1f8e089791d0f3cb915bc"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 100; i++ {
			x1 := rnd.Float64() * 100
			y1 := rnd.Float64() * 100
			x2 := rnd.Float64() * 100
			y2 := rnd.Float64() * 100
			dc.DrawLine(x1, y1, x2, y2)
			dc.SetLineWidth(rnd.Float64() * 3)
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Stroke()
		}
		saveImage(dc, "TestLines")
		checkHash(t, dc, tc.hash)
	}
}

func TestCircles(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.White, "c52698000df96fabafe7863701afe922"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "0e1d09f170a9f8c048fc7975afe67c35"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 10; i++ {
			x := rnd.Float64() * 100
			y := rnd.Float64() * 100
			r := rnd.Float64()*10 + 5
			dc.DrawCircle(x, y, r)
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.FillPreserve()
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.SetLineWidth(rnd.Float64() * 3)
			dc.Stroke()
		}
		saveImage(dc, "TestCircles")
		checkHash(t, dc, tc.hash)
	}
}

func TestQuadratic(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.NRGBA{63, 63, 63, 255}, "56b842d814aee94b52495addae764a77"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "a6daf2aa412e1fdc2511ca7df3cf3934"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 100; i++ {
			x1 := rnd.Float64() * 100
			y1 := rnd.Float64() * 100
			x2 := rnd.Float64() * 100
			y2 := rnd.Float64() * 100
			x3 := rnd.Float64() * 100
			y3 := rnd.Float64() * 100
			dc.MoveTo(x1, y1)
			dc.QuadraticTo(x2, y2, x3, y3)
			dc.SetLineWidth(rnd.Float64() * 3)
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Stroke()
		}
		saveImage(dc, "TestQuadratic")
		checkHash(t, dc, tc.hash)
	}
}

func TestCubic(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.NRGBA{191, 191, 191, 255}, "4a7960fc4eaaa33ce74131c5ce0afca8"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "9a08af24bb02b2bbacb3eb0c9ebe17dc"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 100; i++ {
			x1 := rnd.Float64() * 100
			y1 := rnd.Float64() * 100
			x2 := rnd.Float64() * 100
			y2 := rnd.Float64() * 100
			x3 := rnd.Float64() * 100
			y3 := rnd.Float64() * 100
			x4 := rnd.Float64() * 100
			y4 := rnd.Float64() * 100
			dc.MoveTo(x1, y1)
			dc.CubicTo(x2, y2, x3, y3, x4, y4)
			dc.SetLineWidth(rnd.Float64() * 3)
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Stroke()
		}
		saveImage(dc, "TestCubic")
		checkHash(t, dc, tc.hash)
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.White, "7ccb3a2443906a825e57ab94db785467"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "d1594288f3a31d58172fb355c6c244b1"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 10; i++ {
			dc.NewSubPath()
			for j := 0; j < 10; j++ {
				x := rnd.Float64() * 100
				y := rnd.Float64() * 100
				dc.LineTo(x, y)
			}
			dc.ClosePath()
			dc.SetRGBA(rnd.Float64(), rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Fill()
		}
		saveImage(dc, "TestFill")
		checkHash(t, dc, tc.hash)
	}
}

func TestClip(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.White, "762c32374d529fd45ffa038b05be7865"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "32c782c6a0cb723d54edeb7419fe9663"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		dc.DrawCircle(50, 50, 40)
		dc.Clip()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 1000; i++ {
			x := rnd.Float64() * 100
			y := rnd.Float64() * 100
			r := rnd.Float64()*10 + 5
			dc.DrawCircle(x, y, r)
			dc.SetRGBA(rnd.Float64(), rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Fill()
		}
		saveImage(dc, "TestClip")
		checkHash(t, dc, tc.hash)
	}
}

func TestPushPop(t *testing.T) {
	const S = 100
	tests := []struct {
		dc    *Context
		color color.Color
		hash  string
	}{
		{NewContext(S, S), color.NRGBA{0, 0, 0, 25}, "31e908ee1c2ea180da98fd5681a89d05"},
		{NewAlphaContext(S, S), color.NRGBA{255, 255, 255, 96}, "3db885602672f7be40197486342863e9"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.color)
		for i := 0; i < 360; i += 15 {
			dc.Push()
			dc.RotateAbout(Radians(float64(i)), S/2, S/2)
			dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
			dc.Fill()
			dc.Pop()
		}
		saveImage(dc, "TestPushPop")
		checkHash(t, dc, tc.hash)
	}
}

func TestDrawStringWrapped(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.White, "8d92f6aae9e8b38563f171abd00893f8"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "10b4946c822de9777f3f0ee2f2445ec9"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		dc.DrawStringWrapped("Hello, world! How are you?", 50, 50, 0.5, 0.5, 90, 1.5, AlignCenter)
		saveImage(dc, "TestDrawStringWrapped")
		checkHash(t, dc, tc.hash)
	}
}

func TestDrawImage(t *testing.T) {
	src := NewContext(100, 100)
	src.SetRGB(1, 1, 1)
	src.Clear()
	for i := 10; i < 100; i += 10 {
		x := float64(i) + 0.5
		src.DrawLine(x, 0, x, 100)
		src.DrawLine(0, x, 100, x)
	}
	src.SetRGB(0, 0, 0)
	src.Stroke()

	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(200, 200), color.Black, "282afbc134676722960b6bec21305b15"},
		{NewAlphaContext(200, 200), color.NRGBA{}, "0234f7f19198ef5a1a3c6f630e66abbf"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		dc.DrawImage(src.Image(), 50, 50)
		saveImage(dc, "TestDrawImage")
		checkHash(t, dc, tc.hash)
	}
}

func TestSetPixel(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.Black, "27dda6b4b1d94f061018825b11982793"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "36501f522cd66fae458da39d3ac6a1e5"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		dc.SetRGB(0, 1, 0)
		i := 0
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				if i%31 == 0 {
					dc.SetPixel(x, y)
				}
				i++
			}
		}
		saveImage(dc, "TestSetPixel")
		checkHash(t, dc, tc.hash)
	}
}

func TestDrawPoint(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.Black, "55af8874531947ea6eeb62222fb33e0e"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "2b881959975edb768228d9d8bc2081fa"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		dc.SetRGB(0, 1, 0)
		dc.Scale(10, 10)
		for y := 0; y <= 10; y++ {
			for x := 0; x <= 10; x++ {
				dc.DrawPoint(float64(x), float64(y), 3)
				dc.Fill()
			}
		}
		saveImage(dc, "TestDrawPoint")
		checkHash(t, dc, tc.hash)
	}
}

func TestLinearGradient(t *testing.T) {
	tests := []struct {
		dc    *Context
		stops [3]color.Color
		hash  string
	}{
		{NewContext(100, 100), [3]color.Color{
			color.RGBA{0, 255, 0, 255}, // green
			color.RGBA{255, 0, 0, 255}, // red
			color.RGBA{0, 0, 255, 255}, // blue
		}, "75eb9385c1219b1d5bb6f4c961802c7a"},
		{NewAlphaContext(100, 100), [3]color.Color{
			color.RGBA{255, 255, 255, 0},
			color.RGBA{255, 255, 255, 255},
			color.RGBA{255, 255, 255, 192},
		}, "74305fe008a340d01488fd2f443b72d8"},
	}
	for _, tc := range tests {
		dc := tc.dc

		g := NewLinearGradient(0, 0, 100, 100)
		g.AddColorStop(0.0, tc.stops[0])
		g.AddColorStop(1.0, tc.stops[2])
		g.AddColorStop(0.5, tc.stops[1])
		dc.SetFillStyle(g)
		dc.DrawRectangle(0, 0, 100, 100)
		dc.Fill()
		saveImage(dc, "TestLinearGradient")
		checkHash(t, dc, tc.hash)
	}
}

func TestRadialGradient(t *testing.T) {
	tests := []struct {
		dc    *Context
		stops [3]color.Color
		hash  string
	}{
		{NewContext(100, 100), [3]color.Color{
			color.RGBA{0, 255, 0, 255}, // green
			color.RGBA{255, 0, 0, 255}, // red
			color.RGBA{0, 0, 255, 255}, // blue
		}, "f170f39c3f35c29de11e00428532489d"},
		{NewAlphaContext(100, 100), [3]color.Color{
			color.RGBA{255, 255, 255, 0},
			color.RGBA{255, 255, 255, 255},
			color.RGBA{255, 255, 255, 192},
		}, "033794fee3c193b464b4564c038b6242"},
	}
	for _, tc := range tests {
		dc := tc.dc

		g := NewRadialGradient(30, 50, 0, 70, 50, 50)
		g.AddColorStop(0.0, tc.stops[0])
		g.AddColorStop(1.0, tc.stops[2])
		g.AddColorStop(0.5, tc.stops[1])
		dc.SetFillStyle(g)
		dc.DrawRectangle(0, 0, 100, 100)
		dc.Fill()
		saveImage(dc, "TestRadialGradient")
		checkHash(t, dc, tc.hash)
	}
}

func TestDashes(t *testing.T) {
	tests := []struct {
		dc   *Context
		bg   color.Color
		hash string
	}{
		{NewContext(100, 100), color.White, "d188069c69dcc3970edfac80f552b53c"},
		{NewAlphaContext(100, 100), color.NRGBA{}, "d46f9b8d479d85418416f054395bd984"},
	}
	for _, tc := range tests {
		dc := tc.dc

		dc.SetColor(tc.bg)
		dc.Clear()
		rnd := rand.New(rand.NewSource(99))
		for i := 0; i < 100; i++ {
			x1 := rnd.Float64() * 100
			y1 := rnd.Float64() * 100
			x2 := rnd.Float64() * 100
			y2 := rnd.Float64() * 100
			dc.SetDash(rnd.Float64()*3+1, rnd.Float64()*3+3)
			dc.DrawLine(x1, y1, x2, y2)
			dc.SetLineWidth(rnd.Float64() * 3)
			dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
			dc.Stroke()
		}
		saveImage(dc, "TestDashes")
		checkHash(t, dc, tc.hash)
	}
}

func BenchmarkCircles(b *testing.B) {
	dc := NewContext(1000, 1000)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	rnd := rand.New(rand.NewSource(99))
	for i := 0; i < b.N; i++ {
		x := rnd.Float64() * 1000
		y := rnd.Float64() * 1000
		dc.DrawCircle(x, y, 10)
		if i%2 == 0 {
			dc.SetRGB(0, 0, 0)
		} else {
			dc.SetRGB(1, 1, 1)
		}
		dc.Fill()
	}
}
