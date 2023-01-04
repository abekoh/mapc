package code

import (
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
	return &Func{
		fc: &dst.FuncDecl{
			Recv: nil,
			Name: genFuncName(m),
			Type: &dst.FuncType{
				Func:       false,
				TypeParams: nil,
				Params:     genParams(m.From),
				Results:    genResult(m.To),
				Decs:       dst.FuncTypeDecorations{},
			},
			Body: &dst.BlockStmt{
				List: []dst.Stmt{
					genReturn(m),
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

func genParams(fromStr *object.Struct) *dst.FieldList {
	fromObj := dst.NewObj(dst.Var, "from")
	from := dst.NewIdent("from")
	from.Obj = fromObj
	return &dst.FieldList{
		Opening: true,
		List: []*dst.Field{
			{
				Names: []*dst.Ident{from},
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

func genReturn(m *mapping.Mapping) *dst.ReturnStmt {
	var elts []dst.Expr
	fromTyp := genType(m.From)
	for _, fp := range m.FieldPairs {
		elts = append(elts, genElt(&fp, cloneIdent(fromTyp)))
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
		Decs: dst.KeyValueExprDecorations{},
	}
}

func (f Func) ReNew(m *mapping.Mapping) (*Func, error) {
	return nil, nil
}

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
