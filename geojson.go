package gegography

import (
	"encoding/json"
)

type geoJSONGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

type geoJSONFeature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Geometry   geoJSONGeometry        `json:"geometry"`
}

type geoJSON struct {
	Type     string           `json:"type"`
	Name     string           `json:"name,omitempty"`
	Features []geoJSONFeature `json:"features"`
}

type gjPoint []float64
type gjMultiPoint []gjPoint
type gjPolygon []gjMultiPoint
type gjMultiPolygon []gjPolygon

func (g gjPoint) toPoint() Point {
	return Point{X: g[0], Y: g[1]}
}

func (g gjMultiPoint) toMultiPoint() MultiPoint {
	mp := make(MultiPoint, 0)

	for x := range g {
		_g := g[x]
		mp = append(mp, _g.toPoint())
	}

	return mp
}

func (g gjPolygon) toPolygon() Polygon {
	p := make(Polygon, 0)

	for x := range g {
		_g := g[x]
		p = append(p, _g.toMultiPoint())
	}

	return p
}

func (g gjMultiPolygon) toMultiPolygon() MultiPolygon {
	mp := make(MultiPolygon, 0)

	for x := range g {
		_g := g[x]
		mp = append(mp, _g.toPolygon())
	}

	return mp
}

func (f *Feature) toGeoJSONFeature() (geoJSONFeature, error) {
	var coordinates interface{}

	switch f.Type {
	case "Point":
		coordinates = f.Coordinates.(Point).toGeoJSON()
	case "MultiPoint":
		coordinates = f.Coordinates.(MultiPoint).toGeoJSON()
	case "LineString":
		coordinates = f.Coordinates.(MultiPoint).toGeoJSON()
	case "Polygon":
		coordinates = f.Coordinates.(Polygon).toGeoJSON()
	case "MultiLineString":
		coordinates = f.Coordinates.(Polygon).toGeoJSON()
	case "MultiPolygon":
		coordinates = f.Coordinates.(MultiPolygon).toGeoJSON()
	default:
		return geoJSONFeature{}, GeoTypeError{Type: f.Type}
	}

	jc, err := json.Marshal(coordinates)

	if err != nil {
		return geoJSONFeature{}, err
	}

	gjf := geoJSONFeature{
		Type:       "Feature",
		Properties: f.Properties,
	}

	gjf.Geometry.Type = f.Type
	gjf.Geometry.Coordinates = jc

	return gjf, nil
}

// ToGeoJSON exports a Feature to a byte array containing JSON conforming to the GeoJSON format
func (f *Feature) ToGeoJSON() ([]byte, error) {
	gjf, err := f.toGeoJSONFeature()
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(gjf)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (p Point) toGeoJSON() gjPoint {
	return gjPoint{p.X, p.Y}
}

func (mp MultiPoint) toGeoJSON() gjMultiPoint {
	gjmp := make(gjMultiPoint, 0)

	for x := range mp {
		p := mp[x]
		gjmp = append(gjmp, p.toGeoJSON())
	}

	return gjmp
}

func (p Polygon) toGeoJSON() gjPolygon {
	gjp := make(gjPolygon, 0)

	for x := range p {
		mp := p[x]
		gjp = append(gjp, mp.toGeoJSON())
	}

	return gjp
}

func (mp MultiPolygon) toGeoJSON() gjMultiPolygon {
	gjmp := make(gjMultiPolygon, 0)

	for x := range mp {
		p := mp[x]
		gjmp = append(gjmp, p.toGeoJSON())
	}

	return gjmp
}

func (fc *FeatureCollection) toGeoJSONStruct() (geoJSON, error) {
	gj := geoJSON{Name: fc.Name, Type: "FeatureCollection"}

	for x := range fc.Features {
		f := fc.Features[x]

		gjf, err := f.toGeoJSONFeature()

		if err != nil {
			return geoJSON{}, err
		}

		gj.Features = append(gj.Features, gjf)
	}

	return gj, nil
}

// ToGeoJSON exports a FeatureCollection to a byte array containing JSON conforming to the GeoJSON format
func (fc *FeatureCollection) ToGeoJSON() ([]byte, error) {
	gj, err := fc.toGeoJSONStruct()

	if err != nil {
		return nil, err
	}

	return json.Marshal(gj)
}

// ToGeoJSONFeatureArray exports a FeatureCollection to a byte array conforming to the GeoJSON format but discards everything except the feature array
func (fc *FeatureCollection) ToGeoJSONFeatureArray() ([]byte, error) {
	gj, err := fc.toGeoJSONStruct()

	if err != nil {
		return nil, err
	}

	return json.Marshal(gj.Features)
}

// ToPrettyGeoJSON exports a FeatureCollection to a byte array containing indented JSON conforming to the GeoJSON format
func (fc *FeatureCollection) ToPrettyGeoJSON() ([]byte, error) {
	gj, err := fc.toGeoJSONStruct()

	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(gj, "", "\t")
}

// LoadGeoJSON parses an array of bytes conforming to the GeoJSON format to a FeatureCollection
func LoadGeoJSON(input []byte) (FeatureCollection, error) {
	var gj geoJSON

	if err := json.Unmarshal(input, &gj); err != nil {
		return FeatureCollection{}, err
	}

	var fc FeatureCollection
	fc.Features = make([]Feature, 0)
	fc.Name = gj.Name

	var coordinates interface{}
	var err error

	for f := range gj.Features {
		skip := false
		feature := gj.Features[f]
		switch feature.Geometry.Type {
		case "Point":
			var g gjPoint
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toPoint()
		case "MultiPoint":
			var g gjMultiPoint
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toMultiPoint()
		case "LineString":
			var g gjMultiPoint
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toMultiPoint()
		case "MultiLineString":
			var g gjPolygon
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toPolygon()
		case "Polygon":
			var g gjPolygon
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toPolygon()
		case "MultiPolygon":
			var g gjMultiPolygon
			err = json.Unmarshal(feature.Geometry.Coordinates, &g)
			coordinates = g.toMultiPolygon()
		case "":
			skip = true
		default:
			return FeatureCollection{}, GeoTypeError{Type: feature.Type}
		}

		if !skip {
			if err != nil {
				return FeatureCollection{}, GeoFormatError{Msg: "GeoJSON is malformed"}
			}

			fc.Features = append(fc.Features, Feature{Type: feature.Geometry.Type, Properties: feature.Properties, Coordinates: coordinates})
		}
	}

	return fc, nil
}
