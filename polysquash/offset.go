package polysquash

import (
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
	"github.com/peterstace/simplefeatures/geom"
)

type Offset struct {
	Precision float64
}

func (o Offset) String() string { return "Offset" }

func (o Offset) Encode(w io.Writer, poly geom.Polygon) error {
	pts := poly.DumpCoordinates()
	tokens := make([]int64, pts.Length()*2)

	var prevX, prevY int64
	for i := 0; i < pts.Length(); i++ {
		p := pts.GetXY(i)
		x := int64(math.Round(p.X * o.Precision))
		y := int64(math.Round(p.Y * o.Precision))
		tokens[i*2] = x - prevX
		tokens[i*2+1] = y - prevY
		prevX = x
		prevY = y
	}

	bw := bitstream.NewWriter(w)
	if err := HuffmanEncode(bw, tokens); err != nil {
		return err
	}
	return bw.Flush(bitstream.Zero)
}

func (o Offset) Decode(r io.Reader) (*geom.Polygon, error) {
	br := bitstream.NewReader(r)
	tokens, err := HuffmanDecode(br)
	if err != nil {
		return nil, err
	}

	coords := make([]float64, len(tokens))
	var prevX, prevY int64

	for i := 0; i+2 <= len(tokens); i += 2 {
		xt := tokens[i]
		yt := tokens[i+1]
		x := float64(xt+prevX) / o.Precision
		y := float64(yt+prevY) / o.Precision
		coords[i] = x
		coords[i+1] = y
		prevX += xt
		prevY += yt
	}

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
