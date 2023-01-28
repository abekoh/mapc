package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/internal/str"
	"github.com/abekoh/mapc/internal/util"
	"github.com/abekoh/mapc/mapcstd"
)

type Mapper struct {
	FieldMappers []mapcstd.FieldMapper
	TypeMappers  []mapcstd.TypeMapper
}

type Mapping struct {
	Src        *str.Struct
	Dest       *str.Struct
	FieldPairs []*FieldPair
}

type FieldPair struct {
	Src     *str.Field
	Dest    *str.Field
	Casters []mapcstd.Caster
}

func (m Mapper) NewMapping(src, dest any, outPkgPath string) (*Mapping, error) {
	srcStr, err := str.NewStruct(reflect.TypeOf(src))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	if !isAccessible(srcStr, outPkgPath) {
		return nil, fmt.Errorf("%v is not accessible from %v", srcStr.Name(), outPkgPath)
	}
	destStr, err := str.NewStruct(reflect.TypeOf(dest))
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

func isAccessible(id str.Identifier, outPkgPath string) bool {
	if !util.IsPrivate(id.Name()) {
		return true
	}
	return id.PkgPath() == outPkgPath
}

func (m Mapper) newFieldPairs(src, dest *str.Struct, outPkgPath string) []*FieldPair {
	destFieldMap := make(map[string]*str.Field)
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

func (m Mapper) newFieldPair(src, dest *str.Field) (*FieldPair, bool) {
	for _, typeMapper := range m.TypeMappers {
		if caster, ok := typeMapper.Map(src.Typ(), dest.Typ()); ok {
			return &FieldPair{
				Src:     src,
				Dest:    dest,
				Casters: []mapcstd.Caster{caster},
			}, true
		}
	}
	return nil, false
}
