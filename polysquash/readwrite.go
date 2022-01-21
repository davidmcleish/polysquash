package polysquash

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"

	"github.com/peterstace/simplefeatures/geom"
)

type EncoderDecoder interface {
	Encode(w io.Writer, poly geom.Polygon) error
	Decode(r io.Reader) (*geom.Polygon, error)
	String() string
}

type WKT struct{}

func (t WKT) String() string { return "WKT " }

func (t WKT) Encode(w io.Writer, poly geom.Polygon) error {
	_, err := w.Write([]byte(poly.AsText()))
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

type Base64 struct {
	Binary EncoderDecoder
}

func (b Base64) String() string { return b.Binary.String() + "_b64" }

func (b Base64) Encode(w io.Writer, poly geom.Polygon) error {
	enc := base64.NewEncoder(base64.URLEncoding, w)
	if err := b.Binary.Encode(enc, poly); err != nil {
		return err
	}
	return enc.Close()
}

func (b Base64) Decode(r io.Reader) (*geom.Polygon, error) {
	dec := base64.NewDecoder(base64.URLEncoding, r)
	return b.Binary.Decode(dec)
}
