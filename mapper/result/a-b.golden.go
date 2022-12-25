package result

import (
	"github.com/abekoh/mapstructor/mapper/a"
	"github.com/abekoh/mapstructor/mapper/b"
)

func toBUser(aUser a.User) b.User {
	return b.User{
		ID:           aUser.ID,
		Name:         aUser.Name,
		Age:          aUser.Age,
		RegisteredAt: aUser.RegisteredAt,
	}
}

func toAUser(bUser b.User) a.User {
	return a.User{
		ID:           bUser.ID,
		Name:         bUser.Name,
		Age:          bUser.Age,
		RegisteredAt: bUser.RegisteredAt,
	}
}
