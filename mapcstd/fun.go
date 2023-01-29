package mapcstd

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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
		if typ.Out(1).Kind() == reflect.Bool {
			retType = WithOk
		} else if typ.Out(1).Name() == "error" {
			retType = WithError
		} else {
			return nil, fmt.Errorf("second return value must be error or bool")
		}
	default:
		return nil, fmt.Errorf("# of outputs must be one or two")
	}
	pkgPath, name := funcNameAndPkgPath(a)
	return &Fun{
		srcTyp:  NewTyp(typ.In(0)),
		destTyp: destTyp,
		name:    name,
		pkgPath: pkgPath,
		retType: retType,
	}, nil
}

func funcNameAndPkgPath(a any) (pkgPath, name string) {
	// OPTIMIZE: don't use split/join
	val := reflect.ValueOf(a)
	n := runtime.FuncForPC(val.Pointer()).Name()
	sp := strings.Split(n, ".")
	if len(sp) == 1 {
		return "", n
	}
	if len(sp) > 1 {
		return strings.Join(sp[0:len(sp)-1], "."), sp[len(sp)-1]
	}
	return
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
