package model

import "go/types"

type Model struct {
	from *Struct
	to   *Struct
	maps []FieldMap
}

func MatchModel(from, to *Struct) *Model {
	toFields := make(map[string]*types.Var)
	for i := 0; i < to.str.NumFields(); i++ {
		f := to.str.Field(i)
		// TODO: fix key
		toFields[f.Name()] = f
	}

	res := make([]FieldMap, 0)
	for i := 0; i < from.str.NumFields(); i++ {
		fromField := from.str.Field(i)
		// TODO: fix matching logic
		if toField, ok := toFields[fromField.Name()]; ok {
			res = append(res, FieldMap{fromField, toField})
		}
	}
	return &Model{
		from: from,
		to:   to,
		maps: res,
	}
}
