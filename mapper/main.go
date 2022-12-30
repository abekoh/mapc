package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"io"
	"log"
)

type Mapper struct {
}

type StructParam struct {
	Dir    string
	Pkg    string
	Struct string
}

type Struct struct {
	str *types.Struct
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
		return nil, fmt.Errorf("struct %s is not found", param.Struct)
	}
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not struct", param.Struct)
	}
	return &Struct{str: str}, nil
}

func (m Mapper) Map(w io.Writer, from, to StructParam) error {
	fromStr, err := Lookup(from)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", fromStr, err)
	}
	toStr, err := Lookup(to)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", toStr, err)
	}
	fmt.Printf("%+v, %+v", fromStr.str.String(), toStr.str.String())
	return nil
}

func main() {
	m := Mapper{}
	buf := bytes.Buffer{}
	err := m.Map(&buf, StructParam{
		Dir:    "mapper/a",
		Pkg:    "a",
		Struct: "User",
	}, StructParam{
		Dir:    "mapper/b",
		Pkg:    "b",
		Struct: "User",
	})
	if err != nil {
		log.Fatal(err)
	}
}
