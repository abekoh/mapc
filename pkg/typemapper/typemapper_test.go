package typemapper

import (
	"testing"

	"golang.org/x/mod/modfile"
)

func TestPath(t *testing.T) {
	f, err := modfile.Parse("go.mod", file_bytes, nil)
	if err != nil {
		panic(err)
	}
	t.Log(f)
}
