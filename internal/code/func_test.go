package code

import (
	"testing"

	"github.com/abekoh/mapc/internal/mapping"
)

type AUser struct {
	Int int
	Str string
}

type BUser struct {
	Int int
	Str string
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
		Str: from.Str,
	}
}
`
	if got != expected {
		t.Errorf("not matched with expected.\n\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}

func TestFunc_ReNew(t *testing.T) {
	mapper := mapping.NewMapper()
	mp, err := mapper.Map(AUser{}, BUser{})
	if err != nil {
		t.Fatal(err)
	}
	f, err := loadFileFromString(`package main

func ToBUser(from AUser) BUser {
	return BUser{
		// Int: from.Int,
		Str: from.Str,
	}
}`, "main")
	if err != nil {
		t.Fatal(err)
	}
	existedFc, ok := f.FindFunc("ToBUser")
	if !ok {
		t.Fatal("target func is not found")
	}
	gotFc, err := existedFc.ReNew(mp)
	if err != nil {
		t.Fatal(err)
	}
	f.Apply(gotFc)
	got, err := f.sPrint()
	if err != nil {
		t.Fatal(err)
	}
	expected := `package main

func ToBUser(from AUser) BUser {
	return BUser{
		// Int: from.Int,
		Str: from.Str,
	}
}
`
	if got != expected {
		t.Errorf("not matched with expected.\n\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}
