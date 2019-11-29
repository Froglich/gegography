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
