<p align="center">
	<a href="https://www.gnu.org/licenses/gpl-3.0.en.html"><img src="https://img.shields.io/badge/License-GPL3-orange.svg" alt="GPL v3"/></a>
	<a href="https://goreportcard.com/report/github.com/Froglich/gegography"><img src="https://goreportcard.com/badge/github.com/Froglich/gegography" alt="Go report card"/></a>
	<a href="https://www.codefactor.io/repository/github/froglich/gegography"><img src="https://www.codefactor.io/repository/github/froglich/gegography/badge" alt="Codefactor"/></a>
	<a href="https://godoc.org/github.com/Froglich/gegography"><img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="GoDoc link"/></a>
</p>

# Gegography
Gegography is a library for reading, manipulating and converting
geographical formats, written i pure Go. Currently, support is limited to
GeoJSON, WKT (Well-Known-Text) and Shapefiles (read only). I plan on
implementing WKB (Well-Known-Binary).

I do not plan on supporting writing Shapefiles. [Shapefile must die!](http://switchfromshapefile.org/)

Gegography supports the following geographical types:
* Point
* LineString
* Polygon
* MultiPoint
* MultiLineString
* MultiPolygon

Currently, only XY geometries are supported.

## Installing
Install with `go get github.com/Froglich/gegography` .
