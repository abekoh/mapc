package code

import (
	"errors"
	"fmt"
	"go/token"
	"io"
	"os"

	"github.com/abekoh/mapc/internal/util"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
)

type Mode int

const (
	PrioritizeGenerated Mode = iota
	PrioritizeExisted
	Deterministic
)

type File struct {
	pkgPath string
	dstFile *dst.File
}

func NewFile(pkgPath string) *File {
	df := &dst.File{
		Name:  dst.NewIdent(util.PkgNameFromPath(pkgPath)),
		Decls: []dst.Decl{},
	}
	return &File{dstFile: df, pkgPath: pkgPath}
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

func (f *File) Attach(inpFn *Func, mode Mode) error {
	i, existedFn, ok := f.FindFunc(inpFn.name)

	var resFn *Func
	var err error

	if !ok {
		resFn = inpFn
	}
	switch mode {
	case PrioritizeGenerated:
		resFn, err = inpFn.FillMapExprs(existedFn)
	case PrioritizeExisted:
		resFn, err = existedFn.FillMapExprs(inpFn)
	case Deterministic:
		resFn = inpFn
	default:
		return errors.New("invalid mode")
	}
	if err != nil {
		return fmt.Errorf("failed to fill MapExprs: %w", err)
	}
	fnDecl, err := resFn.Decl()
	if err != nil {
		return fmt.Errorf("failed to generate Decl: %w", err)
	}
	f.dstFile.Decls[i] = fnDecl
	return nil
}

func (f *File) Write(w io.Writer) error {
	r := decorator.NewRestorerWithImports(f.pkgPath, guess.New())
	if err := r.Fprint(w, f.dstFile); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (f *File) findFuncDecl(name string) (idx int, fn *dst.FuncDecl, ok bool) {
	for i, decl := range f.dstFile.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if ok && funcDecl.Name != nil && funcDecl.Name.Name == name {
			return i, funcDecl, true
		}
	}
	return -1, nil, false
}

func (f *File) FindFunc(name string) (idx int, fn *Func, ok bool) {
	i, d, ok := f.findFuncDecl(name)
	if !ok {
		return -1, nil, false
	}
	fn, err := newFuncFromDecl(f.pkgPath, d)
	if err != nil {
		return -1, nil, false
	}
	return i, fn, true
}
