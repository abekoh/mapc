package mapper

import (
	"time"

	"github.com/google/uuid"
)

type AUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

type BUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

func ToAUser(inp AUser) BUser {
	return BUser{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
