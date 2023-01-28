package str

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/mapcstd"
)

type Identifier interface {
	Name() string
	PkgPath() string
}

type Struct struct {
	name    string
	pkgPath string
	Fields  []*Field
}

func NewStruct(typ reflect.Type) (*Struct, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("kind must be struct, got %v", typ.Kind())
	}
	var fields []*Field
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		field := &Field{refStructField: &structField, pkgPath: typ.PkgPath()}
		fields = append(fields, field)
	}
	return &Struct{
		name:    typ.Name(),
		pkgPath: typ.PkgPath(),
		Fields:  fields,
	}, nil
}

func (s Struct) Name() string {
	return s.name
}

func (s Struct) PkgPath() string {
	return s.pkgPath
}

type Field struct {
	pkgPath        string
	refStructField *reflect.StructField
}

func (f Field) Name() string {
	return f.refStructField.Name
}

func (f Field) Typ() *mapcstd.Typ {
	return mapcstd.NewTyp(f.refStructField.Type)
}

func (f Field) TypeName() string {
	return f.refStructField.Type.Name()
}

func (f Field) Kind() reflect.Kind {
	return f.refStructField.Type.Kind()
}

func (f Field) PkgPath() string {
	return f.pkgPath
}

func (f Field) IsSameTypeAndPkgWith(x *Field) bool {
	return f.Kind() == x.Kind() && f.Typ().PkgPath() == x.Typ().PkgPath()
}
