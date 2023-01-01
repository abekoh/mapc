package pivot

import (
	"fmt"
	"go/types"
)

type Object struct {
	v *types.Var
}

func (o Object) Name() string {
	return o.v.Name()
}

type FieldPair struct {
	From Object
	To   Object
}

type Pivot struct {
	From *Struct
	To   *Struct
	Maps []FieldPair
}

func match(from, to *Struct) *Pivot {
	toFieldMap := make(map[string]Object)
	for i := 0; i < to.str.NumFields(); i++ {
		f := to.str.Field(i)
		// TODO: fix key
		toFieldMap[f.Name()] = Object{v: f}
	}

	res := make([]FieldPair, 0)
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		// TODO: fix matching logic
		if toObj, ok := toFieldMap[fromField.Name()]; ok {
			res = append(res, FieldPair{Object{fromField}, toObj})
		}
	}
	return &Pivot{
		From: from,
		To:   to,
		Maps: res,
	}
}

func New(from, to StructParam) (*Pivot, error) {
	fromStr, err := newStruct(from)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup %+v: %w", fromStr, err)
	}
	toStr, err := newStruct(to)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup %+v: %w", toStr, err)
	}
	return match(fromStr, toStr), nil
}
