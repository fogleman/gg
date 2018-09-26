package gg

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"image/color"
	"math/rand"
	"testing"
)

var save bool

func init() {
	flag.BoolVar(&save, "save", false, "save PNG output for each test case")
	flag.Parse()
}

func hash(dc *Context) string {
	return fmt.Sprintf("%x", sha256.Sum256(dc.im.Pix))
}

func checkHash(t *testing.T, dc *Context, expected string) {
	actual := hash(dc)
	if actual != expected {
		t.Fatalf("expected hash: %s != actual hash: %s", expected, actual)
	}
}

func saveImage(dc *Context, name string) error {
	if save {
		return SavePNG(name+".png", dc.Image())
	}
	return nil
}

func TestBlank(t *testing.T) {
	dc := NewContext(100, 100)
	saveImage(dc, "TestBlank")
	checkHash(t, dc, "e7e2dcff542de95352682dc186432e98f0188084896773f1973276b0577d5305")
}

func TestGrid(t *testing.T) {
	dc := NewContext(100, 100)
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
	checkHash(t, dc, "c1d18b7ceb06840c424496afc15322f0d9f1f56de99747266b790a52c34a6037")
}

func TestLines(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(0.5, 0.5, 0.5)
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
	checkHash(t, dc, "3f9586f9d3426a81095fdbcc903418b33e771baf8d34fb6f182d681a37b0f803")
}

func TestCircles(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(1, 1, 1)
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
	checkHash(t, dc, "69ccec48400136c74ab0d7a4ee7ca3d73d7da4fcbb760fb2472e220963a33a13")
}

func TestQuadratic(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(0.25, 0.25, 0.25)
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
	checkHash(t, dc, "e1b01f45f4c22d2db8cef53810aaec0bc3fa54f41ed25d8e8ba7388355468b13")
}

func TestCubic(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(0.75, 0.75, 0.75)
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
	checkHash(t, dc, "3ba5b6f78a382333697ad3da47dae354e9fdd729748c822b683e720db46918ce")
}

func TestFill(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(1, 1, 1)
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
	checkHash(t, dc, "2d55dcc91db735c4dba2167f582ee17d159d4400a968db0a5b1d506fdd829a07")
}

func TestClip(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(1, 1, 1)
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
	checkHash(t, dc, "3e53d9a182ae39d61ceb5d9afd2283f1c53eef10bcd21e05243ee06429210f62")
}

func TestPushPop(t *testing.T) {
	const S = 100
	dc := NewContext(S, S)
	dc.SetRGBA(0, 0, 0, 0.1)
	for i := 0; i < 360; i += 15 {
		dc.Push()
		dc.RotateAbout(Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}
	saveImage(dc, "TestPushPop")
	checkHash(t, dc, "f2c911424a06e82af510aac56e83490c911d3d9cec81ecd11f1bb0eddcd924e5")
}

func TestDrawStringWrapped(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped("Hello, world! How are you?", 50, 50, 0.5, 0.5, 90, 1.5, AlignCenter)
	saveImage(dc, "TestDrawStringWrapped")
	checkHash(t, dc, "f5f5c970236a8c66d223343b654469c0cc7e89fe85bca08386e0a2b9364cc344")
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

	dc := NewContext(200, 200)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.DrawImage(src.Image(), 50, 50)
	saveImage(dc, "TestDrawImage")
	checkHash(t, dc, "310b1d1aeaec0b1a066446f4f1cacc530a874f3e174d37c5899f27ed6dc856ed")
}

func TestSetPixel(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(0, 0, 0)
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
	checkHash(t, dc, "fd1f57cab467960643f3dec11083d1378925d55182c4cddd52f00f5e2f2f5766")
}

func TestDrawPoint(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(0, 0, 0)
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
	checkHash(t, dc, "ddd86935aac79117309cca794efd9be536e5ee966547a7e37986a2ab6b8b930a")
}

func TestLinearGradient(t *testing.T) {
	dc := NewContext(100, 100)
	g := NewLinearGradient(0, 0, 100, 100)
	g.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	g.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	g.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	dc.SetFillStyle(g)
	dc.DrawRectangle(0, 0, 100, 100)
	dc.Fill()
	saveImage(dc, "TestLinearGradient")
	checkHash(t, dc, "592a0fe58f2c1c69200a5734dc85db48665660b1af9452871d278dbbc668382c")
}

func TestRadialGradient(t *testing.T) {
	dc := NewContext(100, 100)
	g := NewRadialGradient(30, 50, 0, 70, 50, 50)
	g.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	g.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	g.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	dc.SetFillStyle(g)
	dc.DrawRectangle(0, 0, 100, 100)
	dc.Fill()
	saveImage(dc, "TestRadialGradient")
	checkHash(t, dc, "e567c0920410245c1f0c037907d3855d5eef5ef8b2f424dac871e4e87e9deaed")
}

func TestDashes(t *testing.T) {
	dc := NewContext(100, 100)
	dc.SetRGB(1, 1, 1)
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
	checkHash(t, dc, "6fdeb2653b4aa9a51663ff27d68dddc68f2fcd883ecfdb68974476af6a8d52c1")
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
