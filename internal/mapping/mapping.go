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
	Src        *types.Struct
	Dest       *types.Struct
	FieldPairs []*FieldPair
}

type FieldPair struct {
	From    *types.Field
	To      *types.Field
	Casters []typemapper.Caster
}

func (m Mapper) NewMapping(src, dest any, outPkgPath string) (*Mapping, error) {
	srcStr, err := types.NewStruct(reflect.TypeOf(src))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	if !isAccessible(srcStr, outPkgPath) {
		return nil, fmt.Errorf("%v is not accessible from %v", srcStr.Name(), outPkgPath)
	}
	destStr, err := types.NewStruct(reflect.TypeOf(dest))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	if !isAccessible(destStr, outPkgPath) {
		return nil, fmt.Errorf("%v is not accessible from %v", destStr.Name(), outPkgPath)
	}
	fieldPairs := m.newFieldPairs(srcStr, destStr, outPkgPath)
	return &Mapping{
		Src:        srcStr,
		Dest:       destStr,
		FieldPairs: fieldPairs,
	}, nil
}

func isAccessible(id types.Identifier, outPkgPath string) bool {
	if !util.IsPrivate(id.Name()) {
		return true
	}
	return id.PkgPath() == outPkgPath
}

func (m Mapper) newFieldPairs(src, dest *types.Struct, outPkgPath string) []*FieldPair {
	destFieldMap := make(map[string]*types.Field)
	for _, field := range dest.Fields {
		destFieldMap[field.Name()] = field
	}

	var pairs []*FieldPair
	for _, srcField := range src.Fields {
		for _, fieldMapper := range m.FieldMappers {
			key := fieldMapper.Map(srcField.Name())
			if destField, ok := destFieldMap[key]; ok {
				if !isAccessible(srcField, outPkgPath) || !isAccessible(destField, outPkgPath) {
					break
				}
				if pair, ok := m.newFieldPair(srcField, destField); ok {
					pairs = append(pairs, pair)
					break
				}
			}
		}
	}
	return pairs
}

func (m Mapper) newFieldPair(src, dest *types.Field) (*FieldPair, bool) {
	for _, typeMapper := range m.TypeMappers {
		if caster, ok := typeMapper.Map(src.Typ(), dest.Typ()); ok {
			return &FieldPair{
				From:    src,
				To:      dest,
				Casters: []typemapper.Caster{caster},
			}, true
		}
	}
	return nil, false
}
