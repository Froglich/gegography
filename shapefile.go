package gegography

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const shpMin float64 = -1e38

func parseShpPoint(in []byte) (Point, error) {
	if len(in) != 16 {
		return Point{}, GeoFormatError{Msg: "point has wrong number of bytes"}
	}

	var x float64
	var y float64

	err := parseValue(in[0:8], binary.LittleEndian, &x)
	if err != nil {
		return Point{}, err
	}

	err = parseValue(in[8:16], binary.LittleEndian, &y)
	if err != nil {
		return Point{}, err
	}

	return Point{X: x, Y: y}, nil
}

func parseShpMultiPoint(in []byte) (MultiPoint, error) {
	if len(in) < 52 {
		return nil, GeoFormatError{Msg: "multipoint with too few bytes"}
	}

	mp := make(MultiPoint, 0)

	var nr int32
	err := parseValue(in[32:36], binary.LittleEndian, &nr)
	if err != nil {
		return nil, err
	}

	if len(in) < 36+int(nr)*16 {
		return nil, GeoFormatError{Msg: "multipoint is malformed"}
	}

	for z := 0; z < int(nr); z++ {
		s := 36 + z*16

		p, err := parseShpPoint(in[s : s+16])
		if err != nil {
			return nil, err
		}

		mp = append(mp, p)
	}

	return mp, err
}

