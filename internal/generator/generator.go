package generator

import (
	"fmt"
	"github.com/abekoh/mapstructor/internal/exporter"
	"github.com/abekoh/mapstructor/internal/pivot"
	"io"
)

type Generator struct {
}

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

func (g Generator) Generate(w io.Writer, from, to Param) error {
	p, err := pivot.New(toStructParam(from), toStructParam(to))
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
