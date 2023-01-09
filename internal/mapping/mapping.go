package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/fieldmapper"
	"github.com/abekoh/mapc/internal/types"
	"github.com/abekoh/mapc/typemapper"
)

type Mapper struct {
	FieldMappers []fieldmapper.FieldMapper
	TypeMappers  []typemapper.TypeMapper
}

type Mapping struct {
	From       *types.Struct
	To         *types.Struct
	FieldPairs []*FieldPair
}

type FieldPair struct {
	From   *types.Field
	To     *types.Field
	Caster typemapper.Caster
}

func (m Mapper) NewMapping(from, to any) (*Mapping, error) {
	fromStr, err := types.NewStruct(reflect.TypeOf(from))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	toStr, err := types.NewStruct(reflect.TypeOf(to))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	fieldPairs := m.newFieldPairs(fromStr, toStr)
	return &Mapping{
		From:       fromStr,
		To:         toStr,
		FieldPairs: fieldPairs,
	}, nil
}

func (m Mapper) newFieldPairs(from, to *types.Struct) []*FieldPair {
	toFieldMap := make(map[string]*types.Field)
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

func (m Mapper) newFieldPair(from, to *types.Field) (*FieldPair, bool) {
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
