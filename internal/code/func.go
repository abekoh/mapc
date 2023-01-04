package code

import (
	"errors"
	"regexp"
	"strings"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/object"
	"github.com/dave/dst"
)

type Func struct {
	fc *dst.FuncDecl
}

func (f Func) Name() string {
	if f.fc.Name == nil {
		return ""
	}
	return f.fc.Name.Name
}

func NewFuncSand(m *mapping.Mapping) *Func {
	//name := m.Name()
	funcNameObj := dst.NewObj(dst.Fun, "ToBUser")
	funcName := dst.NewIdent("ToBUser")
	funcName.Obj = funcNameObj

	inpObj := dst.NewObj(dst.Var, "inp")
	inp := dst.NewIdent("inp")
	inp.Obj = inpObj
	inpTypObj := dst.NewObj(dst.Typ, "AUser")
	inpTyp := dst.NewIdent("AUser")
	inpTyp.Obj = inpTypObj

	params := &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{inp},
				Type:  inpTyp,
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: true,
		Decs:    dst.FieldListDecorations{},
	}

	retTypObj := dst.NewObj(dst.Typ, "BUser")
	retTyp := dst.NewIdent("BUser")
	retTyp.Obj = retTypObj

	results := &dst.FieldList{
		Opening: false,
		List: []*dst.Field{
			{
				Names: nil,
				Type:  retTyp,
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: false,
		Decs:    dst.FieldListDecorations{},
	}

	typ := &dst.FuncType{
		Func:       false,
		TypeParams: nil,
		Params:     params,
		Results:    results,
		Decs:       dst.FuncTypeDecorations{},
	}

	el1Key := dst.NewIdent("ID")
	el1Sel := dst.NewIdent("ID")
	el1 := &dst.KeyValueExpr{
		Key: el1Key,
		Value: &dst.SelectorExpr{
			X:    dst.Clone(inp).(dst.Expr),
			Sel:  el1Sel,
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{},
	}

	el2Key := dst.NewIdent("Name")
	el2Sel := dst.NewIdent("Name")
	el2 := &dst.KeyValueExpr{
		Key: el2Key,
		Value: &dst.SelectorExpr{
			X:    dst.Clone(inp).(dst.Expr),
			Sel:  el2Sel,
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{},
	}

	el3Key := dst.NewIdent("Age")
	el3Sel := dst.NewIdent("Age")
	el3 := &dst.KeyValueExpr{
		Key: el3Key,
		Value: &dst.SelectorExpr{
			X:    dst.Clone(inp).(dst.Expr),
			Sel:  el3Sel,
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{},
	}

	el4Key := dst.NewIdent("RegisteredAt")
	el4Sel := dst.NewIdent("RegisteredAt")
	el4 := &dst.KeyValueExpr{
		Key: el4Key,
		Value: &dst.SelectorExpr{
			X:    dst.Clone(inp).(dst.Expr),
			Sel:  el4Sel,
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{},
	}

	ret := &dst.ReturnStmt{
		Results: []dst.Expr{
			&dst.CompositeLit{
				Type: dst.Clone(retTyp).(dst.Expr),
				Elts: []dst.Expr{
					el1,
					el2,
					el3,
					el4,
				},
				Incomplete: false,
				Decs:       dst.CompositeLitDecorations{},
			},
		},
		Decs: dst.ReturnStmtDecorations{},
	}

	body := &dst.BlockStmt{
		List: []dst.Stmt{
			ret,
		},
		RbraceHasNoPos: false,
		Decs:           dst.BlockStmtDecorations{},
	}

	fc := &dst.FuncDecl{
		Recv: nil,
		Name: funcName,
		Type: typ,
		Body: body,
		Decs: dst.FuncDeclDecorations{},
	}
	return &Func{fc: fc}
}

func NewFunc(m *mapping.Mapping) *Func {
	argIdent := genVar("from")
	return &Func{
		fc: &dst.FuncDecl{
			Recv: nil,
			Name: genFuncName(m),
			Type: &dst.FuncType{
				Func:       false,
				TypeParams: nil,
				Params:     genParams(m.From, argIdent),
				Results:    genResult(m.To),
				Decs:       dst.FuncTypeDecorations{},
			},
			Body: &dst.BlockStmt{
				List: []dst.Stmt{
					genReturn(m, argIdent),
				},
				RbraceHasNoPos: false,
				Decs:           dst.BlockStmtDecorations{},
			},
			Decs: dst.FuncDeclDecorations{},
		},
	}
}

func cloneIdent(i *dst.Ident) *dst.Ident {
	return dst.Clone(i).(*dst.Ident)
}

func genFuncName(m *mapping.Mapping) *dst.Ident {
	name := m.Name()
	o := dst.NewObj(dst.Fun, name)
	i := dst.NewIdent(name)
	i.Obj = o
	return i
}

func genParams(fromStr *object.Struct, argIdent *dst.Ident) *dst.FieldList {
	return &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{cloneIdent(argIdent)},
				Type:  genType(fromStr),
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: true,
		Decs:    dst.FieldListDecorations{},
	}
}

func genResult(toStr *object.Struct) *dst.FieldList {
	return &dst.FieldList{
		Opening: false,
		List: []*dst.Field{
			{
				Names: nil,
				Type:  genType(toStr),
				Tag:   nil,
				Decs:  dst.FieldDecorations{},
			},
		},
		Closing: false,
		Decs:    dst.FieldListDecorations{},
	}
}

func genType(str *object.Struct) *dst.Ident {
	o := dst.NewObj(dst.Typ, str.Name)
	i := dst.NewIdent(str.Name)
	i.Obj = o
	return i
}

func genVar(name string) *dst.Ident {
	o := dst.NewObj(dst.Var, name)
	i := dst.NewIdent(name)
	i.Obj = o
	return i
}

func genReturn(m *mapping.Mapping, argIdent *dst.Ident) *dst.ReturnStmt {
	var elts []dst.Expr
	for _, fp := range m.FieldPairs {
		elts = append(elts, genElt(fp, cloneIdent(argIdent)))
	}
	return &dst.ReturnStmt{
		Results: []dst.Expr{
			&dst.CompositeLit{
				Type:       genType(m.To),
				Elts:       elts,
				Incomplete: false,
				Decs:       dst.CompositeLitDecorations{},
			},
		},
		Decs: dst.ReturnStmtDecorations{},
	}
}

func genElt(fp *mapping.FieldPair, from *dst.Ident) *dst.KeyValueExpr {
	return &dst.KeyValueExpr{
		Key: dst.NewIdent(fp.To.Name()),
		Value: &dst.SelectorExpr{
			X:    from,
			Sel:  dst.NewIdent(fp.From.Name()),
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.NewLine,
				After:  dst.NewLine,
			},
		},
	}
}

func (f Func) ReNew(m *mapping.Mapping) (*Func, error) {
	decl := dst.Clone(f.fc).(*dst.FuncDecl)
	if decl.Body == nil || decl.Body.List == nil || len(decl.Body.List) == 0 {
		return nil, errors.New("failed to find return statement")
	}
	retStmt, ok := decl.Body.List[len(decl.Body.List)-1].(*dst.ReturnStmt)
	if !ok {
		return nil, errors.New("failed cast to *dst.ReturnStmt")
	}
	if retStmt.Results == nil || len(retStmt.Results) == 0 {
		return nil, errors.New("failed to get return results")
	}
	compLit, ok := retStmt.Results[0].(*dst.CompositeLit)
	if !ok {
		return nil, errors.New("failed to cast *dst.CompositeLit")
	}
	feMap := newFieldExprs(compLit.Elts)
	var resFieldExprs []*fieldExpr
	for _, pair := range m.FieldPairs {
		// TODO: get 'from ident' from params
		fromIdent := genVar("from")
		resFieldExprs = append(resFieldExprs, &fieldExpr{
			expr: genElt(pair, fromIdent),
		})

		key := pair.To.Name()
		if _, ok := feMap.m[key]; ok {
			delete(feMap.m, key)
		}
	}
	for _, e := range feMap.sortedFieldExprList() {
		resFieldExprs = append(resFieldExprs, e)
	}
	var resExprs []dst.Expr
	var prevExpr *dst.KeyValueExpr
	var nextComment string
	for _, resFieldExpr := range resFieldExprs {
		if resFieldExpr.isComment {
			if prevExpr != nil {
				// TODO
			}
		} else {
			resExprs = append(resExprs, resFieldExpr.expr)
			prevExpr = resFieldExpr.expr
		}
	}
	return nil, errors.New("not implemented")
}

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

type fieldExpr struct {
	index     int
	expr      *dst.KeyValueExpr
	comment   string
	isComment bool
}

type fieldExprs struct {
	m     map[string]*fieldExpr
	count int
}

func (fe *fieldExprs) inc() int {
	r := fe.count
	fe.count++
	return r
}

func (fe fieldExprs) sortedFieldExprList() []*fieldExpr {
	// FIXME
	tmp := make(map[int]*fieldExpr)
	for _, v := range fe.m {
		tmp[v.index] = v
	}
	var res []*fieldExpr
	i := 0
	for len(res) != len(fe.m) {
		if t, ok := tmp[i]; ok {
			res = append(res, t)
		}
		i++
	}
	return res
}

func newFieldExprs(exprs []dst.Expr) *fieldExprs {
	res := &fieldExprs{
		m:     make(map[string]*fieldExpr),
		count: 0,
	}
	for _, expr := range exprs {
		e, ok := expr.(*dst.KeyValueExpr)
		if !ok {
			continue
		}
		if e.Decorations() != nil {
			if e.Decorations().Start != nil {
				res.addFromComments(e.Decorations().Start.All())
			}
			if e.Decorations().End != nil {
				res.addFromComments(e.Decorations().End.All())
			}
		}
		res.addKeyValueExpWithoutComment(e)
	}
	return res
}

func (fe *fieldExprs) addFromComments(comments []string) {
	for _, comment := range comments {
		key := commentToFieldExprKey(comment)
		if key == "" {
			continue
		}
		if _, ok := fe.m[key]; !ok {
			fe.m[key] = &fieldExpr{
				index:     fe.inc(),
				comment:   comment,
				isComment: true,
			}
		}
	}
}
func (fe *fieldExprs) addKeyValueExpWithoutComment(e *dst.KeyValueExpr) {
	if key, ok := e.Key.(*dst.Ident); ok {
		nExpr := dst.Clone(e).(*dst.KeyValueExpr)
		nExpr.Decs = dst.KeyValueExprDecorations{}
		fe.m[key.Name] = &fieldExpr{
			index: fe.inc(),
			expr:  nExpr,
		}
	}
}

func commentToFieldExprKey(comment string) string {
	re := regexp.MustCompile("^\\/\\/\\s(.+):\\s.*$")
	r := re.FindStringSubmatch(comment)
	if len(r) == 2 {
		return r[1]
	}
	return ""
}
