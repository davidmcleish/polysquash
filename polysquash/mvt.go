package polysquash

import (
	"io"
	"io/ioutil"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/peterstace/simplefeatures/geom"
)

type MVT struct {
	Precision float64
}

func (m MVT) String() string { return "MVT_gz" }

func (m MVT) Encode(w io.Writer, poly geom.Polygon) error {
	pts := poly.DumpCoordinates()
	ls := make(orb.LineString, 0, pts.Length())
	for i := 0; i < pts.Length(); i++ {
		p := pts.GetXY(i)
		x := math.Round(p.X * m.Precision)
		y := math.Round(p.Y * m.Precision)
		ls = append(ls, orb.Point{x, y})
	}
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(ls))
	layers := mvt.NewLayers(map[string]*geojson.FeatureCollection{"p": fc})
	buf, err := mvt.MarshalGzipped(layers)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func (m MVT) Decode(r io.Reader) (*geom.Polygon, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	layers, err := mvt.UnmarshalGzipped(buf)
	if err != nil {
		return nil, err
	}
	ols := layers[0].Features[0].Geometry.(orb.LineString)
	coords := make([]float64, 0, len(ols)*2)
	for _, p := range ols {
		coords = append(coords, p.X()/m.Precision, p.Y()/m.Precision)
	}
	// Fix cumulative error to make sure the poly is closed
	coords[len(coords)-2], coords[len(coords)-1] = coords[0], coords[1]

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
