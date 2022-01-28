package polysquash

import (
	"bufio"
	"io"
	"math"

	"github.com/peterstace/simplefeatures/geom"
)

type Offset struct {
	Precision float64
}

func (o Offset) String() string { return "Offset" }

func (o Offset) Encode(w io.Writer, poly geom.Polygon) error {
	pts := poly.DumpCoordinates()
	tw := TokenWriter{w}

	var prevX, prevY int64
	for i := 0; i < pts.Length(); i++ {
		p := pts.GetXY(i)
		x := int64(math.Round(p.X * o.Precision))
		y := int64(math.Round(p.Y * o.Precision))
		if err := tw.WriteTokens(x-prevX, y-prevY); err != nil {
			return err
		}
		prevX = x
		prevY = y
	}
	return nil
}

func (o Offset) Decode(r io.Reader) (*geom.Polygon, error) {
	tr := TokenReader{bufio.NewReader(r)}
	var coords []float64
	var prevX, prevY int64

	for {
		xt, err := tr.ReadToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		yt, err := tr.ReadToken()
		if err != nil {
			return nil, err
		}
		x := float64(xt+prevX) / o.Precision
		y := float64(yt+prevY) / o.Precision
		coords = append(coords, x, y)
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
