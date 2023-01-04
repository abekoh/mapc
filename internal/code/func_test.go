package code

import (
	"testing"

	"github.com/abekoh/mapc/internal/mapping"
)

type AUser struct {
	Int int
}

type BUser struct {
	Int int
}

func TestNewFunc(t *testing.T) {
	mapper := mapping.NewMapper()
	mp, err := mapper.Map(AUser{}, BUser{})
	if err != nil {
		t.Fatal(err)
	}
	f := NewFile("github.com/abekoh/mapc/main")
	fc := NewFunc(mp)
	f.Apply(fc)
	got, err := f.sPrint()
	if err != nil {
		t.Fatal(err)
	}
	expected := `package main

func ToBUser(from AUser) BUser {
	return BUser{
		Int: from.Int,
	}
}
`
	if got != expected {
		t.Errorf("not matched with expected.\n\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}
