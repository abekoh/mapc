package function

import (
	"regexp"

	"github.com/dave/dst"
)

type FieldMapper interface {
	From() string
	To() string
	DstExpr(arg string) (dst.Expr, bool)
	Comment() (string, bool)
}

type FieldMapperList []FieldMapper

func (fl FieldMapperList) DstExprs(arg string) (exprs []dst.Expr, comments []string) {
	for _, f := range fl {
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

func (s SimpleFieldMapper) DstExpr(arg string) (dst.Expr, bool) {
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

func (c CommentedFieldMapper) DstExpr(_ string) (dst.Expr, bool) {
	return nil, false
}

func (c CommentedFieldMapper) Comment() (string, bool) {
	return c.comment, true
}

func ParseComments(comments ...string) (res FieldMapperList) {
	for _, c := range comments {
		if key := commentToKey(c); key != "" {
			res = append(res, &CommentedFieldMapper{
				to:      key,
				comment: c,
			})
		}
	}
	return
}

func commentToKey(comment string) string {
	re := regexp.MustCompile("^\\/\\/\\s(.+):\\s.*$")
	r := re.FindStringSubmatch(comment)
	if len(r) == 2 {
		return r[1]
	}
	return ""
}
