package pivot

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

type Struct struct {
	pkg        *packages.Package
	str        *types.Struct
	structName string
}

func (s Struct) PkgName() string {
	return s.pkg.Name
}

func (s Struct) PkgPath() string {
	return s.pkg.PkgPath
}

func (s Struct) StructName() string {
	return s.structName
}

func (s Struct) String() string {
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

func newStruct(param StructParam) (*Struct, error) {
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
	return &Struct{pkg: pkg, str: str, structName: param.Struct}, nil
}

func (s Struct) tokenFieldMap() tokenFieldMap {
	res := make(tokenFieldMap)
	for i := 0; i < s.str.NumFields(); i++ {
		f := s.str.Field(i)
		token := tokenizer.Tokenize(f.Name())
		res[token] = Var{v: f}
	}
	return res
}
