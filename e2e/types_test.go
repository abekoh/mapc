package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/various"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AssignableType(t *testing.T) {
	m := mapc.New()
	m.Option.WithoutComment = true
	m.Register(various.S{}, various.S2{},
		func(option *mapc.Option) {
			option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/various"
		},
	)
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package various

func ToS2(x S) S2 {
	return S2{
		Bool:            x.Bool,
		Int:             x.Int,
		Int8:            x.Int8,
		Int16:           x.Int16,
		Int32:           x.Int32,
		Int64:           x.Int64,
		Uint:            x.Uint,
		Uint8:           x.Uint8,
		Uint16:          x.Uint16,
		Uint32:          x.Uint32,
		Uint64:          x.Uint64,
		Uintptr:         x.Uintptr,
		Float32:         x.Float32,
		Float64:         x.Float64,
		Complex64:       x.Complex64,
		Complex128:      x.Complex128,
		IntArray:        x.IntArray,
		IntChan:         x.IntChan,
		IntToStringFunc: x.IntToStringFunc,
		Interface:       x.Interface,
		StringIntMap:    x.StringIntMap,
		Slice:           x.Slice,
		String:          x.String,
		EmptyStruct:     x.EmptyStruct,
		ExternalType:    x.ExternalType,
	}
}
`, got)
}
