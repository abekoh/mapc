package mapping

import (
	"errors"
	"fmt"
	"github.com/abekoh/mapc/internal/tokenizer"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type StructParam struct {
	Dir    string
	Pkg    string
	Struct string
}

type Struct interface {
	PkgName() string
	PkgPath() string
	StructName() string
	Var(fieldName string) *types.Var
	String() string
}

type PkgStruct struct {
	pkg        *packages.Package
	str        *types.Struct
	structName string
}

func (s PkgStruct) PkgName() string {
	return s.pkg.Name
}

func (s PkgStruct) PkgPath() string {
	return s.pkg.PkgPath
}

func (s PkgStruct) StructName() string {
	return s.structName
}

func (s PkgStruct) Var(fieldName string) *types.Var {
	for i := 0; i < s.str.NumFields(); i++ {
		if s.str.Field(i).Name() == fieldName {
			return s.str.Field(i)
		}
	}
	return nil
}

func (s PkgStruct) String() string {
	return fmt.Sprintf("%+v", s.str.String())
}

func findPkg(param StructParam) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Dir:  param.Dir,
	})
	if err != nil {
		return nil, fmt.Errorf("error is occured when loading package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, errors.New("package not found")
	}
	for _, pkg := range pkgs {
		if pkg.Name == param.Pkg {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("package %s is not found in %s", param.Pkg, param.Dir)
}

func newStruct(param StructParam) (*PkgStruct, error) {
	pkg, err := findPkg(param)
	if err != nil {
		return nil, fmt.Errorf("failed to find package: %w", err)
	}
	obj := pkg.Types.Scope().Lookup(param.Struct)
	if obj == nil {
		return nil, fmt.Errorf("model %s is not found", param.Struct)
	}
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not *types.struct", param.Struct)
	}
	return &PkgStruct{pkg: pkg, str: str, structName: param.Struct}, nil
}

func (s PkgStruct) tokenFieldMap() tokenFieldMap {
	res := make(tokenFieldMap)
	for i := 0; i < s.str.NumFields(); i++ {
		f := s.str.Field(i)
		token := tokenizer.Tokenize(f.Name())
		res[token] = Var{v: f}
	}
	return res
}
