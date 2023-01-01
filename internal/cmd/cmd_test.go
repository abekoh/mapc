package cmd

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/ab/ab.golden.go
var golden string

func TestGenerator_Generate(t *testing.T) {
	var w bytes.Buffer
	from := Param{
		Dir:    "testdata/a",
		Pkg:    "a",
		Struct: "User",
	}
	to := Param{
		Dir:    "testdata/b",
		Pkg:    "b",
		Struct: "User",
	}
	distPkgName := "ab"
	err := Generate(&w, from, to, distPkgName)
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
