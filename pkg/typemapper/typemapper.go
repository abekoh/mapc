package typemapper

import (
	"github.com/abekoh/mapc/internal/object"
)

type TypeMapper interface {
	Map(from, to *object.Typ) (Caster, bool)
}

type AssignMapper struct {
}

func (a AssignMapper) Map(from, to *object.Typ) (Caster, bool) {
	if from.AssignableTo(to) {
		return &NopCaster{}, true
	}
	return nil, false
}

type ConvertMapper struct {
}

func (c ConvertMapper) Map(from, to *object.Typ) (Caster, bool) {
	if from.ConvertibleTo(to) {
		return &SimpleCaster{
			pkgPath: to.PkgPath(),
			fun:     to.Name(),
		}, true
	}
	return nil, false
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

var Default = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
}

type Caster interface {
	PkgPath() string
	Fun() string
}

type NopCaster struct{}

func (n NopCaster) PkgPath() string {
	return ""
}

func (n NopCaster) Fun() string {
	return ""
}

type SimpleCaster struct {
	pkgPath string
	fun     string
}

func (s SimpleCaster) PkgPath() string {
	return s.pkgPath
}

func (s SimpleCaster) Fun() string {
	return s.fun
}
