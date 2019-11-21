package gegography

import "fmt"

type Point struct {
	X float64
	Y float64
}

type MultiPoint []Point
type LineString []Point
type Polygon []MultiPoint
type MultiLineString []LineString
type MultiPolygon []Polygon

type Feature struct {
	Type        string
	Properties  map[string]interface{}
	Coordinates interface{}
}

type FeatureCollection struct {
	Name     string
	CRS      string
	Features []Feature
}

type GeoTypeError struct {
	Type string
}

func (g GeoTypeError) Error() string {
	return fmt.Sprintf("%s: bad or unsupported geographical type.", g.Type)
}

type GeoFormatError struct {
	Msg string
}

func (g GeoFormatError) Error() string {
	return g.Msg
}

func NewFeatureCollection() FeatureCollection {
	fc := FeatureCollection{}
	fc.Features = make([]Feature, 0)

	return fc
}

func (fc *FeatureCollection) AddFeature(f Feature) {
	fc.Features = append(fc.Features, f)
}
