package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/a"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WithSamePackage(t *testing.T) {
	m := mapc.New()
	m.Option.FuncComment = false
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func MapAUserToBUser(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, got)
}

func Test_SrcIsPrivate(t *testing.T) {
	t.Run("outPkgPath is src's, success all fields", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = append(m.Option.FieldMappers,
			&mapcstd.UpperFirst{},
			&mapcstd.HashMap{"id": "ID"},
		)
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/a"
		a.RegisterPrivateAUserToBUser(t, m)
		gs, errs := m.Generate()
		requireNoErrors(t, errs)
		got, err := gs[0].Sprint()
		require.Nil(t, err)
		assert.Equal(t, `package a

import "github.com/abekoh/mapc/e2e/testdata/b"

func MapUserToUser(x user) b.User {
	return b.User{
		ID:           x.id,
		Name:         x.name,
		Age:          x.age,
		RegisteredAt: x.registeredAt,
	}
}
`, got)
	})
	t.Run("outPkgPath is dest's, fail", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = append(m.Option.FieldMappers,
			&mapcstd.UpperFirst{},
			&mapcstd.HashMap{"id": "ID"},
		)
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/b"
		a.RegisterPrivateAUserToBUser(t, m)
		_, errs := m.Generate()
		require.Len(t, errs, 1)
	})
}

func Test_DestIsPrivate(t *testing.T) {
	t.Run("outPkgPath is dest's, success all fields", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = append(m.Option.FieldMappers,
			&mapcstd.LowerFirst{},
			&mapcstd.HashMap{"ID": "id"},
		)
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/a"
		a.RegisterBUserToPrivateAUser(t, m)
		gs, errs := m.Generate()
		requireNoErrors(t, errs)
		got, err := gs[0].Sprint()
		require.Nil(t, err)
		assert.Equal(t, `package a

import "github.com/abekoh/mapc/e2e/testdata/b"

func MapUserToUser(x b.User) user {
	return user{
		id:           x.ID,
		name:         x.Name,
		age:          x.Age,
		registeredAt: x.RegisteredAt,
	}
}
`, got)
	})
	t.Run("outPkgPath is dest's, fail", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = append(m.Option.FieldMappers,
			&mapcstd.LowerFirst{},
			&mapcstd.HashMap{"ID": "id"},
		)
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/b"
		a.RegisterPrivateAUserToBUser(t, m)
		_, errs := m.Generate()
		require.Len(t, errs, 1)
	})
}
