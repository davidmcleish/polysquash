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
	candidates := []polysquash.EncoderDecoder{
		// polysquash.WKT{},
		// polysquash.Base64{Data: polysquash.Zip{Data: polysquash.WKT{}}},
		// polysquash.Base64{Data: polysquash.Gzip{Data: polysquash.WKT{}}},
		// polysquash.Base64{Data: polysquash.WKB{}},
		// polysquash.Base64{Data: polysquash.Gzip{Data: polysquash.WKB{}}},
		// polysquash.Base64{Data: polysquash.HuffmanWKT{}},
		// polysquash.Base64{Data: polysquash.Offset{Precision: 1 << 24}},
		polysquash.Base64{Data: polysquash.Huffman{Data: polysquash.Offset{Precision: 1 << 22}}},
		polysquash.Base64{Data: polysquash.Gzip{Data: polysquash.Offset{Precision: 1 << 22}}},
		// polysquash.Base64{Data: polysquash.Gradient{Precision: 1 << 24}},
		// polysquash.Base64{Data: polysquash.Geobuf{}},
		// polysquash.Base64{Data: polysquash.Gzip{Data: polysquash.Geobuf{}}},
		polysquash.Base64{Data: polysquash.MVT{Precision: 1 << 22}},
	}

	polys := []struct {
		name string
		pstr string
	}{
		{"square", "POLYGON ((2000 3000, 2010 3006, 2007 3001, 1997 2995, 2000 3000))"},
		{"small step", "POLYGON ((0 5, 0 6, 1 6, 1 7, 2 7, 2 8, 3 8, 3 9, 4 9, 4 10, 5 10, 5 9, 6 9, 6 8, 7 8, 7 7, 8 7, 9 6, 9 5, 10 5, 10 4, 9 4, 9 3, 8 3, 8 2, 7 2, 7 1, 6 1, 6 0, 5 0, 5 1, 4 1, 4 2, 3 2, 3 3, 2 3, 2 4, 1 4, 1 5, 0 5))"},
		{"big step", stepPoly(151.196, -33.865, 0.00002, 0.00003, 100).AsText()},
		{"big step with jitter", jitter(stepPoly(151.196, -33.865, 0.00002, 0.00003, 100), 1e-6).AsText()},
		{"star", starPoly(151.196, -33.865, 0.001, 0.001, 100).AsText()},
		{"circle", circlePoly(151.196, -33.865, 0.001, 100).AsText()},
	}

	for _, tc := range polys {
		log.Println(tc.name)
		// log.Println(tc.pstr)
		g, _ := geom.UnmarshalWKT(tc.pstr)
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
			diff, err := geom.SymmetricDifference(poly.AsGeometry(), dec.AsGeometry())
			if err != nil {
				log.Printf("%v\tSymmetric difference: %v", c, err)
			}
			log.Printf("%v\t%d bytes\tdiff %.4g", c, len(enc), diff.Area())
			if false {
				log.Printf("%v\t%s", c, dec.AsText())
			}
		}
	}
}

func circlePoly(ox, oy, radius float64, npoints int) geom.Polygon {
	return starPoly(ox, oy, radius, 0, npoints)
}

func starPoly(ox, oy, inRadius, outRadius float64, npoints int) geom.Polygon {
	r := rand.New(rand.NewSource(1234))
	var pts []float64
	for i := 0; i < npoints; i++ {
		angle := math.Pi * 2 * (float64(i) / float64(npoints))
		dist := inRadius + r.Float64()*outRadius
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

func jitter(inPoly geom.Polygon, amount float64) geom.Polygon {
	r := rand.New(rand.NewSource(2345))
	seq := inPoly.DumpCoordinates()
	pts := make([]float64, 0, seq.Length()*2)
	for i := 0; i < seq.Length(); i++ {
		p := seq.GetXY(i)
		pts = append(pts, p.X+r.NormFloat64()*amount)
		pts = append(pts, p.Y+r.NormFloat64()*amount)
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
