package mapc

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/internal/object"
)

type fieldMap map[string]*object.Field

type FieldPair struct {
	From   *object.Field
	To     *object.Field
	Caster Caster
}

type Mapper struct {
	FieldMappers []FieldMapper
	TypeMappers  []TypeMapper
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
				if pair, ok := m.newFieldPairOLD(fromField, toField); ok {
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

func (m Mapper) newFieldPairOLD(from, to *object.Field) (*FieldPair, bool) {
	pair := &FieldPair{
		From: from,
		To:   to,
	}
	fromKind, toKind := from.Kind(), to.Kind()
	if fromKind == toKind {
		if from.PkgPath() != "" && from.PkgPath() == to.PkgPath() {
			pair.Caster = Caster{
				pkgPath:   from.PkgPath(),
				fmtString: fmt.Sprintf("%s(%%s)", from.TypeName()),
			}
		}
		return pair, true
	}
	switch fromKind {
	case reflect.Int:
		switch toKind {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			pair.Caster = &Caster{
				fmtString: fmt.Sprintf("%s(%%s)", toKind.String()),
			}
			return pair, true
		default:
		}
	default:
	}
	return pair, false
}
