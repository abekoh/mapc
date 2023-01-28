package code

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/util"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/dave/dst"
)

type Typ struct {
	name    string
	pkgPath string
}

func (t Typ) Equal(x *Typ) bool {
	return x.name == t.name && x.pkgPath == t.pkgPath
}

type Func struct {
	name            string
	argName         string
	srcTyp          *Typ
	destTyp         *Typ
	mapExprs        MapExprList
	withFuncComment bool
	editable        bool
}

func (f Func) Name() string {
	return f.name
}

func (f Func) funcComments() []string {
	comments := []string{
		"\n",
		fmt.Sprintf("// %s maps %s to %s.", f.name, f.srcTyp.name, f.destTyp.name),
		"// This function is generated by mapc.",
	}
	if f.editable {
		comments = append(comments, "// You can edit mapping fields.")
	} else {
		comments = append(comments, "// DO NOT EDIT this function.")
	}
	return comments
}

type FuncOption struct {
	Name                 string
	NameTemplate         *template.Template
	Private              bool
	ArgName              string
	FuncComment          bool
	NoMapperFieldComment bool
	Editable             bool
}

func NewFuncFromMapping(m *mapping.Mapping, opt *FuncOption) *Func {
	if opt == nil {
		opt = &FuncOption{}
	}

	mapExprs := MapExprList{}
	mappedFieldMap := make(map[string]struct{})
	for _, p := range m.FieldPairs {
		var casters []Caster
		for _, c := range p.Casters {
			caller := c.Caller()
			if caller == nil {
				continue
			}
			switch caller.CallerType {
			case mapcstd.Unary:
				casters = append(casters, UnaryCaster(caller.Name[0]))
			case mapcstd.Type:
				casters = append(casters, TypeCaster{
					name:    caller.Name,
					pkgPath: caller.PkgPath,
				})
			case mapcstd.Func:
				casters = append(casters, FuncCallCaster{
					name:    caller.Name,
					pkgPath: caller.PkgPath,
				})
			default:
			}
		}
		mapExprs = append(mapExprs, &SimpleMapExpr{
			src:     p.Src.Name(),
			dest:    p.Dest.Name(),
			casters: casters,
		})
		mappedFieldMap[p.Dest.Name()] = struct{}{}
	}
	if opt.NoMapperFieldComment {
		for _, destField := range m.Dest.Fields {
			destFieldName := destField.Name()
			if _, ok := mappedFieldMap[destFieldName]; !ok {
				mapExprs = append(mapExprs, &CommentedMapExpr{
					dest:    destFieldName,
					comment: fmt.Sprintf("// %s:", destFieldName),
				})
			}
		}
	}
	return &Func{
		name:            funcName(m, opt),
		argName:         argName(m, opt),
		srcTyp:          &Typ{name: m.Src.Name(), pkgPath: m.Src.PkgPath()},
		destTyp:         &Typ{name: m.Dest.Name(), pkgPath: m.Dest.PkgPath()},
		mapExprs:        mapExprs,
		withFuncComment: opt.FuncComment,
		editable:        opt.Editable,
	}
}

func funcName(m *mapping.Mapping, opt *FuncOption) (res string) {
	if opt.Name != "" {
		res = opt.Name
	} else if opt.NameTemplate != nil {
		var buf bytes.Buffer
		err := opt.NameTemplate.Execute(&buf, struct {
			Src  string
			Dest string
		}{
			Src:  m.Src.Name(),
			Dest: m.Dest.Name(),
		})
		if err == nil {
			res = buf.String()
		}
	}
	if res == "" {
		res = fmt.Sprintf("Map%sTo%s", util.UpperFirst(m.Src.Name()), util.UpperFirst(m.Dest.Name()))
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
	res.srcTyp = &Typ{name: paramTyp.Name, pkgPath: getPkgPath(paramTyp)}

	if len(d.Type.Results.List) != 1 {
		return nil, fmt.Errorf("length of results must be 1")
	}
	resultField := d.Type.Results.List[0]
	resultTyp, ok := resultField.Type.(*dst.Ident)
	if !ok {
		return nil, fmt.Errorf("failed to cast result type")
	}
	res.destTyp = &Typ{name: resultTyp.Name, pkgPath: getPkgPath(resultTyp)}

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
		// TODO: with caster
		res.mapExprs = append(res.mapExprs, &SimpleMapExpr{
			src:  selectorExpr.Sel.Name,
			dest: keyIdent.Name,
		})
		res.mapExprs = append(res.mapExprs, parseComments(expr.Decorations().End.All()...)...)
	}
	return res, nil
}

func (f Func) FillMapExprs(x *Func) (*Func, error) {
	if !f.srcTyp.Equal(x.srcTyp) {
		return nil, errors.New("srcTyp must be equal")
	}
	if !f.destTyp.Equal(x.destTyp) {
		return nil, errors.New("destTyp must be equal")
	}
	existedNormal, existedCommented := f.mapExprs.SeparateCommented()
	xNormal, xCommented := x.mapExprs.SeparateCommented()

	resMapExprs := make(MapExprList, 0)
	resMapExprsKeys := make(map[string]struct{})

	addExprs := func(exprList MapExprList) {
		for _, e := range exprList {
			if _, ok := resMapExprsKeys[e.Dest()]; !ok {
				resMapExprs = append(resMapExprs, e)
				resMapExprsKeys[e.Dest()] = struct{}{}
			}
		}
	}
	addExprs(existedNormal)
	addExprs(xNormal)
	addExprs(existedCommented)
	addExprs(xCommented)

	f.mapExprs = resMapExprs
	return &f, nil
}

func (f Func) Decl() (*dst.FuncDecl, error) {
	var fnComments []string
	if f.withFuncComment {
		fnComments = f.funcComments()
	}
	return &dst.FuncDecl{
		Recv: nil,
		Name: genFuncName(f.name),
		Type: &dst.FuncType{
			Func:       false,
			TypeParams: nil,
			Params:     genParams(f.srcTyp, f.argName),
			Results:    genResult(f.destTyp),
			Decs:       dst.FuncTypeDecorations{},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				genReturn(f.destTyp, f.argName, f.mapExprs),
			},
			RbraceHasNoPos: false,
			Decs:           dst.BlockStmtDecorations{},
		},
		Decs: dst.FuncDeclDecorations{
			NodeDecs: dst.NodeDecs{
				Start: fnComments,
			},
		},
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

func genTypeWithName(name, pkgPath string) *dst.Ident {
	o := dst.NewObj(dst.Typ, name)
	i := dst.NewIdent(name)
	i.Obj = o
	i.Path = pkgPath
	return i
}

func genFunc(name, pkgPath string) *dst.Ident {
	o := dst.NewObj(dst.Fun, name)
	i := dst.NewIdent(name)
	i.Obj = o
	i.Path = pkgPath
	return i
}

func genVar(name string) *dst.Ident {
	o := dst.NewObj(dst.Var, name)
	i := dst.NewIdent(name)
	i.Obj = o
	return i
}

func genParams(destTyp *Typ, arg string) *dst.FieldList {
	return &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{genVar(arg)},
				Type:  genType(destTyp),
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
		Decs: dst.CompositeLitDecorations{
			NodeDecs: dst.NodeDecs{},
			Lbrace:   comments,
		},
	}
	return &dst.ReturnStmt{
		Results: []dst.Expr{lit},
		Decs:    dst.ReturnStmtDecorations{},
	}
}
