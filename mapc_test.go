package mapc_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/internal/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AUser is source for mapping test
type AUser struct {
	id           uuid.UUID
	name         string
	age          int
	registeredAt time.Time
}

// BUser is destination for mapping test
type BUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

func TestMapC(t *testing.T) {
	m := mapc.New()
	m.Register(AUser{}, BUser{}, &mapc.Option{
		FieldMappers: []mapc.FieldMapper{
			func(s string) string {
				return util.UpperFirst(s)
			},
		},
	})
	gs, errs := m.Generate()
	for _, err := range errs {
		require.Nil(t, err)
	}
	for _, g := range gs {
		var buf bytes.Buffer
		err := g.Write(&buf)
		require.Nil(t, err)
		assert.Equal(t, `
`, buf.String())
	}
}

type TypeMapperFunc[T, U any] func(T) U

func (tm TypeMapperFunc[T, U]) Map(inp any) any {
	return tm(inp)
}

type TypeMapper interface {
	Map(any) any
}

func TestTypeMapper(t *testing.T) {
	mappers := []TypeMapper{
		TypeMapperFunc[string, string](func(i string) string {
			return i
		}),
		TypeMapperFunc[int, int64](func(i int) int64 {
			return int64(i)
		}),
	}
	_ = mappers
}
