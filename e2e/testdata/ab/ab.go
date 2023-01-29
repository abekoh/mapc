package ab

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

type AUser2 struct {
	ID           string
	Name         string
	Age          int
	RegisteredAt string
}

func MapStringToUUID(x string) (uuid.UUID, error) {
	u, err := uuid.Parse(x)
	if err != nil {
		return [16]byte{}, err
	}
	return u, err
}

func MapStringToTime(x string) (time.Time, bool) {
	t, err := time.Parse(x, time.RFC3339)
	if err != nil {
		return t, false
	}
	return t, true
}
