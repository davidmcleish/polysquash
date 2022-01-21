package polysquash

import (
	"archive/zip"
	"bytes"
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

type WKB struct{}

func (b WKB) String() string { return "WKB" }

func (b WKB) Encode(w io.Writer, poly geom.Polygon) error {
	_, err := w.Write(poly.AsBinary())
	return err
}

func (b WKB) Decode(r io.Reader) (*geom.Polygon, error) {
	wkb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	g, err := geom.UnmarshalWKB(wkb)
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
	Data EncoderDecoder
}

func (b Base64) String() string { return b.Data.String() + "_b64" }

func (b Base64) Encode(w io.Writer, poly geom.Polygon) error {
	enc := base64.NewEncoder(base64.URLEncoding, w)
	if err := b.Data.Encode(enc, poly); err != nil {
		return err
	}
	return enc.Close()
}

func (b Base64) Decode(r io.Reader) (*geom.Polygon, error) {
	dec := base64.NewDecoder(base64.URLEncoding, r)
	return b.Data.Decode(dec)
}

type Zip struct {
	Data EncoderDecoder
}

func (z Zip) String() string { return z.Data.String() + "_z" }

func (z Zip) Encode(w io.Writer, poly geom.Polygon) error {
	zw := zip.NewWriter(w)
	zf, err := zw.Create("poly")
	if err != nil {
		return err
	}
	if err := z.Data.Encode(zf, poly); err != nil {
		return err
	}
	return zw.Close()
}

func (z Zip) Decode(r io.Reader) (*geom.Polygon, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	zf, err := zr.Open("poly")
	if err != nil {
		return nil, err
	}
	return z.Data.Decode(zf)
}
