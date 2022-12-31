package generator

import (
	"fmt"
	"github.com/abekoh/mapstructor/internal/exporter"
	"github.com/abekoh/mapstructor/internal/pivot"
	"io"
)

type Generator struct {
}

func (g Generator) Generate(w io.Writer, from, to pivot.StructParam) error {
	fromStr, err := pivot.Lookup(from)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", fromStr, err)
	}
	toStr, err := pivot.Lookup(to)
	if err != nil {
		return fmt.Errorf("failed to lookup %+v: %w", toStr, err)
	}
	//fmt.Printf("%+v, %+v", fromStr.String(), toStr.String())
	m := pivot.Match(fromStr, toStr)
	tmplParam := exporter.NewTmplParam(m, "ab")
	if err := exporter.Write(w, tmplParam); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
