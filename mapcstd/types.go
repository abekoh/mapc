package mapcstd

import (
	"fmt"
	"reflect"
)

type Typ struct {
	refTyp reflect.Type
}

func NewTypOf(t any) *Typ {
	return &Typ{refTyp: reflect.TypeOf(t)}
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

func (t Typ) Equals(x *Typ) bool {
	return t.Name() == x.Name() && t.PkgPath() == x.PkgPath() && t.refTyp.Kind() == x.refTyp.Kind()
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
