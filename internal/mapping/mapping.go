package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/fieldmapper"
	"github.com/abekoh/mapc/internal/util"
	"github.com/abekoh/mapc/typemapper"
	"github.com/abekoh/mapc/types"
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
	From    *types.Field
	To      *types.Field
	Casters []typemapper.Caster
}

func (m Mapper) NewMapping(from, to any, outPkgPath string) (*Mapping, error) {
	fromStr, err := types.NewStruct(reflect.TypeOf(from))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	if !isAccessible(fromStr, outPkgPath) {
		return nil, fmt.Errorf("%v is not accessible from %v", fromStr.Name(), outPkgPath)
	}
	toStr, err := types.NewStruct(reflect.TypeOf(to))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	if !isAccessible(toStr, outPkgPath) {
		return nil, fmt.Errorf("%v is not accessible from %v", toStr.Name(), outPkgPath)
	}
	fieldPairs := m.newFieldPairs(fromStr, toStr, outPkgPath)
	return &Mapping{
		From:       fromStr,
		To:         toStr,
		FieldPairs: fieldPairs,
	}, nil
}

func isAccessible(id types.Identifier, outPkgPath string) bool {
	if !util.IsPrivate(id.Name()) {
		return true
	}
	return id.PkgPath() == outPkgPath
}

func (m Mapper) newFieldPairs(from, to *types.Struct, outPkgPath string) []*FieldPair {
	toFieldMap := make(map[string]*types.Field)
	for _, field := range to.Fields {
		toFieldMap[field.Name()] = field
	}

	var pairs []*FieldPair
	for _, fromField := range from.Fields {
		for _, fieldMapper := range m.FieldMappers {
			key := fieldMapper.Map(fromField.Name())
			if toField, ok := toFieldMap[key]; ok {
				if !isAccessible(fromField, outPkgPath) || !isAccessible(toField, outPkgPath) {
					break
				}
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
				From:    from,
				To:      to,
				Casters: []typemapper.Caster{caster},
			}, true
		}
	}
	return nil, false
}
