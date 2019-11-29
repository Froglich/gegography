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

## What is it for?
I wrote Gegography primarily as a utility for processing geographical data
uploaded to websites with Go backends. For example, I manage several websites
with interactive mapping capabilities, and a common request is for users to
be able to upload Shapefiles or GeoJSON-documents and have them automatically
displayed on the map, or saved in the underlying database for viewing on demand.

With gegography you can read a Shapefile or a GeoJSON-document and easily export
all features to WKT-strings for storage as geometries inside a database engine
such as PostgreSQL or SQL Server.

## How do I use it?
A typical use-case would look something like this (convert a shapefile to geojson)

```go
package main

import (
	"os"
	"github.com/Froglich/gegography"
	"io/ioutil"
)

func main() {
	infile := os.Args[1]
	fc, err := gegography.ReadShapefile(infile)

	if err != nil {
		panic(err)
	}

	gj, err := fc.ToPrettyGeoJSON()

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(infile[:len(infile)-4] + ".geojson", gj, 0755)

	if err != nil {
		panic(err)
	}
}
```

## Installing
Install with `go get github.com/Froglich/gegography` .
