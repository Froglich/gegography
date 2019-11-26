# Gegography
Gegography is a library for reading, manipulating and converting
geographical formats, written i pure Go. Currently, support is limited to
GeoJSON and WKT (Well-Known-Text). The basis for Shapefile support (read only)
is implemented, and I plan on implementing WKB (Well-Known-Binary) as well.

I do not plan on supporting writing Shapefiles. [Shapefile must die!](http://switchfromshapefile.org/)

Gegography supports the following geographical types:
* Point
* LineString
* Polygon
* MultiPoint
* MultiLineString
* MultiPoint

Currently, only XY geometries are supported.

## Installing
Install with `go get github.com/Froglich/gegography` .

## ToDo
* Add capability to read dBASE files (i.e. shapefile attributes)
* Finish shapefile implementation (export FeatureCollection)
* Implement WKB (Well-Known-Binary)
