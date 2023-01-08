package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/pkg/fieldmapper"
	"github.com/abekoh/mapc/pkg/typemapper"
)

type fieldMap map[string]*object.Field

type FieldPair struct {
	From   *object.Field
	To     *object.Field
	Caster typemapper.Caster
}

type Mapper struct {
	FieldMappers []fieldmapper.FieldMapper
	TypeMappers  []typemapper.TypeMapper
}

type Mapping struct {
	From       *object.Struct
	To         *object.Struct
	FieldPairs []*FieldPair
}

func (m Mapper) NewMapping(from, to any) (*Mapping, error) {
	fromStr, err := object.NewStruct(reflect.TypeOf(from))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	toStr, err := object.NewStruct(reflect.TypeOf(to))
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

func (m Mapper) newFieldPairs(from, to *object.Struct) []*FieldPair {
	toFieldMap := make(fieldMap)
	for _, field := range to.Fields {
		toFieldMap[field.Name()] = field
	}

	var pairs []*FieldPair
	for _, fromField := range from.Fields {
		for _, fieldMapper := range m.FieldMappers {
			key := fieldMapper(fromField.Name())
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
	pair := &FieldPair{
		From: from,
		To:   to,
	}
	if from.IsSameTypeAndPkgWith(to) {
		return pair, true
	}
	panic("todo: impl")
}
