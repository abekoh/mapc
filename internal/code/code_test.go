package code

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/abekoh/mapc/internal/mapping"
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
	if err := f.Write(&buf); err != nil {
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
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	r := buf.String()
	if r != mapperRawFile {
		t.Errorf("output is not matched. expected: %s, got: %s", mapperRawFile, r)
	}
}

func TestFile_FindFunc(t *testing.T) {
	f := loadSample(t)
	funcName := "ToAUser"
	_, ok := f.FindFunc(funcName)
	if !ok {
		t.Errorf("not found %s", funcName)
	}
}

type From struct {
	Int   int
	Int64 int64
}

type To struct {
	Int   int
	Int64 int64
}

func TestNewFunc(t *testing.T) {
	mapper := mapping.NewMapper()
	mp, err := mapper.Map(From{}, To{})
	if err != nil {
		t.Fatal(err)
	}
	f := NewFile("main")
	fc := NewFunc(mp)
	f.Apply(fc)
	var buf bytes.Buffer
	err = f.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", buf.String())
}
