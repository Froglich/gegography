package gegography

import "fmt"

// Point describes a set of coordinates
type Point struct {
	X float64
	Y float64
}

// MultiPoint describes a collection of points
type MultiPoint []Point

// LineString describes a collection of points which together form a line
type LineString []Point

// Polygon describes a collection of point collections which together form a polygon
type Polygon []MultiPoint

// MultiPolygon describes a collection of Polygons
type MultiPolygon []Polygon

// CRS describes a coordinate reference system according to the GeoJSON format specification
type CRS struct {
	Type       string `json:"type"`
	Properties struct {
		Name string `json:"name,omitempty"`
		Href string `json:"href,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"properties"`
}

// Feature represents a geographical feature
type Feature struct {
	Type        string
	Properties  map[string]interface{}
	Coordinates interface{}
}

// FeatureCollection represents a collection of geographical features and accompanying information
type FeatureCollection struct {
	Name                      string
	CoordinateReferenceSystem *CRS
	Features                  []Feature
}

// GeoTypeError describes an error involving an unsupported geographical type
type GeoTypeError struct {
	Type string
}

func (g GeoTypeError) Error() string {
	return fmt.Sprintf("%s: bad or unsupported geographical type.", g.Type)
}

// GeoFormatError describes an error involving badly formatted geographical information
type GeoFormatError struct {
	Msg string
}

func (g GeoFormatError) Error() string {
	return g.Msg
}

// NewFeatureCollection returns a new blank FeatureCollection with an instantiated feature array
func NewFeatureCollection() FeatureCollection {
	fc := FeatureCollection{}
	fc.Features = make([]Feature, 0)

	return fc
}

// AddFeature adds a feature to a FeatureCollection
func (fc *FeatureCollection) AddFeature(f Feature) {
	fc.Features = append(fc.Features, f)
}
