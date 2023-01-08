package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/internal/tokenizer"
)

type FieldMapper func(string) string

type Mapper struct {
	tzr *tokenizer.Tokenizer
}

func NewMapper() *Mapper {
	return &Mapper{tzr: tokenizer.NewTokenizer()}
}

type fieldMap map[string]*object.Field

func newFieldPair(from, to *object.Field) (*FieldPair, bool) {
	pair := &FieldPair{
		From: from,
		To:   to,
	}
	fromKind, toKind := from.Kind(), to.Kind()
	if fromKind == toKind {
		if from.PkgPath() != "" && from.PkgPath() == to.PkgPath() {
			pair.Caster = &Caster{
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

type Caster struct {
	pkgPath   string
	fmtString string
}

type FieldPair struct {
	From   *object.Field
	To     *object.Field
	Caster *Caster
}

type Mapping struct {
	From       *object.Struct
	To         *object.Struct
	FieldPairs []*FieldPair
}

func NewMapping(from, to any, fieldMappers []FieldMapper) (*Mapping, error) {
	fromStr, err := object.NewStruct(reflect.TypeOf(from))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	toStr, err := object.NewStruct(reflect.TypeOf(to))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	fieldPairs := newFieldPairs(fromStr, toStr, fieldMappers)
	return &Mapping{
		From:       fromStr,
		To:         toStr,
		FieldPairs: fieldPairs,
	}, nil
}

func newFieldPairs(from, to *object.Struct, fieldMappers []FieldMapper) []*FieldPair {
	toFieldMap := make(fieldMap)
	for _, field := range to.Fields {
		toFieldMap[field.Name()] = field
	}

	var pairs []*FieldPair
	for _, fromField := range from.Fields {
		for _, mapper := range fieldMappers {
			key := mapper(fromField.Name())
			if toField, ok := toFieldMap[key]; ok {
				if pair, ok := newFieldPair(fromField, toField); ok {
					pairs = append(pairs, pair)
					break
				}
			}
		}
	}
	return pairs
}
