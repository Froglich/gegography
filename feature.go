/*
 * This file is part of Gegography. Copyright 2019 © Kim Lindgren,
 * Unit for Field based forest research - Swedish University for
 * Agricultural sciences
 *
 * Gegography is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * Gegography is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Gegography.  If not, see <https://www.gnu.org/licenses/>.
 */

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
