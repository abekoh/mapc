package code

import (
	_ "embed"
	"testing"
)

//go:embed testdata/mapper/mapper.go
var mapperFile string

func TestNewFile(t *testing.T) {
	got, err := NewFile("testdata/mapper/mapper.go")
	if err != nil {
		t.Fatal(err)
	}
	if got.code != mapperFile {
		t.Errorf("got.code is not match with expected")
	}
	importLen := 2
	if len(got.imports) != importLen {
		t.Errorf("len(got.imports) must be %d", importLen)
	}
	expectedImports := []string{
		"time",
		"github.com/google/uuid",
	}
	for _, expected := range expectedImports {
		if _, ok := got.imports[expected]; !ok {
			t.Errorf("%s is not found in import()", expected)
		}
	}
}
