package mapcstd

import (
	"fmt"
	"reflect"
)

type Fun struct {
	srcTyp  *Typ
	destTyp *Typ
	name    string
	pkgPath string
	retType ReturnType
}

type ReturnType int

const (
	OnlyValue ReturnType = iota
	WithError
	WithOk
)

func NewFunOf(a any) (*Fun, error) {
	// TODO: add test, refactoring
	typ := reflect.TypeOf(a)
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("kind must be func, got %v", typ.Kind())
	}
	if typ.NumIn() != 1 {
		return nil, fmt.Errorf("# of inputs must be one")
	}
	var destTyp *Typ
	var retType ReturnType
	switch typ.NumOut() {
	case 1:
		destTyp = NewTyp(typ.Out(0))
		retType = OnlyValue
	case 2:
		destTyp = NewTyp(typ.Out(0))
		e := new(error)
		if typ.Out(1).Kind() == reflect.Bool {
			retType = WithOk
		} else if typ.Out(1).Implements(reflect.TypeOf(e)) {
			retType = WithError
		} else {
			return nil, fmt.Errorf("second return value must be error or bool")
		}
	default:
		return nil, fmt.Errorf("# of outputs must be one or two")
	}
	return &Fun{
		srcTyp:  NewTyp(typ.In(0)),
		destTyp: destTyp,
		name:    typ.Name(),
		pkgPath: typ.PkgPath(),
		retType: retType,
	}, nil
}

func (f Fun) SameInOut(src, dest *Typ) bool {
	return src.Equals(f.srcTyp) && dest.Equals(f.destTyp)
}

func (f Fun) Name() string {
	return f.name
}

func (f Fun) PkgPath() string {
	return f.pkgPath
}
