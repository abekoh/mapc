package typemapper

import (
	"github.com/abekoh/mapc/internal/object"
)

type TypeMapper interface {
	Map(from, to *object.Typ) (Caster, bool)
}

var Default = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
}
