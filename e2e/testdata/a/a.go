package a

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

type user struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}
