package typemapper

import (
	"github.com/abekoh/mapc/types"
)

type TypeMapper interface {
	Map(src, dest *types.Typ) (Caster, bool)
}

var Defaults = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
	&RefMapper{},
	&DerefMapper{},
}
