package code

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/mapper/mapper.go
var mapperRawFile string

func loadMapper(t *testing.T) *File {
	t.Helper()
	f, err := LoadFile("testdata/mapper/mapper.go", "github.com/abekoh/mapc/internal/code/testdata/mapper")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestLoadFile(t *testing.T) {
	loadMapper(t)
}

func TestFile_Write(t *testing.T) {
	f := loadMapper(t)
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	if buf.String() != mapperRawFile {
		t.Errorf("failed")
	}
}
