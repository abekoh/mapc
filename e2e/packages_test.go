package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/a"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/abekoh/mapc/fieldmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var commonOptFn = func(o *mapc.Option) {
	o.WithoutComment = true
}

func Test_same_package(t *testing.T) {
	m := mapc.New()
	m.Register(ab.AUser{}, ab.BUser{},
		commonOptFn,
		func(option *mapc.Option) {
			option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
		},
	)
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func ToBUser(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, got)
}

func Test_from_package(t *testing.T) {
	m := mapc.New()
	m.Option.WithoutComment = true
	m.Option.FieldMappers = append(m.Option.FieldMappers, &fieldmapper.UpperFirst{})
	a.RegisterPrivateAUserToBUser(t, m)
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package a

import "github.com/abekoh/mapc/e2e/testdata/b"

func ToUser(x User) b.User {
	return b.User{
		ID:           x.id,
		Name:         x.name,
		Age:          x.age,
		RegisteredAt: x.registeredAt,
	}
}
`, got)
}
