package model

import (
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type StructParam struct {
	Path    string
	Package string
	Struct  string
}

type Struct struct {
	str *types.Struct
}

func findPkg(param StructParam) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Dir:  param.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("error is occured when loading package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, errors.New("package not found")
	}
	for _, pkg := range pkgs {
		if pkg.Name == param.Package {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("package %s is not found in %s", param.Package, param.Path)
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
	return &Struct{str: str}, nil
}
