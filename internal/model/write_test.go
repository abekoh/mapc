package model

import (
	"testing"
)

func Test_camelize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "convert 1",
			args: args{
				s: "age",
			},
			want: "Age",
		},
		{
			name: "convert 2",
			args: args{
				s: "zod",
			},
			want: "Zod",
		},
		{
			name: "not convert 1",
			args: args{
				s: "AST",
			},
			want: "AST",
		},
		{
			name: "not convert 2",
			args: args{
				s: "Zoo",
			},
			want: "Zoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelize(tt.args.s); got != tt.want {
				t.Errorf("camelize() = %v, want %v", got, tt.want)
			}
		})
	}
}
