package e2e

import (
	"bytes"
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ab_same_package(t *testing.T) {
	m := mapc.New()
	m.Register(ab.AUser{}, bb.BUser{}, func(option *mapc.Option) {
		option.OutPath = "src/foo/bar.go"
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
	"github.com/abekoh/mapc/e2e/testdata/a"
	"github.com/abekoh/mapc/e2e/testdata/b"
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
