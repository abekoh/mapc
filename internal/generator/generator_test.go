package generator

import (
	"bytes"
	_ "embed"
	"github.com/abekoh/mapstructor/internal/pivot"
	"testing"
)

//go:embed testdata/ab/ab.golden.go
var golden string

func TestGenerator_Generate(t *testing.T) {
	var w bytes.Buffer
	from := pivot.StructParam{
		Dir:    "testdata/a",
		Pkg:    "a",
		Struct: "User",
	}
	to := pivot.StructParam{
		Dir:    "testdata/b",
		Pkg:    "b",
		Struct: "User",
	}
	g := Generator{}
	err := g.Generate(&w, from, to)
	if err != nil {
		t.Fatal(err)
	}
	actual := w.String()
	if golden != actual {
		t.Errorf(`not equal generated file:

expected:
%s

actual:
%s
`, golden, actual)
	}
}
