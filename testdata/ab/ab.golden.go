package ab

import (
	"github.com/abekoh/mapc/testdata/a"
	"github.com/abekoh/mapc/testdata/b"
)

func ToBUser(inp a.User) b.User {
	return b.User{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
