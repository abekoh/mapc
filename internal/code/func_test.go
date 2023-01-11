package code

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/abekoh/mapc/typemapper"
	"github.com/abekoh/mapc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFuncFromMapping(t *testing.T) {
	from, _ := types.NewStruct(reflect.TypeOf(sample.SrcUser{}))
	to, _ := types.NewStruct(reflect.TypeOf(sample.DestUser{}))
	m := &mapping.Mapping{
		From: from,
		To:   to,
		FieldPairs: []*mapping.FieldPair{
			{
				From:   from.Fields[0],
				To:     to.Fields[0],
				Caster: &typemapper.NopCaster{},
			},
		},
	}
	got := NewFuncFromMapping(m, &FuncOption{})
	assert.Equal(t, "ToDestUser", got.name)
	assert.Equal(t, "x", got.argName)
	assert.Equal(t, &Typ{name: "SrcUser", pkgPath: "github.com/abekoh/mapc/testdata/sample"}, got.fromTyp)
	assert.Equal(t, &Typ{name: "DestUser", pkgPath: "github.com/abekoh/mapc/testdata/sample"}, got.toTyp)
	assert.Len(t, got.mapExprs, 1)
	assert.Equal(t, "ID", got.mapExprs[0].From())
	assert.Equal(t, "ID", got.mapExprs[0].To())
}

func Test_funcName(t *testing.T) {
	type args struct {
		m   *mapping.Mapping
		opt *FuncOption
	}
	m := func() *mapping.Mapping {
		from, _ := types.NewStruct(reflect.TypeOf(sample.SrcUser{}))
		to, _ := types.NewStruct(reflect.TypeOf(sample.DestUser{}))
		return &mapping.Mapping{
			From: from,
			To:   to,
		}
	}()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "opt.Name is given",
			args: args{
				opt: &FuncOption{
					Name: "Convert",
				},
			},
			want: "Convert",
		},
		{
			name: "opt.Name is given, private",
			args: args{
				opt: &FuncOption{
					Name:    "Convert",
					Private: true,
				},
			},
			want: "convert",
		},
		{
			name: "opt.NameTemplate is given",
			args: args{
				m: m,
				opt: &FuncOption{
					NameTemplate: func() *template.Template {
						t, _ := template.New("FuncName").Parse("{{ .From }}To{{ .To }}")
						return t
					}(),
				},
			},
			want: "SrcUserToDestUser",
		},
		{
			name: "opt.NameTemplate is given, private",
			args: args{
				m: m,
				opt: &FuncOption{
					NameTemplate: func() *template.Template {
						t, _ := template.New("FuncName").Parse("{{ .From }}To{{ .To }}")
						return t
					}(),
					Private: true,
				},
			},
			want: "srcUserToDestUser",
		},
		{
			name: "option is not set",
			args: args{
				m:   m,
				opt: &FuncOption{},
			},
			want: "ToDestUser",
		},
		{
			name: "only private",
			args: args{
				m: m,
				opt: &FuncOption{
					Private: true,
				},
			},
			want: "toDestUser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, funcName(tt.args.m, tt.args.opt), "funcName(%v, %v)", tt.args.m, tt.args.opt)
		})
	}
}

func Test_argName(t *testing.T) {
	type args struct {
		m   *mapping.Mapping
		opt *FuncOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "opt.ArgName is set",
			args: args{
				opt: &FuncOption{
					ArgName: "input",
				},
			},
			want: "input",
		},
		{
			name: "default",
			args: args{
				opt: &FuncOption{},
			},
			want: "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, argName(tt.args.m, tt.args.opt), "argName(%v, %v)", tt.args.m, tt.args.opt)
		})
	}
}

func TestFunc_AppendNotSetExprs(t *testing.T) {
	genTarget := func() *Func {
		return &Func{
			name:    "AFunc",
			argName: "x",
			fromTyp: &Typ{
				name: "int",
			},
			toTyp: &Typ{
				name: "int",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					from: "z",
					to:   "z",
				},
			},
		}
	}
	t.Run("mapExpr is appended", func(t *testing.T) {
		x := &Func{
			name:    "AFunc",
			argName: "x",
			fromTyp: &Typ{
				name: "int",
			},
			toTyp: &Typ{
				name: "int",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					to: "y",
				},
			},
		}
		target := genTarget()
		assert.Len(t, target.mapExprs, 1)
		err := target.AppendNotSetExprs(x)
		require.Nil(t, err)
		assert.Len(t, target.mapExprs, 2)
	})
	t.Run("same key is existed", func(t *testing.T) {
		x := &Func{
			name:    "AFunc",
			argName: "x",
			fromTyp: &Typ{
				name: "int",
			},
			toTyp: &Typ{
				name: "int",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					to: "z",
				},
			},
		}
		target := genTarget()
		assert.Len(t, target.mapExprs, 1)
		err := target.AppendNotSetExprs(x)
		require.Nil(t, err)
		assert.Len(t, target.mapExprs, 1)
	})
	t.Run("type is not match", func(t *testing.T) {
		x := &Func{
			name:    "AFunc",
			argName: "x",
			fromTyp: &Typ{
				name: "int",
			},
			toTyp: &Typ{
				name: "string",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					to: "y",
				},
			},
		}
		target := genTarget()
		assert.Len(t, target.mapExprs, 1)
		err := target.AppendNotSetExprs(x)
		require.NotNil(t, err)
	})
}
