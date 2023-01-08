package function

import (
	"fmt"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/util"
	"github.com/dave/dst"
)

type Typ struct {
	name    string
	pkgPath string
}

func (t Typ) Name() string {
	//sp := strings.Split(t.pkgPath, "/")
	//if len(sp) == 0 {
	//	return ""
	//}
	//return sp[len(sp)-1]
	return t.name
}

type Caster struct {
	fc string
}

type Function struct {
	name     string
	argName  string
	fromTyp  *Typ
	toTyp    *Typ
	mapExprs FieldMapperList
}

func (f Function) Name() string {
	return f.name
}

func NewFromMapping(m *mapping.Mapping) *Function {
	var fieldMappers FieldMapperList
	for _, p := range m.FieldPairs {
		fieldMappers = append(fieldMappers, &SimpleFieldMapper{
			from: p.From.Name(),
			to:   p.To.Name(),
		})
	}
	return &Function{
		name:     fmt.Sprintf("To%s", util.UpperFirst(m.To.Name)),
		argName:  "from",
		fromTyp:  &Typ{name: m.From.Name, pkgPath: m.From.PkgPath},
		toTyp:    &Typ{name: m.To.Name, pkgPath: m.To.PkgPath},
		mapExprs: fieldMappers,
	}
}

func NewFromDecl(pkgPath string, d *dst.FuncDecl) (*Function, error) {
	getPkgPath := func(ident *dst.Ident) string {
		if ident.Path == "" {
			return pkgPath
		}
		return ident.Path
	}
	res := &Function{}
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
	res.fromTyp = &Typ{pkgPath: getPkgPath(paramTyp)}

	if len(d.Type.Results.List) != 1 {
		return nil, fmt.Errorf("length of results must be 1")
	}
	resultField := d.Type.Results.List[0]
	resultTyp, ok := resultField.Type.(*dst.Ident)
	if !ok {
		return nil, fmt.Errorf("failed to cast result type")
	}
	res.toTyp = &Typ{pkgPath: getPkgPath(resultTyp)}

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
	res.mapExprs = FieldMapperList{}
	for _, expr := range body.Results {
		res.mapExprs = append(res.mapExprs, ParseComments(expr.Decorations().Start.All()...)...)
		kvExpr, ok := body.Results[0].(*dst.KeyValueExpr)
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
		res.mapExprs = append(res.mapExprs, &SimpleFieldMapper{
			from: selectorExpr.Sel.Name,
			to:   keyIdent.Name,
		})
		res.mapExprs = append(res.mapExprs, ParseComments(expr.Decorations().End.All()...)...)
	}
	return res, nil
}

func (f Function) Decl() (*dst.FuncDecl, error) {
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
	o := dst.NewObj(dst.Typ, typ.Name())
	i := dst.NewIdent(typ.Name())
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

func genReturn(toTyp *Typ, arg string, exprs FieldMapperList) *dst.ReturnStmt {
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
