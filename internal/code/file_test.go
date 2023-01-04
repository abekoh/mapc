package code

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/mapper/mapper.go
var mapperRawFile string

func loadSample(t *testing.T) *File {
	t.Helper()
	f, err := LoadFile("testdata/mapper/mapper.go", "github.com/abekoh/mapc/internal/code/testdata/mapper")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestLoadFile(t *testing.T) {
	loadSample(t)
}

func TestNew(t *testing.T) {
	f := NewFile("github.com/abekoh/mapc/main")

	// check writing
	var buf bytes.Buffer
	if err := f.Fprint(&buf); err != nil {
		t.Fatal(err)
	}
	expected := `package main
`
	r := buf.String()
	if r != expected {
		t.Errorf("output is must be empty")
	}
}

func TestFile_Write(t *testing.T) {
	f := loadSample(t)
	var buf bytes.Buffer
	if err := f.Fprint(&buf); err != nil {
		t.Fatal(err)
	}
	r := buf.String()
	if r != mapperRawFile {
		t.Errorf("output is not matched. expected: %s, got: %s", mapperRawFile, r)
	}
}

func TestFile_FindFunc(t *testing.T) {
	f := loadSample(t)
	funcName := "ToBUser"
	_, ok := f.FindFunc(funcName)
	if !ok {
		t.Errorf("not found %s", funcName)
	}
}
