package polysquash

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/peterstace/simplefeatures/geom"
)

//import "github.com/dgryski/go-bitstream"

type EncoderDecoder interface {
	Encode(w io.Writer, poly geom.Polygon) error
	Decode(r io.Reader) (*geom.Polygon, error)
}

type WKT struct{}

func (t WKT) String() string { return "WKT" }

func (t WKT) Encode(w io.Writer, poly geom.Polygon) error {
	var wkt []byte
	wkt = poly.AppendWKT(wkt)
	_, err := w.Write(wkt)
	return err
}
func (t WKT) Decode(r io.Reader) (*geom.Polygon, error) {
	wkt, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	g, err := geom.UnmarshalWKT(string(wkt))
	if err != nil {
		return nil, err
	}
	poly, ok := g.AsPolygon()
	if !ok {
		return nil, errors.New("not a polygon George")
	}
	return &poly, nil
}
