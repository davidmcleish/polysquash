package poly

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/davidmcleish/polysquash/huffman"
	"github.com/juju/errors"
)

type Point struct {
	X, Y float64
}

func (a Point) Sub(b Point) Point {
	return Point{b.X - a.X, b.Y - a.Y}
}

const epsilon = 1e-6

func Quantise(pts []Point) []string {
	scale := 0.0
	for i := 0; i < len(pts)-1; i++ {
		p := pts[i]
		q := pts[i+1]
		dx := p.X - q.X
		dy := p.Y - q.Y
		dist := dx*dx + dy*dy
		if dist > scale {
			scale = dist
		}
	}
	scale = math.Sqrt(scale)

	tokens := []string{
		fmt.Sprintf("%.6f", pts[0].X),
		fmt.Sprintf("%.6f", pts[0].Y),
		fmt.Sprintf("%.6f", scale),
	}
	for i := 0; i < len(pts)-1; i++ {
		g := Gradient(pts[i].Sub(pts[(i+1)]), scale)
		tokens = append(tokens, g...)
	}
	return tokens
}

func Gradient(p Point, scale float64) []string {
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
	dist := math.Sqrt(p.X*p.X+p.Y*p.Y) / scale
	log.Println(dir, grad, dist, p)

	result := []string{dir}
	result = append(result, huffman.Quantise(grad)...)
	result = append(result, huffman.Quantise(dist)...)
	return result
}

func FromGradient(tok []string) (Point, error) {
	dir := tok[0]
	grad, err := huffman.FromQuantised(tok[1], tok[2])
	if err != nil {
		return Point{}, errors.Trace(err)
	}
	dist, err := huffman.FromQuantised(tok[3], tok[4])
	if err != nil {
		return Point{}, errors.Trace(err)
	}
	p := Point{}
	p.X = math.Sqrt(dist * dist / (1 + 1/(grad*grad)))
	p.Y = math.Sqrt(dist * dist / (1 + grad*grad))
	if dir[0] == '-' {
		p.X = -p.X
	}
	if dir[2] == '-' {
		p.Y = -p.Y
	}
	if dir[1] == 'y' {
		p.Y, p.X = p.X, p.Y
	}
	log.Println(dir, grad, dist, p)
	return p, nil
}

func Parse(tokens []string) ([]Point, error) {
	var origin Point
	var err error
	var scale float64
	origin.X, err = strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return nil, errors.Trace(err)
	}
	origin.Y, err = strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return nil, errors.Trace(err)
	}
	scale, err = strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return nil, errors.Trace(err)
	}

	pts := []Point{origin}
	for i := 3; i+5 <= len(tokens); i += 5 {
		p, err := FromGradient(tokens[i : i+5])
		if err != nil {
			return nil, errors.Trace(err)
		}
		p.X = p.X*scale + pts[len(pts)-1].X
		p.Y = p.Y*scale + pts[len(pts)-1].Y
		pts = append(pts, p)
	}
	return pts, nil
}
