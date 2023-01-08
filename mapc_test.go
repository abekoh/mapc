package mapc_test

import (
	"bytes"
	"reflect"
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

type TypedInt int

func TestTypeMapper(t *testing.T) {
	a := 1
	b := TypedInt(2)
	t.Logf("%v", reflect.TypeOf(a).AssignableTo(reflect.TypeOf(b)))
	t.Logf("%v", reflect.TypeOf(a).ConvertibleTo(reflect.TypeOf(b)))
}
