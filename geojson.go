package gegography

import (
	"fmt"
	"encoding/json"
)

type geoJSONGeometry struct {
	Type            string          `json:"type"`
	Coordinates     json.RawMessage `json:"coordinates"`
}

type geoJSONFeature struct {
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Geometry    geoJSONGeometry        `json:"geometry"`
}

type geoJSON struct {
	Type     string           `json:"type"`
	Name     string           `json:"name,omitempty"`
	CRS      struct {
		Type string           `json:"type"`
		Properties struct {
			Name string       `json:"name"`
		}                     `json:"properties"`
	}                         `json:"crs"`
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

func (f *Feature) toGeoJSON() (geoJSONFeature, error) {
	var coordinates interface{}

	switch(f.Type) {
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

	gjf := geoJSONFeature {
		Type: "Feature",
		Properties: f.Properties,
	}

	gjf.Geometry.Type = f.Type
	gjf.Geometry.Coordinates = jc

	return gjf, nil
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

//ToGeoJSON writes a byte array containing JSON conforming to the GeoJSON format
func (fc *FeatureCollection) ToGeoJSON() ([]byte, error) {
	gj := geoJSON{Name: fc.Name, Type: "FeatureCollection"}

	if fc.CRS != "" {
		gj.CRS.Type = "name"
		gj.CRS.Properties.Name = fc.CRS
	}


	for x := range fc.Features {
		f := fc.Features[x]

		gjf, err := f.toGeoJSON()

		if err != nil {
			return nil, err
		}

		gj.Features = append(gj.Features, gjf)
	}

	return json.Marshal(gj)
}

//LoadGeoJSON parses an array of bytes conforming to the GeoJSON format to a FeatureCollection
func LoadGeoJSON(input []byte) (FeatureCollection, error) {
	var gj geoJSON

	if err := json.Unmarshal(input, &gj); err != nil {
		return FeatureCollection{}, err
	}

	var fc FeatureCollection
	fc.Features = make([]Feature, 0)
	fc.CRS = gj.CRS.Properties.Name
	fc.Name = gj.Name

	for f := range gj.Features {
		feature := gj.Features[f]
		switch feature.Geometry.Type {
			case "Point":
				var g gjPoint
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v point geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "Point", Properties: feature.Properties, Coordinates: g.toPoint()})
			case "MultiPoint":
				var g gjMultiPoint
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v multipoint geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "MultiPoint", Properties: feature.Properties, Coordinates: g.toMultiPoint()})
			case "LineString":
				var g gjMultiPoint
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v linestring geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "LineString", Properties: feature.Properties, Coordinates: g.toMultiPoint()})
			case "MultiLineString":
				var g gjPolygon
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v multilinestring geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "MultiLineString", Properties: feature.Properties, Coordinates: g.toPolygon()})
			case "Polygon":
				var g gjPolygon
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v polygon geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "Polygon", Properties: feature.Properties, Coordinates: g.toPolygon()})
			case "MultiPolygon":
				var g gjMultiPolygon
				if err := json.Unmarshal(feature.Geometry.Coordinates, &g); err != nil {
					return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("geojson feature %v multipolygon geometry coordinates malformed", f)}
				}
				fc.Features = append(fc.Features, Feature{Type: "MultiPolygon", Properties: feature.Properties, Coordinates: g.toMultiPolygon()})
			case "":
			default:
				return FeatureCollection{}, GeoTypeError{Type: feature.Type}
		}
	}

	return fc, nil
}
