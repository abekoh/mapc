package generator

import (
	"github.com/abekoh/mapstructor/internal/model"
	"io"
)

type Generator struct {
}

func (g Generator) Generate(w io.Writer, from, to model.StructParam) error {
	return nil
}
