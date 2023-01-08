package object

import (
	"fmt"
	"reflect"
)

type Struct struct {
	Name    string
	PkgPath string
	Fields  []*Field
}

type Field struct {
	f *reflect.StructField
}

func (f Field) Name() string {
	return f.f.Name
}

func (f Field) Typ() *Typ {
	return &Typ{t: f.f.Type}
}

func (f Field) TypeName() string {
	return f.f.Type.Name()
}

func (f Field) Kind() reflect.Kind {
	return f.f.Type.Kind()
}

func (f Field) PkgPath() string {
	return f.f.Type.PkgPath()
}

func (f Field) IsSameTypeAndPkgWith(x *Field) bool {
	return f.Kind() == x.Kind() && f.PkgPath() == x.PkgPath()
}

func NewStruct(t reflect.Type) (*Struct, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("kind must be struct, got %v", t.Kind())
	}
	var fs []*Field
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		s := &Field{f: &f}
		fs = append(fs, s)
	}
	return &Struct{
		Name:    t.Name(),
		PkgPath: t.PkgPath(),
		Fields:  fs,
	}, nil
}

type Typ struct {
	t reflect.Type
}

func NewTyp(t reflect.Type) *Typ {
	return &Typ{t: t}
}

func (t Typ) PkgPath() string {
	return t.t.PkgPath()
}

func (t Typ) Name() string {
	return t.t.Name()
}

func (t Typ) AssignableTo(to Typ) bool {
	return t.t.AssignableTo(to.t)
}

func (t Typ) ConvertibleTo(to Typ) bool {
	return t.t.ConvertibleTo(to.t)
}
