package polysquash

import (
	"fmt"
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
	"github.com/peterstace/simplefeatures/geom"
)

type Gradient struct {
	Precision float64
}

func (g Gradient) String() string { return "Gradnt" }

func (g Gradient) Encode(w io.Writer, poly geom.Polygon) error {
	pts := poly.DumpCoordinates()
	tokens := make([]int64, 0, pts.Length()*3-1)

	p0 := pts.Get(0)
	x := int64(math.Round(p0.X * g.Precision))
	y := int64(math.Round(p0.Y * g.Precision))
	tokens = append(tokens, x, y)
	prevX := float64(x) / g.Precision
	prevY := float64(y) / g.Precision
	var prevGrad, prevDist int64

	for i := 1; i < pts.Length(); i++ {
		p := pts.GetXY(i)
		dir, grad, dist := g.calcGradient(p.X-prevX, p.Y-prevY)
		tokens = append(tokens, dir, grad-prevGrad, dist-prevDist)
		prevX, prevY = g.addGradient(prevX, prevY, dir, grad, dist)
		prevGrad = grad
		prevDist = dist
	}

	fmt.Println(tokens)

	bw := bitstream.NewWriter(w)
	if err := HuffmanEncode(bw, tokens); err != nil {
		return err
	}
	return bw.Flush(bitstream.Zero)
}

func (g Gradient) Decode(r io.Reader) (*geom.Polygon, error) {
	br := bitstream.NewReader(r)
	tokens, err := HuffmanDecode(br)
	if err != nil {
		return nil, err
	}

	coords := make([]float64, 0, (len(tokens)+1)/3)

	x0 := float64(tokens[0]) / g.Precision
	y0 := float64(tokens[1]) / g.Precision
	coords = append(coords, x0, y0)
	prevX := x0
	prevY := y0
	var prevGrad, prevDist int64

	for i := 2; i+3 <= len(tokens); i += 3 {
		dir := tokens[i]
		grad := tokens[i+1] + prevGrad
		dist := tokens[i+2] + prevDist
		x, y := g.addGradient(prevX, prevY, dir, grad, dist)
		coords = append(coords, x, y)
		prevX = x
		prevY = y
		prevGrad = grad
		prevDist = dist
	}

	// fmt.Println(coords)
	// Fix cumulative error to make sure the poly is closed
	coords[len(coords)-2], coords[len(coords)-1] = x0, y0

	ls, err := geom.NewLineString(geom.NewSequence(coords, geom.DimXY))
	if err != nil {
		return nil, err
	}
	poly, err := geom.NewPolygon([]geom.LineString{ls})
	if err != nil {
		return nil, err
	}
	return &poly, nil
}

// type Point struct {
// 	X, Y float64
// }

// func (a Point) Sub(b Point) Point {
// 	return Point{b.X - a.X, b.Y - a.Y}
// }

// const epsilon = 1e-6

// func QuantisePts(pts []Point) []string {
// 	scale := 0.0
// 	for i := 0; i < len(pts)-1; i++ {
// 		p := pts[i]
// 		q := pts[i+1]
// 		dx := p.X - q.X
// 		dy := p.Y - q.Y
// 		dist := dx*dx + dy*dy
// 		if dist > scale {
// 			scale = dist
// 		}
// 	}
// 	scale = math.Sqrt(scale)

// 	tokens := []string{
// 		fmt.Sprintf("%.6f", pts[0].X),
// 		fmt.Sprintf("%.6f", pts[0].Y),
// 		fmt.Sprintf("%.6f", scale),
// 	}
// 	for i := 0; i < len(pts)-1; i++ {
// 		g := Gradient(pts[i].Sub(pts[(i+1)]), scale)
// 		tokens = append(tokens, g...)
// 	}
// 	return tokens
// }

func (g *Gradient) calcGradient(dx, dy float64) (int64, int64, int64) {
	ax := math.Abs(dx)
	ay := math.Abs(dy)
	if ax*g.Precision < 1 && ay*g.Precision < 1 {
		// ...or just skip this point
		return 0, 0, 0
	}
	var dir int64
	if dx < 0 {
		dir |= 1
	}
	if dy < 0 {
		dir |= 2
	}
	var grad float64
	if ax < ay {
		dir |= 4
		grad = ax / ay
	} else {
		grad = ay / ax
	}
	qgrad := int64(math.Round(grad * g.Precision))
	qdist := int64(math.Round(math.Sqrt(dx*dx+dy*dy) * g.Precision))
	return dir, qgrad, qdist
}

func (g *Gradient) addGradient(x, y float64, dir, grad, dist int64) (float64, float64) {
	if grad == 0 {
		return x, y
	}
	gf := float64(grad) / g.Precision
	df := float64(dist) / g.Precision
	dx := math.Sqrt(df * df / (1 + 1/(gf*gf)))
	dy := math.Sqrt(df * df / (1 + gf*gf))
	if dir&4 == 0 {
		dx, dy = dy, dx
	}
	if dir&1 != 0 {
		dx = -dx
	}
	if dir&2 != 0 {
		dy = -dy
	}
	return x + dx, y + dy
}

// func Parse(tokens []string) ([]Point, error) {
// 	var origin Point
// 	var err error
// 	var scale float64
// 	origin.X, err = strconv.ParseFloat(tokens[0], 64)
// 	if err != nil {
// 		return nil, errors.Trace(err)
// 	}
// 	origin.Y, err = strconv.ParseFloat(tokens[1], 64)
// 	if err != nil {
// 		return nil, errors.Trace(err)
// 	}
// 	scale, err = strconv.ParseFloat(tokens[2], 64)
// 	if err != nil {
// 		return nil, errors.Trace(err)
// 	}

// 	pts := []Point{origin}
// 	for i := 3; i+5 <= len(tokens); i += 5 {
// 		p, err := FromGradient(tokens[i : i+5])
// 		if err != nil {
// 			return nil, errors.Trace(err)
// 		}
// 		p.X = p.X*scale + pts[len(pts)-1].X
// 		p.Y = p.Y*scale + pts[len(pts)-1].Y
// 		pts = append(pts, p)
// 	}
// 	return pts, nil
// }
