package mapping

import (
	"github.com/abekoh/mapc/fieldmapper"
	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/typemapper"
)

type (
	Mapper struct {
		FieldMappers []fieldmapper.FieldMapper
		TypeMappers  []typemapper.TypeMapper
	}

	Mapping struct {
		From       *object.Struct
		To         *object.Struct
		FieldPairs []*FieldPair
	}

	FieldPair struct {
		From   *object.Field
		To     *object.Field
		Caster typemapper.Caster
	}
)

func (m Mapper) newFieldPairs(from, to *object.Struct) []*FieldPair {
	toFieldMap := make(map[string]*object.Field)
	for _, field := range to.Fields {
		toFieldMap[field.Name()] = field
	}

	var pairs []*FieldPair
	for _, fromField := range from.Fields {
		for _, fieldMapper := range m.FieldMappers {
			key := fieldMapper.Map(fromField.Name())
			if toField, ok := toFieldMap[key]; ok {
				if pair, ok := m.newFieldPair(fromField, toField); ok {
					pairs = append(pairs, pair)
					break
				}
			}
		}
	}
	return pairs
}

func (m Mapper) newFieldPair(from, to *object.Field) (*FieldPair, bool) {
	for _, typeMapper := range m.TypeMappers {
		if caster, ok := typeMapper.Map(from.Typ(), to.Typ()); ok {
			return &FieldPair{
				From:   from,
				To:     to,
				Caster: caster,
			}, true
		}
	}
	return nil, false
}
