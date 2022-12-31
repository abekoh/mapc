package model

import "go/types"

type FieldPair struct {
	from *types.Var
	to   *types.Var
}

type Model struct {
	from *Struct
	to   *Struct
	maps []FieldPair
}

func MatchModel(from, to *Struct) *Model {
	toFieldMap := make(map[string]*types.Var)
	for i := 0; i < to.str.NumFields(); i++ {
		f := to.str.Field(i)
		// TODO: fix key
		toFieldMap[f.Name()] = f
	}

	res := make([]FieldPair, 0)
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		// TODO: fix matching logic
		if toField, ok := toFieldMap[fromField.Name()]; ok {
			res = append(res, FieldPair{fromField, toField})
		}
	}
	return &Model{
		from: from,
		to:   to,
		maps: res,
	}
}
