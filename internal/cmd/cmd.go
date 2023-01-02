package cmd

import (
	"fmt"
	"github.com/abekoh/mapc/internal/exporter"
	"github.com/abekoh/mapc/internal/mapping"
	"io"
)

type Param struct {
	Dir    string
	Pkg    string
	Struct string
}

func toStructParam(inp Param) mapping.StructParam {
	return mapping.StructParam{
		Dir:    inp.Dir,
		Pkg:    inp.Pkg,
		Struct: inp.Struct,
	}
}

func Generate(w io.Writer, from, to Param, distPkgName string) error {
	p, err := mapping.New(toStructParam(from), toStructParam(to), distPkgName)
	if err != nil {
		return fmt.Errorf("failed to construct mapping: %w", err)
	}
	tmplParam := exporter.NewTmplParam(p, distPkgName)
	if err := exporter.Run(w, tmplParam); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}
	return nil
}
