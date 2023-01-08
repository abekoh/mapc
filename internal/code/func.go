package code

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/abekoh/mapc"
	//"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/internal/util"
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

func NewFuncSand(m *mapc.Mapping) *Func {
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

func NewFunc(m *mapc.Mapping) *Func {
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

func genFuncName(m *mapc.Mapping) *dst.Ident {
	name := fmt.Sprintf("To%s", util.UpperFirst(m.To.Name))
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

func genReturn(m *mapc.Mapping, argIdent *dst.Ident) *dst.ReturnStmt {
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

func genElt(fp *mapc.FieldPair, from *dst.Ident) *dst.KeyValueExpr {
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

func (f Func) ReNew(m *mapc.Mapping) (*Func, error) {
	//decl := dst.Clone(f.fc).(*dst.FuncDecl)
	//if decl.Body == nil || decl.Body.List == nil || len(decl.Body.List) == 0 {
	//	return nil, errors.New("failed to find return statement")
	//}
	//retStmt, ok := decl.Body.List[len(decl.Body.List)-1].(*dst.ReturnStmt)
	//if !ok {
	//	return nil, errors.New("failed cast to *dst.ReturnStmt")
	//}
	//if retStmt.Results == nil || len(retStmt.Results) == 0 {
	//	return nil, errors.New("failed to get return results")
	//}
	//compLit, ok := retStmt.Results[0].(*dst.CompositeLit)
	//if !ok {
	//	return nil, errors.New("failed to cast *dst.CompositeLit")
	//}
	//eFieldExprs := newStructFieldExprs(compLit.Elts)
	//var resFieldExprList []*structFieldExpr
	//for _, pair := range m.FieldPairs {
	//	// TODO: get 'from ident' from params
	//	fromIdent := genVar("from")
	//	resFieldExprList = append(resFieldExprList, &structFieldExpr{
	//		expr: genElt(pair, fromIdent),
	//	})
	//	key := pair.To.Name()
	//	if _, ok := eFieldExprs.m[key]; ok {
	//		eFieldExprs.m[key].used = true
	//	}
	//}
	//for _, e := range eFieldExprs.sortedFieldExprList() {
	//	resFieldExprs = append(resFieldExprs, e)
	//}
	//var resExprs []dst.Expr
	//var prevExpr *dst.KeyValueExpr
	//var nextComment string
	//for _, resFieldExpr := range resFieldExprs {
	//	if resFieldExpr.isComment {
	//		if prevExpr != nil {
	//			// TODO
	//		}
	//	} else {
	//		resExprs = append(resExprs, resFieldExpr.expr)
	//		prevExpr = resFieldExpr.expr
	//	}
	//}
	return nil, errors.New("not implemented")
}

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

type structFieldExpr struct {
	expr      *dst.KeyValueExpr
	comment   string
	isComment bool
	used      bool
	next      *structFieldExpr
}

type structFieldExprs struct {
	m        map[string]*structFieldExpr
	ordering []*structFieldExpr
}

func newStructFieldExprs(exprs []dst.Expr) *structFieldExprs {
	res := &structFieldExprs{
		m:        make(map[string]*structFieldExpr),
		ordering: make([]*structFieldExpr, 0),
	}
	for _, expr := range exprs {
		keyValueExpr, ok := expr.(*dst.KeyValueExpr)
		if !ok {
			continue
		}
		if keyValueExpr.Decorations() != nil {
			if keyValueExpr.Decorations().Start != nil {
				res.addComments(keyValueExpr.Decorations().Start.All())
			}
			if keyValueExpr.Decorations().End != nil {
				res.addComments(keyValueExpr.Decorations().End.All())
			}
		}
		res.addKeyValueExp(keyValueExpr)
	}
	return res
}

func (ste *structFieldExprs) last() *structFieldExpr {
	if len(ste.ordering) == 0 {
		return nil
	}
	return ste.ordering[len(ste.ordering)-1]
}

func (ste *structFieldExprs) addComments(comments []string) {
	for _, comment := range comments {
		k := commentToFieldExprKey(comment)
		if k == "" {
			continue
		}
		if _, ok := ste.m[k]; !ok {
			res := &structFieldExpr{
				comment:   comment,
				isComment: true,
			}
			ste.m[k] = res
			if last := ste.last(); last != nil {
				last.next = res
			}
			ste.ordering = append(ste.ordering, res)
		}
	}
}

func (ste *structFieldExprs) addKeyValueExp(e *dst.KeyValueExpr) {
	if key, ok := e.Key.(*dst.Ident); ok {
		clonedE := dst.Clone(e).(*dst.KeyValueExpr)
		clonedE.Decs = dst.KeyValueExprDecorations{}
		res := &structFieldExpr{
			expr: clonedE,
		}
		ste.m[key.Name] = res
		if last := ste.last(); last != nil {
			last.next = res
		}
		ste.ordering = append(ste.ordering, res)
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
