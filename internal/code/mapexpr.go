package code

import (
	"regexp"

	"github.com/dave/dst"
)

type MapExpr interface {
	From() string
	To() string
	DstExpr(arg string) (dst.Expr, bool)
	Comment() (string, bool)
}

type MapExprList []MapExpr

func (mel MapExprList) DstExprs(arg string) (exprs []dst.Expr, comments []string) {
	for _, f := range mel {
		if e, ok := f.DstExpr(arg); ok {
			e.Decorations().Start.Append(comments...)
			comments = []string{}
			exprs = append(exprs, e)
		} else if c, ok := f.Comment(); ok {
			comments = append(comments, c)
		}
	}
	return
}

type SimpleMapExpr struct {
	from    string
	to      string
	casters []Caster
}

type Caster interface {
	Ident() *dst.Ident
}

type TypeCaster struct {
	typ *Typ
}

func (t TypeCaster) Ident() *dst.Ident {
	return genType(t.typ)
}

type CallerCaster struct {
	caller *Caller
}

func (c CallerCaster) Ident() *dst.Ident {
	return genCaller(c.caller)
}

func (e SimpleMapExpr) From() string {
	return e.from
}

func (e SimpleMapExpr) To() string {
	return e.to
}

func (e SimpleMapExpr) Comment() (string, bool) {
	return "", false
}

func (e SimpleMapExpr) DstExpr(arg string) (dst.Expr, bool) {
	return &dst.KeyValueExpr{
		Key:   dst.NewIdent(e.to),
		Value: e.valueExpr(arg),
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.NewLine,
				After:  dst.NewLine,
			},
		},
	}, true
}

func (e SimpleMapExpr) valueExpr(arg string) dst.Expr {
	var el dst.Expr
	el = &dst.SelectorExpr{
		X:    genVar(arg),
		Sel:  dst.NewIdent(e.from),
		Decs: dst.SelectorExprDecorations{},
	}
	for _, caster := range e.casters {
		newEl := &dst.CallExpr{
			Fun:      caster.Ident(),
			Args:     []dst.Expr{el},
			Ellipsis: false,
			Decs:     dst.CallExprDecorations{},
		}
		el = newEl
	}
	return el
}

type CommentedMapExpr struct {
	to      string
	comment string
}

func (c CommentedMapExpr) From() string {
	return ""
}

func (c CommentedMapExpr) To() string {
	return c.to
}

func (c CommentedMapExpr) DstExpr(_ string) (dst.Expr, bool) {
	return nil, false
}

func (c CommentedMapExpr) Comment() (string, bool) {
	return c.comment, true
}

func parseComments(comments ...string) (res MapExprList) {
	for _, c := range comments {
		if key := commentToKey(c); key != "" {
			res = append(res, &CommentedMapExpr{
				to:      key,
				comment: c,
			})
		}
	}
	return
}

func commentToKey(comment string) string {
	re := regexp.MustCompile(`^\\/\\/\\s(.+):\\s.*$`)
	r := re.FindStringSubmatch(comment)
	if len(r) == 2 {
		return r[1]
	}
	return ""
}
