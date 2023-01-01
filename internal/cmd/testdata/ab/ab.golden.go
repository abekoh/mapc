package ab

import (
	"github.com/abekoh/mapc/internal/cmd/testdata/a"
	"github.com/abekoh/mapc/internal/cmd/testdata/b"
)

func ToBUser(inp a.User) b.User {
	return b.User{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
