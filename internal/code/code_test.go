package code

import (
	_ "embed"
	"testing"
)

//go:embed testdata/mapper/mapper.go
var mapperFile string

func TestNewFile(t *testing.T) {
	_, err := NewFile("testdata/mapper/mapper.go")
	if err != nil {
		t.Fatal(err)
	}
}
