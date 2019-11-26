package gegography

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"strings"
	"encoding/binary"
)

const shpMin float64 = -1e38

func parseShpPoint(in []byte) (Point, error) {
	if len(in) != 16 {
		return Point{}, GeoFormatError{Msg: "point has wrong number of bytes"}
	}

	var x float64
	var y float64

	err := parseValue(in[0:8], binary.LittleEndian, &x)
	if err != nil {
		return Point{}, err
	}

	err = parseValue(in[8:16], binary.LittleEndian, &y)
	if err != nil {
		return Point{}, err
	}

	return Point{X: x, Y: y}, nil
}

func parseShpMultiPoint(in []byte) (MultiPoint, error) {
	if len(in) < 52 {
		return nil, GeoFormatError{Msg: "multipoint with too few bytes"}
	}

	mp := make(MultiPoint, 0)

	var nr int32
	err := parseValue(in[32:36], binary.LittleEndian, &nr)
	if err != nil {
		return nil, err
	}

	if len(in) < 36+int(nr)*16 {
		return nil, GeoFormatError{Msg: "multipoint is malformed"}
	}

	for z := 0; z < int(nr); z++ {
		s := 36 + z*16

		p, err := parseShpPoint(in[s:s+16])
		if err != nil {
			return nil, err
		}

		mp = append(mp, p)
	}

	return mp, err
}

func parseShpPolyLine(in []byte) (Polygon, error) {
	if len(in) < 44 {
		return nil, GeoFormatError{Msg: "polyline with too few bytes"}
	}

	p := make(Polygon, 0)

	var nparts int32
	err := parseValue(in[32:36], binary.LittleEndian, &nparts)
	if err != nil {
		return nil, err
	}

	var npoints int32
	err = parseValue(in[36:40], binary.LittleEndian, &npoints)
	if err != nil {
		return nil, err
	}

	parts := make([]int32, int(nparts))
	points := make([]Point, int(npoints))

	var v int32
	for x := 0; x < int(nparts); x++ {
		s := 40 + x*4
		err = parseValue(in[s:s+4], binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		parts[x] = v
	}

	for x := 0; x < int(npoints); x++ {
		s := 40 + 4*int(nparts) + 16*x
		point, err := parseShpPoint(in[s:s+16])
		if err != nil {
			return nil, err
		}
		points[x] = point
	}

	i := len(parts) - 1
	for x := 0; x < i; x++ {
		p = append(p, points[parts[x]:parts[x+1]])
	}

	p = append(p, points[parts[i]:])

	return p, nil
}

func parseValue(in []byte, order binary.ByteOrder, out interface{}) error {
	buf := bytes.NewReader(in)
	return binary.Read(buf, order, out)
}

type shpGeoParser struct {
	Features []Feature
	Error error
}

func shpGeographyParser(filename string, result chan shpGeoParser) {
	var r io.Reader
	r, err := os.Open(filename)

	if err != nil {
		result <- shpGeoParser{Error: err}
	}

	header := make([]byte, 100)

	_, err = r.Read(header)
	if err != nil {
		result <- shpGeoParser{Error: err}
	}

	var fileLength int32
	err = parseValue(header[24:28], binary.BigEndian, &fileLength)
	if err != nil {
		result <- shpGeoParser{Error: err}
	}

	var pos int32 = 100
	var t int32
	rh := make([]byte, 8)

	var cl int32

	features := make([]Feature, 0)

	for pos < fileLength*2 {
		_, err = r.Read(rh)
		if err != nil {
			result <- shpGeoParser{Error: err}
		}
		pos += 8

		err = parseValue(rh[4:8], binary.BigEndian, &cl)
		if err != nil {
			result <- shpGeoParser{Error: err}
		}

		content := make([]byte, cl*2)
		_, err = r.Read(content)
		if err != nil {
			result <- shpGeoParser{Error: err}
		}
		pos += cl*2

		err = parseValue(content[0:4], binary.LittleEndian, &t)
		if err != nil {
			result <- shpGeoParser{Error: err}
		}

		var c interface{}

		switch(t) {
			case 1:
				c, err = parseShpPoint(content[4:])
				features = append(features, Feature{Type: "Point", Coordinates: c})
			case 8:
				c, err = parseShpMultiPoint(content[4:])
				features = append(features, Feature{Type: "MultiPoint", Coordinates: c})
			case 3:
				c, err = parseShpPolyLine(content[4:])
				features = append(features, Feature{Type: "MultiLineString", Coordinates: c})
			case 5:
				c, err = parseShpPolyLine(content[4:])
				features = append(features, Feature{Type: "Polygon", Coordinates: c})
			default:
				result <- shpGeoParser{Error: GeoTypeError{Type: fmt.Sprintf("unsupported shapefile geographical type '%v'", t)}}
		}

		if err != nil {
			result <- shpGeoParser{Error: err}
		}
	}

	result <- shpGeoParser{Features: features}
}

func ParseShapefile(filename string) {
	geoChan := make(chan shpGeoParser)

	tabFile := strings.Replace(filename, ".shp", ".dbf", 1)

	go shpGeographyParser(filename, geoChan)

	geo := <- geoChan

	fmt.Println(geo.Features)
}
