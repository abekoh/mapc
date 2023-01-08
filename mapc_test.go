package mapc_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/internal/util"
	"github.com/google/uuid"
)

// AUser is source for mapping test
type AUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// BUser is destination for mapping test
type BUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// CUser is destination for mapping test
type CUser struct {
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
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, g := range gs {
		var buf bytes.Buffer
		err := g.Write(&buf)
		if err != nil {
			t.Fatal(err)
		}
	}
}