func parseShpPolyLine(in []byte) (Polygon, error) {
	if len(in) < 44 {
		return nil, GeoFormatError{Msg: "polyline with too few bytes"}
	}

	p := make(Polygon, 0)

	var nparts int32
	err := parseValue(in[32:36], binary.LittleEndian, &nparts)
	if err != nil {
		return nil, err
	}

	var npoints int32
	err = parseValue(in[36:40], binary.LittleEndian, &npoints)
	if err != nil {
		return nil, err
	}

	parts := make([]int32, int(nparts))
	points := make([]Point, int(npoints))

	var v int32
	for x := 0; x < int(nparts); x++ {
		s := 40 + x*4
		err = parseValue(in[s:s+4], binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		parts[x] = v
	}

	for x := 0; x < int(npoints); x++ {
		s := 40 + 4*int(nparts) + 16*x
		point, err := parseShpPoint(in[s : s+16])
		if err != nil {
			return nil, err
		}
		points[x] = point
	}

	i := len(parts) - 1
	for x := 0; x < i; x++ {
		p = append(p, points[parts[x]:parts[x+1]])
	}

	p = append(p, points[parts[i]:])

	return p, nil
}

func parseValue(in []byte, order binary.ByteOrder, out interface{}) error {
	buf := bytes.NewReader(in)
	return binary.Read(buf, order, out)
}

type shpGeography struct {
	Features []Feature
	Error    error
}

func shpGeographyReader(filename string, result chan shpGeography) {
	r, err := os.Open(filename)

	if err != nil {
		result <- shpGeography{Error: err}
		return
	}

	header := make([]byte, 100)

	_, err = r.Read(header)
	if err != nil {
		result <- shpGeography{Error: err}
		return
	}

	var fileLength int32
	err = parseValue(header[24:28], binary.BigEndian, &fileLength) //the only part of the header I care about for reading.
	if err != nil {
		result <- shpGeography{Error: err}
		return
	}

	var pos int32 = 100
	var t int32
	rh := make([]byte, 8)

	var cl int32

	features := make([]Feature, 0)

	for pos < fileLength*2 {
		_, err = r.Read(rh)
		if err != nil {
			result <- shpGeography{Error: err}
			return
		}
		pos += 8

		err = parseValue(rh[4:8], binary.BigEndian, &cl)
		if err != nil {
			result <- shpGeography{Error: err}
			return
		}

		content := make([]byte, cl*2)
		_, err = r.Read(content)
		if err != nil {
			result <- shpGeography{Error: err}
			return
		}
		pos += cl * 2

		err = parseValue(content[0:4], binary.LittleEndian, &t)
		if err != nil {
			result <- shpGeography{Error: err}
			return
		}

		var c interface{}

		switch t {
		case 1:
			c, err = parseShpPoint(content[4:])
			features = append(features, Feature{Type: "Point", Coordinates: c})
		case 8:
			c, err = parseShpMultiPoint(content[4:])
			features = append(features, Feature{Type: "MultiPoint", Coordinates: c})
		case 3:
			c, err = parseShpPolyLine(content[4:])
			features = append(features, Feature{Type: "MultiLineString", Coordinates: c})
		case 5:
			c, err = parseShpPolyLine(content[4:])
			features = append(features, Feature{Type: "Polygon", Coordinates: c})
		default:
			result <- shpGeography{Error: GeoTypeError{Type: fmt.Sprintf("unsupported shapefile geographical type '%v'", t)}}
			return
		}

		if err != nil {
			result <- shpGeography{Error: err}
			return
		}
	}

	result <- shpGeography{Features: features}
}

type dBASEColumn struct {
	Name     string
	Index    int
	DataType byte
	Size     int
}

func (dbc *dBASEColumn) castValue(inVal string) (outVal interface{}, err error) {
	switch dbc.DataType {
	case 'N':
		outVal, err = strconv.ParseFloat(inVal, 64)
	case 'F':
		outVal, err = strconv.ParseFloat(inVal, 64)
	case 'O':
		outVal, err = strconv.ParseFloat(inVal, 64)
	default:
		outVal = inVal
	}

	return
}

type dBASETable struct {
	FileExists bool
	Properties []map[string]interface{}
	Error      error
	Columns    []dBASEColumn
}

func (dbr *dBASETable) addColumn(name string, index int, dt byte, size int) {
	dbr.Columns = append(dbr.Columns, dBASEColumn{
		Name:     name,
		Index:    index,
		DataType: dt,
		Size:     size,
	})
}

func (dbr *dBASETable) addRow(row map[string]interface{}) {
	dbr.Properties = append(dbr.Properties, row)
}

func dBASEReader(filename string, result chan dBASETable) {
	ret := dBASETable{
		Columns:    make([]dBASEColumn, 0),
		Properties: make([]map[string]interface{}, 0),
	}

	returnError := func(err error) {
		ret.Error = err
		result <- ret
	}

	r, err := os.Open(filename)
	if err != nil {
		returnError(err)
		return
	}

	d, err := r.Stat()
	if err != nil {
		returnError(err)
		return
	}

	ret.FileExists = true //technically, it may exist before this point as well.

	c := make([]byte, d.Size())

	_, err = r.Read(c)
	if err != nil {
		returnError(err)
		return
	}

	var nrOfRecords uint32
	var headerSize uint16
	var recordLength uint16

	err = parseValue(c[4:8], binary.LittleEndian, &nrOfRecords)
	if err != nil {
		returnError(err)
		return
	}

	err = parseValue(c[8:10], binary.LittleEndian, &headerSize)
	if err != nil {
		returnError(err)
		return
	}

	err = parseValue(c[10:12], binary.LittleEndian, &recordLength)
	if err != nil {
		returnError(err)
		return
	}

	nrOfColumns := int((headerSize - 33) / 32)

	var size uint8
	for x := 0; x < nrOfColumns; x++ {
		offset := x*32 + 32
		fieldName := string(c[offset : offset+10])

		err := parseValue(c[offset+16:offset+17], binary.LittleEndian, &size)
		if err != nil {
			returnError(err)
			return
		}

		ret.addColumn(fieldName, x, c[offset+11], int(size))
	}

	for row := 0; row < int(nrOfRecords); row++ {
		newRow := make(map[string]interface{})
		rowOffset := int(headerSize) + row*int(recordLength)
		recordOffset := 1
		prevRecordSize := 0

		for x := range ret.Columns {
			column := ret.Columns[x]
			recordOffset += prevRecordSize
			prevRecordSize += column.Size

			temp := c[rowOffset+recordOffset : rowOffset+recordOffset+column.Size]
			record := string(temp)
			for i := 0; i < len(temp); i++ {
				if temp[i] == 0x00 {
					record = string(temp[:i])
					break
				}
			}

			val, err := column.castValue(strings.TrimSpace(record))
			if err != nil {
				returnError(err)
				return
			}

			newRow[column.Name] = val
		}

		ret.addRow(newRow)
	}

	result <- ret
	return
}

//ParseShapefile reads a shapefile (and accompanying dBASE-table, if any) into a FeatureCollection
func ReadShapefile(shapeFile string) (FeatureCollection, error) {
	if !strings.HasSuffix(shapeFile, ".shp") {
		return FeatureCollection{}, GeoFormatError{Msg: fmt.Sprintf("%v does not appear to be a shapefile", shapeFile)}
	}

	tabFile := shapeFile[:len(shapeFile)-4] + ".dbf"
	tabChan := make(chan dBASETable)
	geoChan := make(chan shpGeography)

	go shpGeographyReader(shapeFile, geoChan)
	go dBASEReader(tabFile, tabChan)
	tab := <-tabChan
	geo := <-geoChan

	fc := FeatureCollection{}

	if geo.Error != nil {
		return FeatureCollection{}, geo.Error
	}

	fc.Features = geo.Features

	if tab.FileExists {
		if tab.Error != nil {
			return FeatureCollection{}, tab.Error
		}

		if len(tab.Properties) != len(geo.Features) {
			return FeatureCollection{}, GeoFormatError{Msg: "mismatching number of rows in attribute table and shapefile"}
		}

		for x := range fc.Features {
			fc.Features[x].Properties = tab.Properties[x]
		}
	}

	return fc, nil
}
