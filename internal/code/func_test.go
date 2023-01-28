package code

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/str"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFuncFromMapping(t *testing.T) {
	src, _ := str.NewStruct(reflect.TypeOf(sample.SrcUser{}))
	dest, _ := str.NewStruct(reflect.TypeOf(sample.DestUser{}))
	m := &mapping.Mapping{
		Src:  src,
		Dest: dest,
		FieldPairs: []*mapping.FieldPair{
			{
				Src:     src.Fields[0],
				Dest:    dest.Fields[0],
				Casters: []mapcstd.Caster{&mapcstd.NopCaster{}},
			},
		},
	}
	got := NewFuncFromMapping(m, &FuncOption{})
	assert.Equal(t, "MapSrcUserToDestUser", got.name)
	assert.Equal(t, "x", got.argName)
	assert.Equal(t, &Typ{name: "SrcUser", pkgPath: "github.com/abekoh/mapc/testdata/sample"}, got.srcTyp)
	assert.Equal(t, &Typ{name: "DestUser", pkgPath: "github.com/abekoh/mapc/testdata/sample"}, got.destTyp)
	assert.Len(t, got.mapExprs, 1)
	assert.Equal(t, "ID", got.mapExprs[0].Src())
	assert.Equal(t, "ID", got.mapExprs[0].Dest())
}

func Test_funcName(t *testing.T) {
	type args struct {
		m   *mapping.Mapping
		opt *FuncOption
	}
	m := func() *mapping.Mapping {
		src, _ := str.NewStruct(reflect.TypeOf(sample.SrcUser{}))
		dest, _ := str.NewStruct(reflect.TypeOf(sample.DestUser{}))
		return &mapping.Mapping{
			Src:  src,
			Dest: dest,
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
						t, _ := template.New("FuncName").Parse("{{ .Src }}To{{ .Dest }}")
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
						t, _ := template.New("FuncName").Parse("{{ .Src }}To{{ .Dest }}")
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
			want: "MapSrcUserToDestUser",
		},
		{
			name: "only private",
			args: args{
				m: m,
				opt: &FuncOption{
					Private: true,
				},
			},
			want: "mapSrcUserToDestUser",
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

func TestFunc_FillMapExprs(t *testing.T) {
	genTarget := func() *Func {
		return &Func{
			name:    "MapAToB",
			argName: "x",
			srcTyp: &Typ{
				name: "A",
			},
			destTyp: &Typ{
				name: "B",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					src:  "s",
					dest: "s",
				},
				&CommentedMapExpr{
					dest:    "t",
					comment: "t",
				},
			},
		}
	}
	t.Run("SimpleMapExpr is appended", func(t *testing.T) {
		x := &Func{
			name:    "MapAToB",
			argName: "x",
			srcTyp: &Typ{
				name: "A",
			},
			destTyp: &Typ{
				name: "B",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					src:  "u",
					dest: "u",
				},
			},
		}
		got, err := genTarget().FillMapExprs(x)
		require.Nil(t, err)
		assert.Equal(t, MapExprList{
			&SimpleMapExpr{
				src:  "s",
				dest: "s",
			},
			&SimpleMapExpr{
				src:  "u",
				dest: "u",
			},
			&CommentedMapExpr{
				dest:    "t",
				comment: "t",
			},
		}, got.mapExprs)
	})
	t.Run("CommentedMepExpr is appended", func(t *testing.T) {
		x := &Func{
			name:    "MapAToB",
			argName: "x",
			srcTyp: &Typ{
				name: "A",
			},
			destTyp: &Typ{
				name: "B",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					dest:    "v",
					comment: "v",
				},
			},
		}
		got, err := genTarget().FillMapExprs(x)
		require.Nil(t, err)
		assert.Equal(t, MapExprList{
			&SimpleMapExpr{
				src:  "s",
				dest: "s",
			},
			&CommentedMapExpr{
				dest:    "t",
				comment: "t",
			},
			&CommentedMapExpr{
				dest:    "v",
				comment: "v",
			},
		}, got.mapExprs)
	})
	t.Run("same key is existed", func(t *testing.T) {
		x := &Func{
			name:    "MapAToB",
			argName: "x",
			srcTyp: &Typ{
				name: "A",
			},
			destTyp: &Typ{
				name: "B",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					src:  "s",
					dest: "s",
				},
			},
		}
		got, err := genTarget().FillMapExprs(x)
		require.Nil(t, err)
		assert.Equal(t, MapExprList{
			&SimpleMapExpr{
				src:  "s",
				dest: "s",
			},
			&CommentedMapExpr{
				dest:    "t",
				comment: "t",
			},
		}, got.mapExprs)
	})
	t.Run("type is not match", func(t *testing.T) {
		x := &Func{
			name:    "MapAToB",
			argName: "x",
			srcTyp: &Typ{
				name: "A",
			},
			destTyp: &Typ{
				name: "C",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					dest: "y",
				},
			},
		}
		_, err := genTarget().FillMapExprs(x)
		require.NotNil(t, err)
	})
}
