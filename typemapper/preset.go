package typemapper

import (
	"github.com/abekoh/mapc/types"
)

type AssignMapper struct {
}

func (a AssignMapper) Map(from, to *types.Typ) (Caster, bool) {
	if from.AssignableTo(to) {
		return &NopCaster{}, true
	}
	return nil, false
}

type ConvertMapper struct {
}

func (c ConvertMapper) Map(from, to *types.Typ) (Caster, bool) {
	if from.ConvertibleTo(to) {
		return &SimpleCaster{
			pkgPath: to.PkgPath(),
			fn:      to.Name(),
		}, true
	}
	return nil, false
}

type RefMapper struct {
}

func (p RefMapper) Map(from, to *types.Typ) (Caster, bool) {
	if from.IsPointer() || !to.IsPointer() {
		return nil, false
	}
	if toElm, ok := to.Elem(); ok && from.AssignableTo(toElm) {
		return &SimpleCaster{
			fn: "&",
		}, true
	}
	return nil, false
}

type DerefMapper struct {
}

func (p DerefMapper) Map(from, to *types.Typ) (Caster, bool) {
	if !from.IsPointer() || to.IsPointer() {
		return nil, false
	}
	if fromElm, ok := from.Elem(); ok && fromElm.AssignableTo(to) {
		return &SimpleCaster{
			fn: "*",
		}, true
	}
	return nil, false
}

type MapTypeMapper map[types.Typ]map[types.Typ]Caster

func (m MapTypeMapper) Map(from, to types.Typ) (Caster, bool) {
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
