package gegography

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func readToBuffer(zf *zip.File) (*bytes.Reader, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func TestReadShapefile(t *testing.T) {
	fc, err := ReadShapefile("test_data/test_shapefile.shp")
	if err != nil {
		t.Error(err)
	}

	if len(fc.Features) != 1 {
		t.Errorf("ReadShapefile('test_data/test_shapefile.shp'), want 1 feature got %d", len(fc.Features))
	}

	if v, ok := fc.Features[0].Properties["TestField"]; !ok || v != "Hello!" {
		t.Error("ReadShapefile('test_data/test_shapefile.shp'), only feature should have property 'TestField' with value 'Hello!'")
	}
}

func TestReadZippedShapefile(t *testing.T) {
	r, err := zip.OpenReader("test_data/test_shapefile.zip")
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	var shpReader *bytes.Reader
	var dbReader *bytes.Reader

	for _, f := range r.File {
		switch f.Name {
		case "test_shapefile.shp":
			shpReader, err = readToBuffer(f)
		case "test_shapefile.dbf":
			dbReader, err = readToBuffer(f)
		}

		if err != nil {
			t.Errorf("Unable to read zipped test data")
		}
	}

	if shpReader == nil || dbReader == nil {
		t.Error("Zipped shapefile should contain test_shapefile.shp and test_shapefile.dbf")
	}

	fc, err := ReadShapefileData(shpReader, dbReader)
	if err != nil {
		t.Error(err)
	}

	if len(fc.Features) != 1 {
		t.Errorf("ReadShapefileData(shapefileReader, databaseReader), want 1 feature got %d", len(fc.Features))
	}

	if v, ok := fc.Features[0].Properties["TestField"]; !ok || v != "Hello!" {
		t.Error("ReadShapefileData(shapefileReader, databaseReader), only feature should have property 'TestField' with value 'Hello!'")
	}
}
