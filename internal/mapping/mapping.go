package mapping

import (
	"fmt"
	"reflect"

	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/internal/tokenizer"
	"github.com/abekoh/mapc/internal/util"
)

//type Var struct {
//	v *types.Var
//}
//
//function (v Var) Name() string {
//	return v.v.Name()
//}
//
//type Var2 struct {
//	v *types.Var
//}

//function newFieldPair(from, to Var) (FieldPairOld, bool) {
//	pair := FieldPairOld{
//		From: from,
//		To:   to,
//	}
//	if from.v.Type().String() == to.v.Type().String() {
//		return pair, true
//	}
//	switch from.v.Type().Underlying().(type) {
//	case *types.Basic:
//		fromBasic := from.v.Type().Underlying().(*types.Basic)
//		toBasic := to.v.Type().Underlying().(*types.Basic)
//		if isCastableBasicInfo(fromBasic.Info()) && fromBasic.Info() == toBasic.Info() {
//			pair.Caster = &Caster{
//				fmtString: fmt.Sprintf("%s(%%s)", to.v.Type().String()),
//			}
//			return pair, true
//		}
//		return pair, false
//	default:
//		return pair, false
//	}
//}
//
//function isCastableBasicInfo(i types.BasicInfo) bool {
//	return i&(types.IsBoolean|types.IsInteger|types.IsUnsigned|types.IsFloat|types.IsComplex) > 0
//}
//
//type tokenFieldMapOld map[string]Var

// new

type Mapper struct {
	tzr *tokenizer.Tokenizer
}

func NewMapper() *Mapper {
	return &Mapper{tzr: tokenizer.NewTokenizer()}
}

type tokenFieldMap map[string]*object.Field

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

func (m Mapping) Name() string {
	// TODO: fix
	return fmt.Sprintf("To%s", util.UpperFirst(m.To.Name))
}

func (m Mapper) Map(from, to any) (*Mapping, error) {
	fromStr, err := object.NewStruct(reflect.TypeOf(from))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	toStr, err := object.NewStruct(reflect.TypeOf(to))
	if err != nil {
		return nil, fmt.Errorf("failed to construct struct: %w", err)
	}
	toTokenFieldMap := m.newTokenFieldMap(toStr)
	fieldPairs := m.newFieldPairs(fromStr, toTokenFieldMap)
	return &Mapping{
		From:       fromStr,
		To:         toStr,
		FieldPairs: fieldPairs,
	}, nil
}

func (m Mapper) newTokenFieldMap(s *object.Struct) tokenFieldMap {
	r := make(tokenFieldMap)
	for _, f := range s.Fields {
		token := m.tzr.Tokenize(f.Name())
		r[token] = f
	}
	return r
}
func (m Mapper) newFieldPairs(from *object.Struct, toTokenFieldMap tokenFieldMap) []*FieldPair {
	var r []*FieldPair
	for _, fromF := range from.Fields {
		fromToken := m.tzr.Tokenize(fromF.Name())
		if toF, ok := toTokenFieldMap[fromToken]; ok {
			if pair, ok := newFieldPair(fromF, toF); ok {
				r = append(r, pair)
			}
		}
	}
	return r
}
