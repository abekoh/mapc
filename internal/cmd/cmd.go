package cmd

import (
	"fmt"
	"github.com/abekoh/mapc/internal/exporter"
	"github.com/abekoh/mapc/internal/pivot"
	"io"
)

type Param struct {
	Dir    string
	Pkg    string
	Struct string
}

func toStructParam(inp Param) pivot.StructParam {
	return pivot.StructParam{
		Dir:    inp.Dir,
		Pkg:    inp.Pkg,
		Struct: inp.Struct,
	}
}

func Generate(w io.Writer, from, to Param, distPkgName string) error {
	p, err := pivot.New(toStructParam(from), toStructParam(to), distPkgName)
	if err != nil {
		return fmt.Errorf("failed to construct pivot: %w", err)
	}
	tmplParam := exporter.NewTmplParam(p, distPkgName)
	if err := exporter.Run(w, tmplParam); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}
	return nil
}
