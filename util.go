package dd

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"

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
		DPI:     96,
		Hinting: font.HintingFull,
	})
}
