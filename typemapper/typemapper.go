package typemapper

import (
	"github.com/abekoh/mapc/types"
)

type TypeMapper interface {
	Map(from, to *types.Typ) (Caster, bool)
}

var Default = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
}
