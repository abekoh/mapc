package mapcstd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/abekoh/mapc/testdata/sample"
)

type typedInt int

var IntValue = 1
var StringValue = "foo"
var PointerIntValue = &IntValue

var (
	Int              = NewTypOf(1)
	Int64            = NewTypOf(int64(1))
	String           = NewTypOf("foo")
	TypedInt         = NewTypOf(typedInt(1))
	Object           = NewTypOf(sample.Object{})
	PointerInt       = NewTypOf(&IntValue)
	PointerString    = NewTypOf(&StringValue)
	DoublePointerInt = NewTypOf(&PointerIntValue)
)

func TestAssignMapper_Map(t *testing.T) {
	tests := []struct {
		src    *Typ
		dest   *Typ
		want   Caster
		wantOk bool
	}{
		{src: Int, dest: Int, want: &NopCaster{}, wantOk: true},
		{src: String, dest: String, want: &NopCaster{}, wantOk: true},
		{src: Object, dest: Object, want: &NopCaster{}, wantOk: true},
		{src: Int, dest: Int64, want: nil, wantOk: false},
		{src: Int, dest: TypedInt, want: nil, wantOk: false},
		{src: Int, dest: String, want: nil, wantOk: false},
		{src: String, dest: Int, want: nil, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.src, tt.dest), func(t *testing.T) {
			got, got1 := AssignMapper{}.Map(tt.src, tt.dest)
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
		src    *Typ
		dest   *Typ
		want   Caster
		wantOk bool
	}{
		{
			src:  Int,
			dest: Int,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "int",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:  String,
			dest: String,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "string",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:  Object,
			dest: Object,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "github.com/abekoh/mapc/testdata/sample",
					Name:       "Object",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:  Int,
			dest: Int64,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "int64",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:  Int,
			dest: TypedInt,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "github.com/abekoh/mapc/mapcstd",
					Name:       "typedInt",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:  Int,
			dest: String,
			want: &SimpleCaster{
				caller: &Caller{
					PkgPath:    "",
					Name:       "string",
					CallerType: Type,
				},
			},
			wantOk: true,
		},
		{
			src:    String,
			dest:   Int,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.src, tt.dest), func(t *testing.T) {
			got, got1 := ConvertMapper{}.Map(tt.src, tt.dest)
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
		src    *Typ
		dest   *Typ
		want   Caster
		wantOk bool
	}{
		{
			src:    Int,
			dest:   Int,
			want:   nil,
			wantOk: false,
		},
		{
			src:  Int,
			dest: PointerInt,
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
			src:  PointerInt,
			dest: DoublePointerInt,
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
			src:  String,
			dest: PointerString,
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
			src:    String,
			dest:   PointerInt,
			want:   nil,
			wantOk: false,
		},
		{
			src:    PointerInt,
			dest:   Int,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.src, tt.dest), func(t *testing.T) {
			got, got1 := RefMapper{}.Map(tt.src, tt.dest)
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
		src    *Typ
		dest   *Typ
		want   Caster
		wantOk bool
	}{
		{
			src:    Int,
			dest:   Int,
			want:   nil,
			wantOk: false,
		},
		{
			src:    Int,
			dest:   PointerInt,
			want:   nil,
			wantOk: false,
		},
		{
			src:    String,
			dest:   PointerString,
			want:   nil,
			wantOk: false,
		},
		{
			src:  PointerString,
			dest: String,
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
			src:  PointerInt,
			dest: Int,
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
			src:  DoublePointerInt,
			dest: PointerInt,
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
			src:    PointerInt,
			dest:   String,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.src, tt.dest), func(t *testing.T) {
			got, got1 := DerefMapper{}.Map(tt.src, tt.dest)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("Map() gotOk = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}
