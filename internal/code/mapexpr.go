package code

import (
	"go/token"
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

func (mel MapExprList) DstExprs(arg string) (exprs []dst.Expr, startComments []string) {
	var prevExpr dst.Expr
	firstComment := true
	aboveLineIsExpr := false
	for _, f := range mel {
		if e, ok := f.DstExpr(arg); ok {
			exprs = append(exprs, e)
			prevExpr = e
			aboveLineIsExpr = true
		} else if c, ok := f.Comment(); ok {
			if prevExpr == nil {
				if firstComment {
					startComments = append(startComments, "\n")
					firstComment = false
				}
				startComments = append(startComments, c)
			} else {
				if aboveLineIsExpr {
					prevExpr.Decorations().End.Append("\n")
				}
				prevExpr.Decorations().End.Append(c)
			}
			aboveLineIsExpr = false
		}
	}
	// returned comments are appended to lbrace
	return
}

type SimpleMapExpr struct {
	src     string
	dest    string
	casters []Caster
}

func (e SimpleMapExpr) From() string {
	return e.src
}

func (e SimpleMapExpr) To() string {
	return e.dest
}

func (e SimpleMapExpr) Comment() (string, bool) {
	return "", false
}

func (e SimpleMapExpr) DstExpr(arg string) (dst.Expr, bool) {
	return &dst.KeyValueExpr{
		Key:   dst.NewIdent(e.dest),
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
		Sel:  dst.NewIdent(e.src),
		Decs: dst.SelectorExprDecorations{},
	}
	for _, caster := range e.casters {
		el = caster.Expr(el)
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

type Caster interface {
	Expr(arg dst.Expr) dst.Expr
}

type TypeCaster struct {
	name    string
	pkgPath string
}

func (t TypeCaster) Expr(arg dst.Expr) dst.Expr {
	return &dst.CallExpr{
		Fun:      genTypeWithName(t.name, t.pkgPath),
		Args:     []dst.Expr{arg},
		Ellipsis: false,
		Decs:     dst.CallExprDecorations{},
	}
}

type FuncCallCaster struct {
	name    string
	pkgPath string
}

func (c FuncCallCaster) Expr(arg dst.Expr) dst.Expr {
	return &dst.CallExpr{
		Fun:      genFunc(c.name, c.pkgPath),
		Args:     []dst.Expr{arg},
		Ellipsis: false,
		Decs:     dst.CallExprDecorations{},
	}
}

type UnaryCaster byte

func (u UnaryCaster) Expr(arg dst.Expr) dst.Expr {
	var op token.Token
	switch rune(u) {
	case '&':
		op = token.AND
	case '*':
		op = token.MUL
	default:
		return nil
	}
	return &dst.UnaryExpr{
		Op:   op,
		X:    arg,
		Decs: dst.UnaryExprDecorations{},
	}
}
