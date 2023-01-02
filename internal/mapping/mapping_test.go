package mapping

import (
	"reflect"
	"testing"
)

type Hoge struct {
	Fuga string
}

type From struct {
	Int   int
	Int64 int64
	Hoge  Hoge
}

type To struct {
	Int   int
	Int64 int64
	Hoge  Hoge
}

func loadField(t *testing.T, str any, fieldName string) Var {
	t.Helper()

	s := loadStruct(t, str)
	f := s.Var(fieldName)
	if f == nil {
		t.Fatalf("field '%s' is not found", fieldName)
	}
	return Var{
		v: f,
	}
}

func Test_match(t *testing.T) {
	type args struct {
		from            *Struct
		toTokenFieldMap tokenFieldMap
	}
	tests := []struct {
		name         string
		args         args
		wantElements int
	}{
		{
			name: "int -> int",
			args: args{
				from:            loadStruct(t, From{}),
				toTokenFieldMap: tokenFieldMap{"Int": loadField(t, To{}, "Int")},
			},
			wantElements: 1,
		},
		{
			name: "int64 -> int",
			args: args{
				from:            loadStruct(t, From{}),
				toTokenFieldMap: tokenFieldMap{"Int64": loadField(t, To{}, "Int")},
			},
			wantElements: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := match(tt.args.from, tt.args.toTokenFieldMap); len(got) != tt.wantElements {
				t.Errorf("match() = %v, wantElements = %v", got, tt.wantElements)
			}
		})
	}
}

func Test_newFieldPair(t *testing.T) {
	type args struct {
		from Var
		to   Var
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantCaster *Caster
	}{
		{
			name: "int -> int",
			args: args{
				from: loadField(t, From{}, "Int"),
				to:   loadField(t, To{}, "Int"),
			},
			wantExists: true,
			wantCaster: nil,
		},
		{
			name: "int -> int64",
			args: args{
				from: loadField(t, From{}, "Int"),
				to:   loadField(t, To{}, "Int64"),
			},
			wantExists: true,
			wantCaster: &Caster{
				fmtString: "int64(%s)",
			},
		},
		{
			name: "int64 -> int",
			args: args{
				from: loadField(t, From{}, "Int64"),
				to:   loadField(t, To{}, "Int"),
			},
			wantExists: true,
			wantCaster: &Caster{
				fmtString: "int(%s)",
			},
		},
		{
			name: "struct -> struct",
			args: args{
				from: loadField(t, From{}, "Hoge"),
				to:   loadField(t, To{}, "Hoge"),
			},
			wantExists: true,
			wantCaster: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newFieldPair(tt.args.from, tt.args.to)
			if !reflect.DeepEqual(got.Caster, tt.wantCaster) {
				t.Errorf("newFieldPair().Caster. got = %v, want %v", got.Caster, tt.wantCaster)
			}
			if got1 != tt.wantExists {
				t.Errorf("newFieldPair() got1 = %v, want %v", got1, tt.wantExists)
			}
		})
	}
}
