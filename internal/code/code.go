package code

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/abekoh/mapc/internal/mapping"
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

func New(pkgPath string) *File {
	df := &dst.File{
		Name:  dst.NewIdent(pkgName(pkgPath)),
		Decls: []dst.Decl{},
	}
	return &File{dstFile: df, pkgPath: pkgPath}
}

func (f File) Write(w io.Writer) error {
	r := decorator.NewRestorerWithImports(f.pkgPath, guess.New())
	if err := r.Fprint(w, f.dstFile); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (f File) FindFunc(name string) (*Func, bool) {
	for _, decl := range f.dstFile.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if ok && funcDecl.Name != nil && funcDecl.Name.Name == name {
			return &Func{fc: funcDecl}, true
		}
	}
	return nil, false
}

type Func struct {
	fc *dst.FuncDecl
}

func NewFunc(m *mapping.Mapping) *Func {
	return &Func{}
}

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
