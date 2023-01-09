package sample

import (
	"time"

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

// ToBUser is mapper from AUser into BUser
func ToBUser(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
