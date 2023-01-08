package mapc_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/pkg/fieldmapper"
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
		OutPath: "src/foo/bar.go",
		FieldMappers: []fieldmapper.FieldMapper{
			&fieldmapper.UpperFirst{},
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
		assert.Equal(t, `package foo

import "github.com/abekoh/mapc_test"

func ToBUser(x mapc_test.AUser) mapc_test.BUser {
	return mapc_test.BUser{
		Name:         x.name,
		Age:          x.age,
		RegisteredAt: x.registeredAt,
	}
}
`, buf.String())
	}
}

func TestPath(t *testing.T) {
	f, _ := os.Getwd()
	t.Log(f)
	t.Log(filepath.Base(f))
	t.Log(filepath.Dir(f))
}
