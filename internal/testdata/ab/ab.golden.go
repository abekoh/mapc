package ab

import (
	"github.com/abekoh/mapstructor/internal/testdata/a"
	"github.com/abekoh/mapstructor/internal/testdata/b"
)

func toBUser(aUser a.User) b.User {
	return b.User{
		ID:           aUser.ID,
		Name:         aUser.Name,
		Age:          aUser.Age,
		RegisteredAt: aUser.RegisteredAt,
	}
}
