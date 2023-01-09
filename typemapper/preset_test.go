package typemapper

import (
	"reflect"
	"testing"

	"github.com/abekoh/mapc/internal/testdata/sample"
	"github.com/abekoh/mapc/internal/types"
)

type TypedInt int

func TestAssignMapper_Map(t *testing.T) {
	tests := []struct {
		name   string
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{
			name:   "int -> int",
			from:   types.NewTyp(reflect.TypeOf(1)),
			to:     types.NewTyp(reflect.TypeOf(2)),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "string -> string",
			from:   types.NewTyp(reflect.TypeOf("foo")),
			to:     types.NewTyp(reflect.TypeOf("bar")),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "Object -> Object",
			from:   types.NewTyp(reflect.TypeOf(sample.Object{})),
			to:     types.NewTyp(reflect.TypeOf(sample.Object{})),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "int -> int64",
			from:   types.NewTyp(reflect.TypeOf(1)),
			to:     types.NewTyp(reflect.TypeOf(int64(2))),
			want:   nil,
			wantOk: false,
		},
		{
			name:   "int -> TypedInt",
			from:   types.NewTyp(reflect.TypeOf(1)),
			to:     types.NewTyp(reflect.TypeOf(TypedInt(2))),
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AssignMapper{}.Map(tt.from, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
