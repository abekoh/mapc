package pivot

import (
	"fmt"
	"github.com/abekoh/mapc/internal/tokenizer"
	"go/types"
)

type Var struct {
	v *types.Var
}

func (v Var) Name() string {
	return v.v.Name()
}

func newFieldPair(from, to Var) (FieldPair, bool) {
	// TODO: logic
	return FieldPair{
		From:   from,
		To:     to,
		Caster: nil,
	}, true
}

type tokenFieldMap map[string]Var

type Caster struct {
	fc        func(any) any
	fmtString string
}

type FieldPair struct {
	From   Var
	To     Var
	Caster *Caster
}

type Pivot struct {
	From        *Struct
	To          *Struct
	Maps        []FieldPair
	DistPkgName string
}

func newWithMatch(from, to *Struct, distPkgName string) *Pivot {
	toTokenFieldMap := to.tokenFieldMap()
	maps := match(from, toTokenFieldMap)
	return &Pivot{
		From:        from,
		To:          to,
		Maps:        maps,
		DistPkgName: distPkgName,
	}
}

func match(from *Struct, toTokenFieldMap tokenFieldMap) []FieldPair {
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

func New(from, to StructParam, distPkgName string) (*Pivot, error) {
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
