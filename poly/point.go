package poly

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func Sub(a, b Point) Point {
	return Point{b.X - a.X, b.Y - a.Y}
}

const epsilon = 1e-6

func Gradient(p Point) []string {
	ax := math.Abs(p.X)
	ay := math.Abs(p.Y)
	if ax <= epsilon && ay <= epsilon {
		return []string{"0"}
	}
	var xdir, ydir, dir string
	var grad float64
	if p.X < 0 {
		xdir = "-x"
	} else {
		xdir = "+x"
	}
	if p.Y < 0 {
		ydir += "-y"
	} else {
		ydir += "+y"
	}
	if ax < ay {
		dir = xdir + ydir
		grad = ax / ay
	} else {
		dir = ydir + xdir
		grad = ay / ax
	}
	dist := math.Sqrt(p.X*p.X + p.Y*p.Y)
	return []string{dir, fmt.Sprintf("%.5f", grad), fmt.Sprintf("%.5f", dist)}
}
