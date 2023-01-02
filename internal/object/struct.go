package object

import (
	"fmt"
	"reflect"
)

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
	for i := 0; i < t.NumField(); i++ {
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
