package code

import (
	"go/token"
	"regexp"

	"github.com/dave/dst"
)

type MapExpr interface {
	Src() string
	Dest() string
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

func (mel MapExprList) SeparateCommented() (normal, commented MapExprList) {
	for _, me := range mel {
		if _, ok := me.Comment(); ok {
			commented = append(commented, me)
		} else {
			normal = append(normal, me)
		}
	}
	return
}

type SimpleMapExpr struct {
	src     string
	dest    string
	casters []Caster
}

func (e SimpleMapExpr) Src() string {
	return e.src
}

func (e SimpleMapExpr) Dest() string {
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
	dest    string
	comment string
}

func (c CommentedMapExpr) Src() string {
	return ""
}

func (c CommentedMapExpr) Dest() string {
	return c.dest
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
				dest:    key,
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

type UnknownMapExpr struct {
	dest    string
	dstExpr dst.Expr
}

func (e UnknownMapExpr) Src() string {
	return ""
}

func (e UnknownMapExpr) Dest() string {
	return e.dest
}

func (e UnknownMapExpr) DstExpr(_ string) (dst.Expr, bool) {
	return e.dstExpr, true
}

func (e UnknownMapExpr) Comment() (string, bool) {
	return "", false
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
