package mapper

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
func ToBUser(inp AUser) BUser {
	return BUser{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
