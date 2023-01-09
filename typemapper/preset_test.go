package typemapper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/abekoh/mapc/testdata/sample"
	"github.com/abekoh/mapc/types"
)

type typedInt int

var (
	Int      = types.NewTyp(reflect.TypeOf(1))
	Int64    = types.NewTyp(reflect.TypeOf(int64(1)))
	String   = types.NewTyp(reflect.TypeOf("foo"))
	TypedInt = types.NewTyp(reflect.TypeOf(typedInt(1)))
	Object   = types.NewTyp(reflect.TypeOf(sample.Object{}))
)

func TestAssignMapper_Map(t *testing.T) {
	tests := []struct {
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{
			from:   Int,
			to:     Int,
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			from:   String,
			to:     String,
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			from:   Object,
			to:     Object,
			want:   &NopCaster{},
			wantOk: true,
		},
		{
			from:   Int,
			to:     Int64,
			want:   nil,
			wantOk: false,
		},
		{
			from:   Int,
			to:     TypedInt,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.to), func(t *testing.T) {
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
