package a

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/b"
)

func RegisterPrivateAUserToBUser(t *testing.T, m mapc.Registerer) {
	t.Helper()
	m.Register(user{}, b.User{})
}

func RegisterBUserToPrivateAUser(t *testing.T, m mapc.Registerer) {
	t.Helper()
	m.Register(b.User{}, user{})
}
