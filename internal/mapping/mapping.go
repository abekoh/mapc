package mapping

import (
	"fmt"
	"github.com/abekoh/mapc/internal/tokenizer"
	"go/types"
	"reflect"
)

type Var struct {
	v *types.Var
}

func (v Var) Name() string {
	return v.v.Name()
}

func newFieldPair(from, to Var) (FieldPair, bool) {
	pair := FieldPair{
		From: from,
		To:   to,
	}
	if from.v.Type().String() == to.v.Type().String() {
		return pair, true
	}
	switch from.v.Type().Underlying().(type) {
	case *types.Basic:
		fromBasic := from.v.Type().Underlying().(*types.Basic)
		toBasic := to.v.Type().Underlying().(*types.Basic)
		if isCastableBasicInfo(fromBasic.Info()) && fromBasic.Info() == toBasic.Info() {
			pair.Caster = &Caster{
				fmtString: fmt.Sprintf("%s(%%s)", to.v.Type().String()),
			}
			return pair, true
		}
		return pair, false
	default:
		return pair, false
	}
}

func isCastableBasicInfo(i types.BasicInfo) bool {
	return i&(types.IsBoolean|types.IsInteger|types.IsUnsigned|types.IsFloat|types.IsComplex) > 0
}

type tokenFieldMap map[string]Var

type Caster struct {
	fmtString string
}

type FieldPair struct {
	From   Var
	To     Var
	Caster *Caster
}

type Mapping struct {
	From        *PkgStruct
	To          *PkgStruct
	Maps        []FieldPair
	DistPkgName string
}

func newWithMatch(from, to *PkgStruct, distPkgName string) *Mapping {
	toTokenFieldMap := to.tokenFieldMap()
	maps := match(from, toTokenFieldMap)
	return &Mapping{
		From:        from,
		To:          to,
		Maps:        maps,
		DistPkgName: distPkgName,
	}
}

func match(from *PkgStruct, toTokenFieldMap tokenFieldMap) []FieldPair {
	var res []FieldPair
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		fromToken := tokenizer.Tokenize(fromField.Name())
		if toVar, ok := toTokenFieldMap[fromToken]; ok {
			if pair, ok := newFieldPair(Var{v: fromField}, toVar); ok {
				res = append(res, pair)
			}
		}
	}
	return res
}

func NewOld(from, to StructParam, distPkgName string) (*Mapping, error) {
	fromStr, err := newStruct(from)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup %+v: %w", fromStr, err)
	}
	toStr, err := newStruct(to)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup %+v: %w", toStr, err)
	}
	return newWithMatch(fromStr, toStr, distPkgName), nil
}

func New(from, to any) (*Mapping, error) {
	if reflect.TypeOf(from).Kind() != reflect.Struct {
		return nil, fmt.Errorf("from is must be struct, got %T", from)
	}
	if reflect.TypeOf(to).Kind() != reflect.Struct {
		return nil, fmt.Errorf("to is must be struct, got %T", to)
	}
	return nil, nil
}
