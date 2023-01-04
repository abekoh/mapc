package code

import (
	"bytes"
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

func NewFile(pkgPath string) *File {
	df := &dst.File{
		Name:  dst.NewIdent(pkgName(pkgPath)),
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

func (f File) FindFunc(name string) (*Func, bool) {
	if d, ok := f.findFuncDecl(name); ok {
		return &Func{fc: d}, true
	}
	return nil, false
}

func (f *File) Apply(fc *Func) {
	existed, ok := f.findFuncDecl(fc.Name())
	if ok {
		existed = fc.fc
	} else {
		_ = existed
		f.dstFile.Decls = append(f.dstFile.Decls, fc.fc)
	}
}

func (f File) Write(w io.Writer) error {
	r := decorator.NewRestorerWithImports(f.pkgPath, guess.New())
	if err := r.Fprint(w, f.dstFile); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (f File) sPrint() (string, error) {
	var buf bytes.Buffer
	err := f.Write(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to write to buffer: %w", err)
	}
	return buf.String(), nil
}

func (f File) findFuncDecl(name string) (*dst.FuncDecl, bool) {
	for _, decl := range f.dstFile.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if ok && funcDecl.Name != nil && funcDecl.Name.Name == name {
			return funcDecl, true
		}
	}
	return nil, false
}
