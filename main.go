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
		polysquash.Base64{Binary: polysquash.HuffmanWKT{}},
	}

	polys := []string{
		"POLYGON ((2000 3000, 2010 3006, 2007 3001, 1997 2995, 2000 3000))",
		"POLYGON ((0 5, 0 6, 1 6, 1 7, 2 7, 2 8, 3 8, 3 9, 4 9, 4 10, 5 10, 5 9, 6 9, 6 8, 7 8, 7 7, 8 7, 9 6, 9 5, 10 5, 10 4, 9 4, 9 3, 8 3, 8 2, 7 2, 7 1, 6 1, 6 0, 5 0, 5 1, 4 1, 4 2, 3 2, 3 3, 2 3, 2 4, 1 4, 1 5, 0 5))",
	}

	for _, pstr := range polys {
		log.Println(pstr)
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
			log.Printf("%v\t%s", c, dec.AsText())
		}
	}
}
