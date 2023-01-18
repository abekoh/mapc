package typemapper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/abekoh/mapc/testdata/sample"
	"github.com/abekoh/mapc/types"
)

type typedInt int

var IntValue = 1
var StringValue = "foo"
var PointerIntValue = &IntValue

var (
	Int              = types.NewTyp(reflect.TypeOf(1))
	Int64            = types.NewTyp(reflect.TypeOf(int64(1)))
	String           = types.NewTyp(reflect.TypeOf("foo"))
	TypedInt         = types.NewTyp(reflect.TypeOf(typedInt(1)))
	Object           = types.NewTyp(reflect.TypeOf(sample.Object{}))
	PointerInt       = types.NewTyp(reflect.TypeOf(&IntValue))
	PointerString    = types.NewTyp(reflect.TypeOf(&StringValue))
	DoublePointerInt = types.NewTyp(reflect.TypeOf(&PointerIntValue))
)

func TestAssignMapper_Map(t *testing.T) {
	tests := []struct {
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{from: Int, to: Int, want: &NopCaster{}, wantOk: true},
		{from: String, to: String, want: &NopCaster{}, wantOk: true},
		{from: Object, to: Object, want: &NopCaster{}, wantOk: true},
		{from: Int, to: Int64, want: nil, wantOk: false},
		{from: Int, to: TypedInt, want: nil, wantOk: false},
		{from: Int, to: String, want: nil, wantOk: false},
		{from: String, to: Int, want: nil, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.to), func(t *testing.T) {
			got, got1 := AssignMapper{}.Map(tt.from, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}

func TestConvertMapper_Map(t *testing.T) {
	tests := []struct {
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{
			from: Int,
			to:   Int,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "int",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from: String,
			to:   String,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "string",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from: Object,
			to:   Object,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "github.com/abekoh/mapc/testdata/sample",
					Name:       "Object",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from: Int,
			to:   Int64,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "int64",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from: Int,
			to:   TypedInt,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "github.com/abekoh/mapc/typemapper",
					Name:       "typedInt",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from: Int,
			to:   String,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "string",
					CallerType: Typ,
				},
			},
			wantOk: true,
		},
		{
			from:   String,
			to:     Int,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.to), func(t *testing.T) {
			got, got1 := ConvertMapper{}.Map(tt.from, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}

func TestRefMapper_Map(t *testing.T) {
	tests := []struct {
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{
			from:   Int,
			to:     Int,
			want:   nil,
			wantOk: false,
		},
		{
			from: Int,
			to:   PointerInt,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "&",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from: PointerInt,
			to:   DoublePointerInt,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "&",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from: String,
			to:   PointerString,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "&",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from:   String,
			to:     PointerInt,
			want:   nil,
			wantOk: false,
		},
		{
			from:   PointerInt,
			to:     Int,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.to), func(t *testing.T) {
			got, got1 := RefMapper{}.Map(tt.from, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}

func TestDerefMapper_Map(t *testing.T) {
	tests := []struct {
		from   *types.Typ
		to     *types.Typ
		want   Caster
		wantOk bool
	}{
		{
			from:   Int,
			to:     Int,
			want:   nil,
			wantOk: false,
		},
		{
			from:   Int,
			to:     PointerInt,
			want:   nil,
			wantOk: false,
		},
		{
			from:   String,
			to:     PointerString,
			want:   nil,
			wantOk: false,
		},
		{
			from: PointerString,
			to:   String,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "*",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from: PointerInt,
			to:   Int,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "*",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from: DoublePointerInt,
			to:   PointerInt,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "*",
					CallerType: Unary,
				},
			},
			wantOk: true,
		},
		{
			from:   PointerInt,
			to:     String,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.from, tt.to), func(t *testing.T) {
			got, got1 := DerefMapper{}.Map(tt.from, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}
