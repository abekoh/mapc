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
	from string
	to   string
}

func (s SimpleMapExpr) From() string {
	return s.from
}

func (s SimpleMapExpr) To() string {
	return s.to
}

func (s SimpleMapExpr) DstExpr(arg string) (dst.Expr, bool) {
	return &dst.KeyValueExpr{
		Key: dst.NewIdent(s.to),
		Value: &dst.SelectorExpr{
			X:    genVar(arg),
			Sel:  dst.NewIdent(s.from),
			Decs: dst.SelectorExprDecorations{},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.NewLine,
				After:  dst.NewLine,
			},
		},
	}, true
}

func (s SimpleMapExpr) Comment() (string, bool) {
	return "", false
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
