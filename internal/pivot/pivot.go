package pivot

import (
	"fmt"
	"github.com/abekoh/mapc/internal/tokenizer"
	"go/types"
)

type Var struct {
	v *types.Var
}

func (o Var) Name() string {
	return o.v.Name()
}

type tokenFieldMap map[string]Var

type FieldPair struct {
	From Var
	To   Var
}

type Pivot struct {
	From        *Struct
	To          *Struct
	Maps        []FieldPair
	DistPkgName string
}

func newWithMatch(from, to *Struct, distPkgName string) *Pivot {
	toTokenFieldMap := to.tokenFieldMap()
	res := make([]FieldPair, 0)
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		fromToken := tokenizer.Tokenize(fromField.Name())
		// TODO: fix matching logic
		if toObj, ok := toTokenFieldMap[fromToken]; ok {
			res = append(res, FieldPair{Var{fromField}, toObj})
		}
	}
	return &Pivot{
		From:        from,
		To:          to,
		Maps:        res,
		DistPkgName: distPkgName,
	}
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
