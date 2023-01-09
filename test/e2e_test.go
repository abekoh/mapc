package main

import (
	"bytes"
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/test/testdata/a"
	"github.com/abekoh/mapc/test/testdata/b"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//// AUser is source for mapping test
//type AUser struct {
//	id           uuid.UUID
//	name         string
//	age          int
//	registeredAt time.Time
//}
//
//// BUser is destination for mapping test
//type BUser struct {
//	ID           uuid.UUID
//	Name         string
//	Age          int
//	RegisteredAt time.Time
//}

func TestMapC(t *testing.T) {
	m := mapc.New()
	m.Register(a.User{}, b.User{}, &mapc.Option{
		OutPath: "src/foo/bar.go",
	})
	gs, errs := m.Generate()
	for _, err := range errs {
		require.Nil(t, err)
	}
	for _, g := range gs {
		var buf bytes.Buffer
		err := g.Write(&buf)
		require.Nil(t, err)
		assert.Equal(t, `package foo

import (
	"github.com/abekoh/mapc/test/testdata/a"
	"github.com/abekoh/mapc/test/testdata/b"
)

func ToUser(x a.User) b.User {
	return b.User{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, buf.String())
	}
}
