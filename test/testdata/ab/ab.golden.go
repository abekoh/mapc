package ab

import (
	"github.com/abekoh/mapc/test/testdata/a"
	"github.com/abekoh/mapc/test/testdata/b"
)

func ToBUser(inp a.User) b.User {
	return b.User{
		ID:           inp.ID,
		Name:         inp.Name,
		Age:          inp.Age,
		RegisteredAt: inp.RegisteredAt,
	}
}
