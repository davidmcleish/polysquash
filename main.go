package main

import (
	"bytes"
	"log"

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
	}

	polys := []string{
		"POLYGON ((2000 3000, 2010 3006, 2007 3001, 1997 2995, 2000 3000))",
	}

	for _, pstr := range polys {
		log.Println(pstr)
		g, _ := geom.UnmarshalWKT(pstr)
		poly := g.MustAsPolygon()
		for _, c := range candidates {
			var buf bytes.Buffer
			if err := c.Encode(&buf, poly); err != nil {
				log.Printf("%v Encode: %v", c, err)
				continue
			}
			enc := buf.Bytes()
			dec, err := c.Decode(bytes.NewBuffer(enc))
			if err != nil {
				log.Printf("%v Decode: %v", c, err)
				continue
			}
			log.Printf("%v: %d bytes", c, len(enc))
			// TODO: compare dec to poly
			log.Printf("%v: %s", c, dec.AsText())
		}
	}
}
