package mapping

import (
	"fmt"
	"reflect"
)

// old

//type StructParam struct {
//	Dir    string
//	Pkg    string
//	Struct string
//}
//
//type PkgStruct struct {
//	pkg        *packages.Package
//	str        *types.Struct
//	structName string
//}
//
//function (s PkgStruct) PkgName() string {
//	return s.pkg.Name
//}
//
//function (s PkgStruct) PkgPath() string {
//	return s.pkg.PkgPath
//}
//
//function (s PkgStruct) StructName() string {
//	return s.structName
//}
//
//function (s PkgStruct) Var(fieldName string) *types.Var {
//	for i := 0; i < s.str.NumFields(); i++ {
//		if s.str.Field(i).Name() == fieldName {
//			return s.str.Field(i)
//		}
//	}
//	return nil
//}
//
//function (s PkgStruct) String() string {
//	return fmt.Sprintf("%+v", s.str.String())
//}
//
//function findPkg(param StructParam) (*packages.Package, error) {
//	pkgs, err := packages.Load(&packages.Config{
//		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
//		Dir:  param.Dir,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("error is occured when loading package: %w", err)
//	}
//	if len(pkgs) == 0 {
//		return nil, errors.New("package not found")
//	}
//	for _, pkg := range pkgs {
//		if pkg.Name == param.Pkg {
//			return pkg, nil
//		}
//	}
//	return nil, fmt.Errorf("package %s is not found in %s", param.Pkg, param.Dir)
//}
//
//function (s PkgStruct) fieldMap() tokenFieldMapOld {
//	res := make(tokenFieldMapOld)
//	for i := 0; i < s.str.NumFields(); i++ {
//		f := s.str.Field(i)
//		token := tokenizer.Tokenize(f.Name())
//		res[token] = Var{v: f}
//	}
//	return res
//}
//
//function newStruct(param StructParam) (*PkgStruct, error) {
//	pkg, err := findPkg(param)
//	if err != nil {
//		return nil, fmt.Errorf("failed to find package: %w", err)
//	}
//	obj := pkg.Types.Scope().Lookup(param.Struct)
//	if obj == nil {
//		return nil, fmt.Errorf("model %s is not found", param.Struct)
//	}
//	str, ok := obj.Type().Underlying().(*types.Struct)
//	if !ok {
//		return nil, fmt.Errorf("%s is not *types.struct", param.Struct)
//	}
//	return &PkgStruct{pkg: pkg, str: str, structName: param.Struct}, nil
//}

// new

type Struct struct {
	Name    string
	PkgPath string
	Fields  []Field
}

type Field struct {
	f *reflect.StructField
}

func (f Field) Name() string {
	return f.f.Name
}

func NewStruct(t reflect.Type) (*Struct, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("kind must be struct, got %v", t.Kind())
	}
	var fs []Field
	for i := 0; i < t.NumMethod(); i++ {
		f := t.Field(i)
		s := Field{f: &f}
		fs = append(fs, s)
	}
	return &Struct{
		Name:    t.Name(),
		PkgPath: t.PkgPath(),
		Fields:  fs,
	}, nil
}
