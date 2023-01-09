package typemapper

import (
	"reflect"
	"testing"

	"github.com/abekoh/mapc/internal/object"
	"github.com/abekoh/mapc/internal/testdata/sample"
)

type TypedInt int

func TestAssignMapper_Map(t *testing.T) {
	tests := []struct {
		name   string
		from   *object.Typ
		to     *object.Typ
		want   Caster
		wantOk bool
	}{
		{
			name:   "int -> int",
			from:   object.NewTyp(reflect.TypeOf(1)),
			to:     object.NewTyp(reflect.TypeOf(2)),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "string -> string",
			from:   object.NewTyp(reflect.TypeOf("foo")),
			to:     object.NewTyp(reflect.TypeOf("bar")),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "Object -> Object",
			from:   object.NewTyp(reflect.TypeOf(sample.Object{})),
			to:     object.NewTyp(reflect.TypeOf(sample.Object{})),
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			name:   "int -> int64",
			from:   object.NewTyp(reflect.TypeOf(1)),
			to:     object.NewTyp(reflect.TypeOf(int64(2))),
			want:   nil,
			wantOk: false,
		},
		{
			name:   "int -> TypedInt",
			from:   object.NewTyp(reflect.TypeOf(1)),
			to:     object.NewTyp(reflect.TypeOf(TypedInt(2))),
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
