package code

import (
	"fmt"
	"go/token"
	"io"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
)

type File struct {
	pkgPath string
	dstFile *dst.File
}

func LoadFile(filePath, pkgPath string) (*File, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	dec := decorator.NewDecoratorWithImports(token.NewFileSet(), pkgPath, goast.New())
	df, err := dec.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	return &File{dstFile: df, pkgPath: pkgPath}, nil
}

func (f File) Write(w io.Writer) error {
	r := decorator.NewRestorerWithImports(f.pkgPath, guess.New())
	if err := r.Fprint(w, f.dstFile); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
