package code

import (
	_ "embed"
	"testing"
)

//go:embed testdata/mapper/mapper.go
var mapperFile string

func TestLoadFile(t *testing.T) {
	_, err := LoadFile("testdata/mapper/mapper.go")
	if err != nil {
		t.Fatal(err)
	}
}
