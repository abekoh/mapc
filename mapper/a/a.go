package a

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}
