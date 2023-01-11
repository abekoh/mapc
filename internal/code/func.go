package code

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/util"
	"github.com/dave/dst"
)

type Typ struct {
	name    string
	pkgPath string
}

type Func struct {
	name     string
	argName  string
	fromTyp  *Typ
	toTyp    *Typ
	mapExprs MapExprList
}

type FuncOption struct {
	Name         string
	NameTemplate *template.Template
	Private      bool
	ArgName      string
}

func NewFuncFromMapping(m *mapping.Mapping, opt *FuncOption) *Func {
	if opt == nil {
		opt = &FuncOption{}
	}
	fieldMappers := MapExprList{}
	for _, p := range m.FieldPairs {
		fieldMappers = append(fieldMappers, &SimpleMapExpr{
			from: p.From.Name(),
			to:   p.To.Name(),
		})
	}
	return &Func{
		name:     funcName(m, opt),
		argName:  argName(m, opt),
		fromTyp:  &Typ{name: m.From.Name, pkgPath: m.From.PkgPath},
		toTyp:    &Typ{name: m.To.Name, pkgPath: m.To.PkgPath},
		mapExprs: fieldMappers,
	}
}

func funcName(m *mapping.Mapping, opt *FuncOption) (res string) {
	if opt.Name != "" {
		res = opt.Name
	} else if opt.NameTemplate != nil {
		var buf bytes.Buffer
		err := opt.NameTemplate.Execute(&buf, struct {
			From string
			To   string
		}{
			From: m.From.Name,
			To:   m.To.Name,
		})
		if err == nil {
			res = buf.String()
		}
	}
	if res == "" {
		res = fmt.Sprintf("To%s", util.UpperFirst(m.To.Name))
	}
	if opt.Private {
		res = util.LowerFirst(res)
	}
	return
}

func argName(m *mapping.Mapping, opt *FuncOption) string {
	if opt.ArgName != "" {
		return opt.ArgName
	}
	return "x"
}

func newFuncFromDecl(pkgPath string, d *dst.FuncDecl) (*Func, error) {
	getPkgPath := func(ident *dst.Ident) string {
		if ident.Path == "" {
			return pkgPath
		}
		return ident.Path
	}
	res := &Func{}
	res.name = d.Name.Name
	if len(d.Type.Params.List) != 1 {
		return nil, fmt.Errorf("length of params must be 1")
	}

	paramField := d.Type.Params.List[0]
	if len(paramField.Names) != 1 {
		return nil, fmt.Errorf("length of names must be 1")
	}
	res.argName = paramField.Names[0].Name
	paramTyp, ok := paramField.Type.(*dst.Ident)
	if !ok {
		return nil, fmt.Errorf("failed to cast param type")
	}
	res.fromTyp = &Typ{name: paramTyp.Name, pkgPath: getPkgPath(paramTyp)}

	if len(d.Type.Results.List) != 1 {
		return nil, fmt.Errorf("length of results must be 1")
	}
	resultField := d.Type.Results.List[0]
	resultTyp, ok := resultField.Type.(*dst.Ident)
	if !ok {
		return nil, fmt.Errorf("failed to cast result type")
	}
	res.toTyp = &Typ{name: resultTyp.Name, pkgPath: getPkgPath(resultTyp)}

	if len(d.Body.List) != 1 {
		return nil, fmt.Errorf("length of body must be 1")
	}
	body, ok := d.Body.List[0].(*dst.ReturnStmt)
	if !ok {
		return nil, fmt.Errorf("body must be ReturnStmt")
	}
	if len(body.Results) != 1 {
		return nil, fmt.Errorf("length of body.Results must be 1")
	}
	comps, ok := body.Results[0].(*dst.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("body.Results[0] must be *dst.CompositeLit")
	}
	res.mapExprs = MapExprList{}
	for _, expr := range comps.Elts {
		res.mapExprs = append(res.mapExprs, parseComments(expr.Decorations().Start.All()...)...)
		kvExpr, ok := expr.(*dst.KeyValueExpr)
		if !ok {
			return nil, fmt.Errorf("body.Results[*] must be KeyValueExpr")
		}
		keyIdent, ok := kvExpr.Key.(*dst.Ident)
		if !ok {
			return nil, fmt.Errorf("key must be Ident")
		}
		selectorExpr, ok := kvExpr.Value.(*dst.SelectorExpr)
		if !ok {
			return nil, fmt.Errorf("value must be SelectorExpr")
		}
		res.mapExprs = append(res.mapExprs, &SimpleMapExpr{
			from: selectorExpr.Sel.Name,
			to:   keyIdent.Name,
		})
		res.mapExprs = append(res.mapExprs, parseComments(expr.Decorations().End.All()...)...)
	}
	return res, nil
}

func (f Func) Decl() (*dst.FuncDecl, error) {
	return &dst.FuncDecl{
		Recv: nil,
		Name: genFuncName(f.name),
		Type: &dst.FuncType{
			Func:       false,
			TypeParams: nil,
			Params:     genParams(f.fromTyp, f.argName),
			Results:    genResult(f.toTyp),
			Decs:       dst.FuncTypeDecorations{},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				genReturn(f.toTyp, f.argName, f.mapExprs),
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
	o := dst.NewObj(dst.Typ, typ.name)
	i := dst.NewIdent(typ.name)
	i.Obj = o
	i.Path = typ.pkgPath
	return i
}

func genVar(name string) *dst.Ident {
	o := dst.NewObj(dst.Var, name)
	i := dst.NewIdent(name)
	i.Obj = o
	return i
}

func genParams(fromTyp *Typ, arg string) *dst.FieldList {
	return &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{genVar(arg)},
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

func genReturn(toTyp *Typ, arg string, exprs MapExprList) *dst.ReturnStmt {
	elts, comments := exprs.DstExprs(arg)
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
