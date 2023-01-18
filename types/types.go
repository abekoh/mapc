package types

import (
	"fmt"
	"reflect"
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
	return f.pkgPath
}

func (f Field) IsSameTypeAndPkgWith(x *Field) bool {
	return f.Kind() == x.Kind() && f.Typ().PkgPath() == x.Typ().PkgPath()
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
	if t.IsPointer() {
		e, _ := t.Elem()
		return fmt.Sprintf("*%s", e.String())
	}
	return t.refTyp.Name()
}

func (t Typ) IsPointer() bool {
	return t.refTyp.Kind() == reflect.Pointer
}

func (t Typ) Elem() (*Typ, bool) {
	if !t.IsPointer() {
		return nil, false
	}
	e := t.refTyp.Elem()
	if e == nil {
		return nil, false
	}
	return NewTyp(e), true
}
