package pivot

import (
	"fmt"
	"go/types"
)

type Var struct {
	v *types.Var
}

func (o Var) Name() string {
	return o.v.Name()
}

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
	toFieldMap := make(map[string]Var)
	for i := 0; i < to.str.NumFields(); i++ {
		f := to.str.Field(i)
		// TODO: fix key
		toFieldMap[f.Name()] = Var{v: f}
	}

	res := make([]FieldPair, 0)
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		// TODO: fix matching logic
		if toObj, ok := toFieldMap[fromField.Name()]; ok {
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
