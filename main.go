package main

import (
	"bytes"
	"log"
	"math"
	"math/rand"

	"github.com/davidmcleish/polysquash/polysquash"
	"github.com/peterstace/simplefeatures/geom"
)

func main() {
	// points := []polysquash.Point{
	// 	{X: 2000, Y: 3000}, {X: 2010, Y: 3006}, {X: 2007, Y: 3001}, {X: 1997, Y: 2995},
	// }
	// tokens := polysquash.QuantisePts(points)
	// fmt.Println(tokens)

	// table := polysquash.CreateHuffman(tokens)
	// fmt.Println(table)

	// parsed, err := polysquash.Parse(tokens)
	// fmt.Println(parsed, err)

	candidates := []polysquash.EncoderDecoder{
		polysquash.WKT{},
		polysquash.Base64{Binary: polysquash.HuffmanWKT{}},
		polysquash.Base64{Binary: polysquash.Offset{Precision: 1 << 24}},
	}

	polys := []string{
		"POLYGON ((2000 3000, 2010 3006, 2007 3001, 1997 2995, 2000 3000))",
		"POLYGON ((0 5, 0 6, 1 6, 1 7, 2 7, 2 8, 3 8, 3 9, 4 9, 4 10, 5 10, 5 9, 6 9, 6 8, 7 8, 7 7, 8 7, 9 6, 9 5, 10 5, 10 4, 9 4, 9 3, 8 3, 8 2, 7 2, 7 1, 6 1, 6 0, 5 0, 5 1, 4 1, 4 2, 3 2, 3 3, 2 3, 2 4, 1 4, 1 5, 0 5))",
		stepPoly(151.196, -33.865, 0.00002, 0.00003, 1000).AsText(),
		randomStarPoly(-33.865, 151.196, 0.001, 1000).AsText(),
	}

	for i, pstr := range polys {
		log.Println(i)
		// log.Println(pstr)
		g, _ := geom.UnmarshalWKT(pstr)
		poly := g.MustAsPolygon()
		for _, c := range candidates {
			var buf bytes.Buffer
			if err := c.Encode(&buf, poly); err != nil {
				log.Printf("%v\tEncode: %v", c, err)
				continue
			}
			enc := buf.Bytes()
			dec, err := c.Decode(bytes.NewBuffer(enc))
			if err != nil {
				log.Printf("%v\tDecode: %v", c, err)
				continue
			}
			log.Printf("%v\t%d bytes", c, len(enc))
			// TODO: compare dec to poly
			if false {
				log.Printf("%v\t%s", c, dec.AsText())
			}
		}
	}
}

func randomStarPoly(ox, oy float64, radius float64, npoints int) geom.Polygon {
	r := rand.New(rand.NewSource(1234))
	var pts []float64
	for i := 0; i < npoints; i++ {
		angle := math.Pi * 2 * (float64(i) / float64(npoints))
		dist := (r.Float64()*0.5 + 0.5) * radius
		x := ox + math.Cos(angle)*dist
		y := oy + math.Sin(angle)*dist
		pts = append(pts, x, y)
	}
	// Close the polygon
	pts = append(pts, pts[0], pts[1])
	ls, err := geom.NewLineString(geom.NewSequence(pts, geom.DimXY))
	if err != nil {
		panic(err)
	}
	poly, err := geom.NewPolygon([]geom.LineString{ls})
	if err != nil {
		panic(err)
	}
	return poly
}

func stepPoly(ox, oy, dx, dy float64, npoints int) geom.Polygon {
	var pts []float64
	x := ox
	y := oy
	for dir := 0; dir < 4; dir++ {
		for i := 0; i < npoints; i++ {
			x += dx
			y += dy
			pts = append(pts, x, y)
			x -= dy
			y += dx
			pts = append(pts, x, y)
		}
		dx, dy = -dy, dx
	}
	// Close the polygon
	pts = append(pts, pts[0], pts[1])
	ls, err := geom.NewLineString(geom.NewSequence(pts, geom.DimXY))
	if err != nil {
		panic(err)
	}
	poly, err := geom.NewPolygon([]geom.LineString{ls})
	if err != nil {
		panic(err)
	}
	return poly
}
