package function

import (
	"fmt"
	"strings"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/util"
	"github.com/dave/dst"
)

type Typ struct {
	pkgPath string
}

func (t Typ) Name() string {
	sp := strings.Split(t.pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

type Caster struct {
	fc string
}

type Function struct {
	name     string
	argName  string
	fromTyp  *Typ
	toTyp    *Typ
	withErr  bool
	mapExprs FieldMapperList
}

func NewFunction(m *mapping.Mapping) *Function {
	var exprs FieldMapperList
	for _, p := range m.FieldPairs {
		exprs = append(exprs, &SimpleFieldMapper{
			from: p.From.Name(),
			to:   p.To.Name(),
		})
	}
	return &Function{
		name:     fmt.Sprintf("To%s", util.UpperFirst(m.To.Name)),
		argName:  "from",
		fromTyp:  &Typ{pkgPath: m.From.PkgPath},
		toTyp:    &Typ{pkgPath: m.To.PkgPath},
		withErr:  false,
		mapExprs: exprs,
	}
}

func (f Function) Decl() (*dst.FuncDecl, error) {
	argIdent := genVar(f.argName)
	return &dst.FuncDecl{
		Recv: nil,
		Name: genFuncName(f.name),
		Type: &dst.FuncType{
			Func:       false,
			TypeParams: nil,
			Params:     genParams(f.fromTyp, argIdent),
			Results:    genResult(f.toTyp),
			Decs:       dst.FuncTypeDecorations{},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				genReturn(f.toTyp, argIdent, f.mapExprs),
			},
			RbraceHasNoPos: false,
			Decs:           dst.BlockStmtDecorations{},
		},
		Decs: dst.FuncDeclDecorations{},
	}, nil
}

func genFuncName(s string) *dst.Ident {
	o := dst.NewObj(dst.Fun, s)
	i := dst.NewIdent(s)
	i.Obj = o
	return i
}

func genType(typ *Typ) *dst.Ident {
	o := dst.NewObj(dst.Typ, typ.Name())
	i := dst.NewIdent(typ.Name())
	i.Obj = o
	return i
}

func cloneIdent(i *dst.Ident) *dst.Ident {
	return dst.Clone(i).(*dst.Ident)
}

func genVar(name string) *dst.Ident {
	o := dst.NewObj(dst.Var, name)
	i := dst.NewIdent(name)
	i.Obj = o
	return i
}

func genParams(fromTyp *Typ, argIdent *dst.Ident) *dst.FieldList {
	return &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{cloneIdent(argIdent)},
				Type:  genType(fromTyp),
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: true,
		Decs:    dst.FieldListDecorations{},
	}
}

func genResult(toTyp *Typ) *dst.FieldList {
	return &dst.FieldList{
		Opening: false,
		List: []*dst.Field{
			{
				Names: nil,
				Type:  genType(toTyp),
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: false,
		Decs:    dst.FieldListDecorations{},
	}
}

func genReturn(toTyp *Typ, argIdent *dst.Ident, exprs FieldMapperList) *dst.ReturnStmt {
	elts, comments := exprs.DstExprs(argIdent)
	lit := &dst.CompositeLit{
		Type:       genType(toTyp),
		Elts:       elts,
		Incomplete: false,
		Decs:       dst.CompositeLitDecorations{},
	}
	lit.Decorations().Start.Append(comments...)
	return &dst.ReturnStmt{
		Results: []dst.Expr{lit},
		Decs:    dst.ReturnStmtDecorations{},
	}
}
