package dd

import (
	"image"
	"image/png"
	"os"

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
