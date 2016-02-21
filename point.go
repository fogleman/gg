package gg

import "golang.org/x/image/math/fixed"

type Point struct {
	X, Y float64
}

func (p Point) Fixed() fixed.Point26_6 {
	return fp(p.X, p.Y)
}
