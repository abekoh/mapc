package code

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/abekoh/mapc/types"
	"github.com/stretchr/testify/assert"
)

func Test_funcName(t *testing.T) {
	type args struct {
		m   *mapping.Mapping
		opt *FuncOption
	}
	m := func() *mapping.Mapping {
		from, _ := types.NewStruct(reflect.TypeOf(sample.AUser{}))
		to, _ := types.NewStruct(reflect.TypeOf(sample.BUser{}))
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
			want: "AUserToBUser",
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
			want: "aUserToBUser",
		},
		{
			name: "option is not set",
			args: args{
				m:   m,
				opt: &FuncOption{},
			},
			want: "ToBUser",
		},
		{
			name: "only private",
			args: args{
				m: m,
				opt: &FuncOption{
					Private: true,
				},
			},
			want: "toBUser",
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
