package function

import (
	"github.com/dave/dst"
)

type FieldMapper interface {
	From() string
	To() string
	DstExpr(argIdent *dst.Ident) (dst.Expr, bool)
	Comment() (string, bool)
}

type FieldMapperList []FieldMapper

func (fl FieldMapperList) DstExprs(argIdent *dst.Ident) (exprs []dst.Expr, comments []string) {
	for _, f := range fl {
		if e, ok := f.DstExpr(cloneIdent(argIdent)); ok {
			e.Decorations().Start.Append(comments...)
			comments = []string{}
			exprs = append(exprs, e)
		} else if c, ok := f.Comment(); ok {
			comments = append(comments, c)
		}
	}
	return
}

type SimpleFieldMapper struct {
	from string
	to   string
}

func (s SimpleFieldMapper) From() string {
	return s.from
}

func (s SimpleFieldMapper) To() string {
	return s.to
}

func (s SimpleFieldMapper) DstExpr(argIdent *dst.Ident) (dst.Expr, bool) {
	return &dst.KeyValueExpr{
		Key: dst.NewIdent(s.to),
		Value: &dst.SelectorExpr{
			X:    argIdent,
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

func (s SimpleFieldMapper) Comment() (string, bool) {
	return "", false
}

type CommentedFieldMapper struct {
	to      string
	comment string
}

func (c CommentedFieldMapper) From() string {
	return ""
}

func (c CommentedFieldMapper) To() string {
	return c.to
}

func (c CommentedFieldMapper) DstExpr(ident *dst.Ident) (dst.Expr, bool) {
	return nil, false
}

func (c CommentedFieldMapper) Comment() (string, bool) {
	return c.comment, true
}
