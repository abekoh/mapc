package model

import (
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type StructParam struct {
	Dir    string
	Pkg    string
	Struct string
}

type Struct struct {
	str   *types.Struct
	param *StructParam
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

func Lookup(param StructParam) (*Struct, error) {
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
		return nil, fmt.Errorf("%s is not model", param.Struct)
	}
	return &Struct{str: str, param: &param}, nil
}

type FieldMap struct {
	from *types.Var
	to   *types.Var
}
