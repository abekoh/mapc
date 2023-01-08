package typemapper

import (
	"github.com/abekoh/mapc/internal/object"
)

type TypeMappers []TypeMapper

type TypeMapper interface {
	Map(from, to object.Typ) (Caster, bool)
}

type MapTypeMapper map[object.Typ]map[object.Typ]Caster

func (m MapTypeMapper) Map(from, to object.Typ) (Caster, bool) {
	m2, ok := m[from]
	if !ok {
		return nil, false
	}
	c, ok := m2[to]
	if !ok {
		return nil, false
	}
	return c, true
}

var DefaultTypeMappers = TypeMappers{
	//MapTypeMapper{
	//	//object.NewTyp(reflect.TypeOf("")): map[object.Typ]&BuiltinCaster{Typ: object.NewTyp()},
	//},
	//Int: map[Typ]Caster{
	//	Int16: &BuiltinCaster{Typ: Int16},
	//	Int32: &BuiltinCaster{Typ: Int32},
	//	Int64: &BuiltinCaster{Typ: Int64},
	//},
}

type Caster interface {
	PkgPath() string
	Fun() string
}

type BuiltinCaster struct {
	Typ BuiltinTyp
}

func (b BuiltinCaster) PkgPath() string {
	return ""
}

func (b BuiltinCaster) Fun() string {
	return b.Typ.Name()
}

type BuiltinTyp string

const (
	Int   BuiltinTyp = "int"
	Int8  BuiltinTyp = "int8"
	Int16 BuiltinTyp = "int16"
	Int32 BuiltinTyp = "int32"
	Int64 BuiltinTyp = "int64"
)

func (bt BuiltinTyp) PkgPath() string {
	return ""
}

func (bt BuiltinTyp) Name() string {
	return string(bt)
}
