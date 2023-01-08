package mapping

import (
	"reflect"
	"testing"

	"github.com/abekoh/mapc/internal/object"
)

type TypedInt int

type Hoge struct {
	Fuga string
}

type From struct {
	Int   int
	Int64 int64
	TInt  TypedInt
	Hoge  Hoge
}

type To struct {
	Int   int
	Int64 int64
	TInt  TypedInt
	Hoge  Hoge
}

func loadStruct(t *testing.T, target any) *object.Struct {
	t.Helper()

	s, err := object.NewStruct(reflect.TypeOf(target))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func loadField(t *testing.T, str any, fieldName string) *object.Field {
	t.Helper()

	s := loadStruct(t, str)
	var field *object.Field
	for _, f := range s.Fields {
		if f.Name() == fieldName {
			field = f
			break
		}
	}
	if field == nil {
		t.Fatalf("field '%s' is not found", fieldName)
	}
	return field
}

func Test_newFieldPair(t *testing.T) {
	type args struct {
		from *object.Field
		to   *object.Field
	}
	tests := []struct {
		name       string
		args       args
		wantOk     bool
		wantCaster *Caster
	}{
		{
			name: "int -> int",
			args: args{
				from: loadField(t, From{}, "Int"),
				to:   loadField(t, To{}, "Int"),
			},
			wantOk:     true,
			wantCaster: nil,
		},
		{
			name: "int -> int64",
			args: args{
				from: loadField(t, From{}, "Int"),
				to:   loadField(t, To{}, "Int64"),
			},
			wantOk: true,
			wantCaster: &Caster{
				pkgPath:   "",
				fmtString: "int64(%s)",
			},
		},
		{
			name: "TypedInt -> TypedInt",
			args: args{
				from: loadField(t, From{}, "TInt"),
				to:   loadField(t, To{}, "TInt"),
			},
			wantOk: true,
			wantCaster: &Caster{
				pkgPath:   "github.com/abekoh/mapc/internal/mapping",
				fmtString: "TypedInt(%s)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := newFieldPair(tt.args.from, tt.args.to)
			if gotOk != tt.wantOk {
				t.Errorf("newFieldPair() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(got.Caster, tt.wantCaster) {
				t.Errorf("newFieldPair().Caster got = %v, want %v", got.Caster, tt.wantCaster)
			}
		})
	}
}
