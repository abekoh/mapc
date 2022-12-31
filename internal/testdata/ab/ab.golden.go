package ab

import (
	"github.com/abekoh/mapstructor/internal/testdata/a"
	"github.com/abekoh/mapstructor/internal/testdata/b"
)

func ToBUser(inp a.User) b.User {
	return b.User{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
