package pivot

import (
	"github.com/abekoh/mapc/internal/tokenizer"
	"go/types"
	"golang.org/x/tools/go/packages"
	"reflect"
	"testing"
)

type Sample struct {
	CapCamel   struct{}
	smallCamel struct{}
	snake_case struct{}
}

func loadStruct(t *testing.T, target any) *Struct {
	t.Helper()
	structName := reflect.TypeOf(target).Name()
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Tests: true,
		Dir:   "",
	})
	if err != nil {
		t.Fatal(err)
	}
	var pkg *packages.Package
	for _, p := range pkgs {
		if p.ID == "github.com/abekoh/mapc/internal/pivot [github.com/abekoh/mapc/internal/pivot.test]" {
			pkg = p
			break
		}
	}
	if pkg == nil {
		t.Fatal("package is not found")
	}
	obj := pkg.Types.Scope().Lookup(structName)
	if obj == nil {
		t.Fatal("object is not found")
	}
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		t.Fatal("failed to cast")
	}
	return &Struct{
		pkg:        pkg,
		str:        str,
		structName: structName,
	}
}

func TestStruct_tokenFieldMap(t *testing.T) {
	t.Run("no tokenizer", func(t *testing.T) {
		s := loadStruct(t, Sample{})
		tokenizer.Initialize() // no tokenizers
		got := s.tokenFieldMap()
		expectedFields := []string{
			"CapCamel",
			"smallCamel",
			"snake_case",
		}
		for _, expected := range expectedFields {
			if _, ok := got[expected]; !ok {
				t.Errorf("field '%s' is not found in got", expected)
			}
		}
	})
	t.Run("upperFirst", func(t *testing.T) {
		s := loadStruct(t, Sample{})
		tokenizer.Initialize(tokenizer.UpperFirst)
		got := s.tokenFieldMap()
		expectedFields := []string{
			"CapCamel",
			"SmallCamel",
			"Snake_case",
		}
		for _, expected := range expectedFields {
			if _, ok := got[expected]; !ok {
				t.Errorf("field '%s' is not found in got", expected)
			}
		}
	})
	t.Run("snakeToCamel", func(t *testing.T) {
		s := loadStruct(t, Sample{})
		tokenizer.Initialize(tokenizer.SnakeToCamel)
		got := s.tokenFieldMap()
		expectedFields := []string{
			"CapCamel",
			"smallCamel",
			"snakeCase",
		}
		for _, expected := range expectedFields {
			if _, ok := got[expected]; !ok {
				t.Errorf("field '%s' is not found in got", expected)
			}
		}
	})
	t.Run("snakeToCamel,upperFirst", func(t *testing.T) {
		s := loadStruct(t, Sample{})
		tokenizer.Initialize(tokenizer.SnakeToCamel, tokenizer.UpperFirst)
		got := s.tokenFieldMap()
		expectedFields := []string{
			"CapCamel",
			"SmallCamel",
			"SnakeCase",
		}
		for _, expected := range expectedFields {
			if _, ok := got[expected]; !ok {
				t.Errorf("field '%s' is not found in got", expected)
			}
		}
	})
}
