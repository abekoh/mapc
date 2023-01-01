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
	p, err := pivot.New(from, to)
	if err != nil {
		return fmt.Errorf("failed to construct pivot: %w", err)
	}
	// TODO: set dstPkgName
	tmplParam := exporter.NewTmplParam(p, "ab")
	if err := exporter.Run(w, tmplParam); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}
	return nil
}
