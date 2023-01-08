package mapc_test

import (
	"testing"
	"time"

	"github.com/abekoh/mapc"
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

func TestMapC(t *testing.T) {
	m := mapc.New()
	m.Register(AUser{}, BUser{})
	errs := m.Generate()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}
