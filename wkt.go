package gegography

import (
	"fmt"
	"strings"
	"strconv"
)

func (p Point) toWKT() string {
	return fmt.Sprintf("%f %f", p.X, p.Y)
}

func (mp MultiPoint) toWKT() string {
	str := make([]string, 0)

	for x := range mp {
		p := mp[x]
		str = append(str, p.toWKT())
	}

	return strings.Join(str, ", ")
}

func (p Polygon) toWKT() string {
	str := make([]string, 0)

	for x := range p {
		mp := p[x]
		str = append(str, fmt.Sprintf("(%s)", mp.toWKT()))
	}

	return strings.Join(str, ", ")
}

func (mp MultiPolygon) toWKT() string {
	str := make([]string, 0)

	for x := range mp {
		p := mp[x]
		str = append(str, fmt.Sprintf("(%s)", p.toWKT()))
	}

	return strings.Join(str, ", ")
}

//ToWKT writes a WKT string representing a Feature
func (f *Feature) ToWKT() (string, error) {
	var str string

	switch(f.Type) {
		case "Point":
			str = f.Coordinates.(Point).toWKT()
		case "MultiPoint":
			str = f.Coordinates.(MultiPoint).toWKT()
		case "LineString":
			str = f.Coordinates.(MultiPoint).toWKT()
		case "Polygon":
			str = f.Coordinates.(Polygon).toWKT()
		case "MultiLineString":
			str = f.Coordinates.(Polygon).toWKT()
		case "MultiPolygon":
			str = f.Coordinates.(MultiPolygon).toWKT()
		default:
			return "", GeoTypeError{Type: f.Type}
	}

	return fmt.Sprintf("%s (%s)", strings.ToUpper(f.Type), str), nil
}

//ToWKT writes a WKT string representing a FeatureCollection
func (fc *FeatureCollection) ToWKT() (string, error) {
	str := make([]string, 0)

	for x := range fc.Features {
		f := fc.Features[x]
		wkt, err := f.ToWKT()

		if err != nil {
			return "", err
		}

		str = append(str, wkt)
	}

	return fmt.Sprintf("GEOMETRYCOLLECTION (%s)", strings.Join(str, ", ")), nil
}

func parseWKTPoint(wkt string) (Point, error) {
	var p Point

	w := strings.Trim(wkt, " ()\t\r\n")

	parts := strings.Split(w, " ")
	length := len(parts)

	if length < 2 || length > 3 {
		return p, GeoFormatError{Msg: fmt.Sprintf("invalid WKT point - '%s' wrong number of separators", w)}
	}

	x, err := strconv.ParseFloat(parts[0], 64)

	if err != nil {
		return p, err
	}

	y, err := strconv.ParseFloat(parts[1], 64)

	if err != nil {
		return p, err
	}

	p.X = x
	p.Y = y

	if length == 3 {
		z, err := strconv.ParseFloat(parts[2], 64)

		if err != nil {
			return p, err
		}

		p.Z = z
	}

	return p, nil
}

func parseWKTMultiPoint(wkt string) (MultiPoint, error) {
	mp := make(MultiPoint, 0)

	parts := strings.Split(wkt, ",")

	for x := range parts {
		p, err := parseWKTPoint(parts[x])

		if err != nil {
			return nil, GeoFormatError{Msg: fmt.Sprintf("bad WKT component - '%s', %v", wkt, err)}
		}

		mp = append(mp, p)
	}

	return mp, nil
}

func parseWKTPolygon(wkt string) (Polygon, error) {
	p := make(Polygon, 0)

	parts := strings.Split(wkt, "),(")

	for x := range parts {
		mp, err := parseWKTMultiPoint(parts[x])

		if err != nil {
			return nil, GeoFormatError{Msg: fmt.Sprintf("bad WKT component - '%s', %v", wkt, err)}
		}

		p = append(p, mp)
	}

	return p, nil
}

func parseWKTMultiPolygon(wkt string) (MultiPolygon, error) {
	mp := make(MultiPolygon, 0)

	parts := strings.Split(wkt, ")),((")

	for x := range parts {
		p, err := parseWKTPolygon(parts[x])

		if err != nil {
			return nil, GeoFormatError{Msg: fmt.Sprintf("bad WKT component - '%s', %v", wkt, err)}
		}

		mp = append(mp, p)
	}

	return mp, nil
}

//ParseWKT parses a WKT string and returns a feature
func ParseWKT(wkt string) (Feature, error) {
	var g interface{}
	var t string
	var err error

	w := strings.ToUpper(wkt)

	w = strings.Replace(w, ", ", ",", -1)
	w = strings.Replace(w, " ,", ",", -1)

	s := strings.Index(w, "(")
	e := strings.LastIndex(w, ")")

	if s < 6 || e < 7 || s > e {
		return Feature{}, GeoFormatError{Msg: fmt.Sprintf("invalid WKT - '%s' not properly enclosed", w)}
	}

	if strings.HasPrefix(w, "POINT") {
		g, err = parseWKTPoint(w[s+1:e])
		t = "Point"
	} else if strings.HasPrefix(w, "MULTIPOINT") {
		g, err = parseWKTMultiPoint(w[s+1:e])
		t = "MultiPoint"
	} else if strings.HasPrefix(w, "LINESTRING") {
		g, err = parseWKTMultiPoint(w[s+1:e])
		t = "LineString"
	} else if strings.HasPrefix(w, "POLYGON") {
		g, err = parseWKTPolygon(w[s+1:e])
		t = "Polygon"
	} else if strings.HasPrefix(w, "MULTILINESTRING") {
		g, err = parseWKTPolygon(w[s+1:e])
		t = "MultiLineString"
	} else if strings.HasPrefix(w, "MULTIPOLYGON") {
		g, err = parseWKTMultiPolygon(w[s+1:e])
		t = "MultiPolygon"
	} else {
		err = GeoTypeError{Type: w}
	}

	if err != nil {
		return Feature{}, err
 	}

 	return Feature{Type: t, Coordinates: g, Properties: make(map[string]interface{})}, nil
}
