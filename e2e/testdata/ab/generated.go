package ab

import "time"

// ModeTestMapper maps AUser to BUser.
// This function is generated by mapc.
// You can edit mapping fields.
func ModeTestMapper(x AUser) BUser {
	return BUser{
		Name:         "Example",
		RegisteredAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
