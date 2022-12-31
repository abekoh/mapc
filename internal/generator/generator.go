package generator

import (
	"fmt"
	"github.com/abekoh/mapstructor/internal/model"
	"io"
)

type Generator struct {
}

func (g Generator) Generate(w io.Writer, from, to model.StructParam) error {
	fromStr, err := model.Lookup(from)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", fromStr, err)
	}
	toStr, err := model.Lookup(to)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", toStr, err)
	}
	//fmt.Printf("%+v, %+v", fromStr.String(), toStr.String())
	m := model.MatchModel(fromStr, toStr)
	tmplParam := model.NewTmplParam(m, "ab")
	if err := model.Write(w, tmplParam); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
