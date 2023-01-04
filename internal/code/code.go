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

func emptyDecs() dst.IdentDecorations {
	return dst.IdentDecorations{
		NodeDecs: dst.NodeDecs{
			Before: dst.None,
			Start:  nil,
			End:    nil,
			After:  dst.None,
		},
		X: nil,
	}
}

func NewFunc(m *mapping.Mapping) *Func {
	//name := m.Name()
	fc := &dst.FuncDecl{
		Recv: nil,
		Name: &dst.Ident{
			Name: "ToBUser",
			Obj: &dst.Object{
				Kind: dst.Fun,
				Name: "ToBUser",
				Decl: &dst.FuncDecl{
					Recv: nil,
					Name: nil, // parent
					Type: nil,
					Body: nil,
					Decs: dst.FuncDeclDecorations{},
				},
				Type: nil,
			},
			Path: "",
			Decs: dst.IdentDecorations{},
		},
		Type: &dst.FuncType{
			Func:       true,
			TypeParams: nil,
			Params: &dst.FieldList{
				Opening: true,
				List: []*dst.Field{
					{
						Names: []*dst.Ident{
							{
								Name: "inp",
								Obj: &dst.Object{
									Kind: dst.Var,
									Name: "inp",
									Decl: nil, // will be inserted
									Data: nil,
									Type: nil,
								},
								Path: "",
								Decs: dst.IdentDecorations{},
							},
						},
						Type: &dst.Ident{
							Name: "AUser",
							Obj: &dst.Object{
								Kind: dst.Typ,
								Name: "AUser",
								Decl: &dst.TypeSpec{
									Name:       nil, // replace
									TypeParams: nil,
									Assign:     false,
									Type:       nil,
									Decs:       dst.TypeSpecDecorations{},
								},
								Data: nil,
								Type: nil,
							},
							Path: "",
							Decs: dst.IdentDecorations{},
						},
						Tag:  nil,
						Decs: dst.FieldDecorations{},
					},
				},
				Closing: true,
				Decs:    dst.FieldListDecorations{},
			},
			Results: &dst.FieldList{
				Opening: false,
				List: []*dst.Field{
					{
						Names: nil,
						Type: &dst.Ident{
							Name: "BUser",
							Obj: &dst.Object{
								Kind: dst.Typ,
								Name: "BUser",
								Decl: nil, // insert
								Data: nil,
								Type: nil,
							},
							Path: "",
							Decs: dst.IdentDecorations{},
						},
						Tag:  nil,
						Decs: dst.FieldDecorations{},
					},
				},
				Closing: false,
				Decs:    dst.FieldListDecorations{},
			},
			Decs: dst.FuncTypeDecorations{},
		},
		Body: nil,
		Decs: dst.FuncDeclDecorations{},
	}
	return &Func{fc: fc}
}

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
