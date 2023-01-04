package code

import (
	"strings"

	"github.com/abekoh/mapc/internal/mapping"
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

func NewFunc(m *mapping.Mapping) *Func {
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

func pkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
