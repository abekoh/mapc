package mapc_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/abekoh/mapc"
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
	m.Register(AUser{}, BUser{})
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
	b := int64(2)
	c := TypedInt(3)
	d := "4"
	e := mapc.TypedInt2(5)
	rs := []reflect.Type{
		reflect.TypeOf(a),
		reflect.TypeOf(b),
		reflect.TypeOf(c),
		reflect.TypeOf(d),
		reflect.TypeOf(e),
	}
	for _, x := range rs {
		for _, y := range rs {
			//t.Logf("assign %v -> %v: %v", x, y, x.AssignableTo(y))
			t.Logf("convert %v -> %v: %v", x, y, x.ConvertibleTo(y))
		}
	}
}
