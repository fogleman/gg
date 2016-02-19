package gg

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func writeToPNG(path string, im image.Image) error {
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

func parseHexColor(x string) (r, g, b int) {
	x = strings.TrimPrefix(x, "#")
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
	return
}

func fp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
}

func fi(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

func loadFontFace(path string, size float64) font.Face {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(f, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	})
}
