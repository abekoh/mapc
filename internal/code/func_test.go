package code

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/abekoh/mapc/internal/mapping"
)

type From struct {
	Int int
}

type To struct {
	Int int
}

func TestNewFunc(t *testing.T) {
	mapper := mapping.NewMapper()
	mp, err := mapper.Map(From{}, To{})
	if err != nil {
		t.Fatal(err)
	}
	f := NewFile("github.com/abekoh/mapc/main")
	fc := NewFunc(mp)
	f.Apply(fc)
	var buf bytes.Buffer
	err = f.Fprint(&buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", buf.String())
}
