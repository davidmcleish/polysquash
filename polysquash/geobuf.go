package polysquash

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/cairnapp/go-geobuf"
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
	geopb "github.com/cairnapp/go-geobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/peterstace/simplefeatures/geom"
)

type Geobuf struct{}

func (g Geobuf) String() string { return "Geobuf" }

func (g Geobuf) Encode(w io.Writer, poly geom.Polygon) error {
	coords := poly.DumpCoordinates()
	pts := make([]geometry.Point, 0, coords.Length())
	for i := 0; i < coords.Length(); i++ {
		p := coords.GetXY(i)
		pts = append(pts, geometry.Point([]float64{p.X, p.Y}))
	}
	pb := geobuf.Encode(geojson.NewGeometry(geometry.LineString(pts)))
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (g Geobuf) Decode(r io.Reader) (*geom.Polygon, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var pb geopb.Data
	err = proto.Unmarshal(data, &pb)
	if err != nil {
		return nil, err
	}
	dg := geobuf.Decode(&pb).(*geojson.Geometry)
	coords := dg.Coordinates

	// ugh
	cstr := fmt.Sprint(coords)
	var pts []float64
	for _, tok := range regexp.MustCompile(`[ \[\]]`).Split(cstr, -1) {
		if tok == "" {
			continue
		}
		v, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return nil, err
		}
		pts = append(pts, v)
	}

	ls, err := geom.NewLineString(geom.NewSequence(pts, geom.DimXY))
	if err != nil {
		return nil, err
	}
	poly, err := geom.NewPolygon([]geom.LineString{ls})
	if err != nil {
		return nil, err
	}
	return &poly, nil
}
