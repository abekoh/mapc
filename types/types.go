package types

import (
	"fmt"
	"reflect"
)

type Info interface {
	Name() string
	PkgPath() string
}

type Struct struct {
	name    string
	PkgPath string
	Fields  []*Field
}

func NewStruct(typ reflect.Type) (*Struct, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("kind must be struct, got %v", typ.Kind())
	}
	var fields []*Field
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		field := &Field{refStructField: &structField}
		fields = append(fields, field)
	}
	return &Struct{
		name:    typ.Name(),
		PkgPath: typ.PkgPath(),
		Fields:  fields,
	}, nil
}

func (s Struct) Name() string {
	return s.name
}

//
//func (s Struct) PkgPath() string {
//	return s.PkgPath
//}

type Field struct {
	refStructField *reflect.StructField
}

func (f Field) Name() string {
	return f.refStructField.Name
}

func (f Field) Typ() *Typ {
	return NewTyp(f.refStructField.Type)
}

func (f Field) TypeName() string {
	return f.refStructField.Type.Name()
}

func (f Field) Kind() reflect.Kind {
	return f.refStructField.Type.Kind()
}

func (f Field) PkgPath() string {
	return f.refStructField.Type.PkgPath()
}

func (f Field) IsSameTypeAndPkgWith(x *Field) bool {
	return f.Kind() == x.Kind() && f.PkgPath() == x.PkgPath()
}

type Typ struct {
	refTyp reflect.Type
}

func NewTyp(refTyp reflect.Type) *Typ {
	return &Typ{refTyp: refTyp}
}

func (t Typ) PkgPath() string {
	return t.refTyp.PkgPath()
}

func (t Typ) Name() string {
	return t.refTyp.Name()
}

func (t Typ) AssignableTo(to *Typ) bool {
	return t.refTyp.AssignableTo(to.refTyp)
}

func (t Typ) ConvertibleTo(to *Typ) bool {
	return t.refTyp.ConvertibleTo(to.refTyp)
}

func (t Typ) String() string {
	return t.refTyp.Name()
}
