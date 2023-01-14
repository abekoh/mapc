package typemapper

import (
	"github.com/abekoh/mapc/types"
)

type TypeMapper interface {
	Map(from, to *types.Typ) (Caster, bool)
}

type TypeProcessor interface {
	Proc(typ *types.Typ) (Caster, bool)
}

var DefaultMappers = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
	&RefMapper{},
	&DerefMapper{},
}
