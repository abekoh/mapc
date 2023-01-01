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

func loadSampleStruct(t *testing.T) *Struct {
	t.Helper()
	structName := reflect.TypeOf(Sample{}).Name()

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
		if p.Name == "pivot" {
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
		s := loadSampleStruct(t)
		tokenizer.Initialize() // no tokenizers
		got := s.tokenFieldMap()

		expectedFields := []string{
			"CapCamel",
		}

		for _, expected := range expectedFields {
			if _, ok := got[expected]; !ok {
				t.Errorf("field '%s' is not found in got", expected)
			}
		}
	})
}
