package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FuncComments(t *testing.T) {
	m := mapc.New()
	m.Option.NoMapperFieldComment = false
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

// MapAUserToBUser maps AUser to BUser.
// This function is generated by mapc.
// You can edit mapping fields.
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

func Test_FieldComments(t *testing.T) {
	t.Run("one field is commented", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = []mapcstd.FieldMapper{
			mapcstd.HashMap{
				"ID":   "ID",
				"Name": "Name",
				"Age":  "Age",
			},
		}
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
		m.Register(ab.AUser{}, ab.BUser{})
		gs, errs := m.Generate()
		requireNoErrors(t, errs)
		got, err := gs[0].Sprint()
		require.Nil(t, err)
		assert.Equal(t, `package ab

func MapAUserToBUser(x AUser) BUser {
	return BUser{
		ID:   x.ID,
		Name: x.Name,
		Age:  x.Age,
		// RegisteredAt:
	}
}
`, got)
	})
	t.Run("two fields are commented", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = []mapcstd.FieldMapper{
			mapcstd.HashMap{
				"ID":   "ID",
				"Name": "Name",
			},
		}
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
		m.Register(ab.AUser{}, ab.BUser{})
		gs, errs := m.Generate()
		requireNoErrors(t, errs)
		got, err := gs[0].Sprint()
		require.Nil(t, err)
		assert.Equal(t, `package ab

func MapAUserToBUser(x AUser) BUser {
	return BUser{
		ID:   x.ID,
		Name: x.Name,
		// Age:
		// RegisteredAt:
	}
}
`, got)
	})
	t.Run("all fields are commented", func(t *testing.T) {
		m := mapc.New()
		m.Option.FuncComment = false
		m.Option.FieldMappers = []mapcstd.FieldMapper{}
		m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
		m.Register(ab.AUser{}, ab.BUser{})
		gs, errs := m.Generate()
		requireNoErrors(t, errs)
		got, err := gs[0].Sprint()
		require.Nil(t, err)
		assert.Equal(t, `package ab

func MapAUserToBUser(x AUser) BUser {
	return BUser{
		// ID:
		// Name:
		// Age:
		// RegisteredAt:
	}
}
`, got)
	})
}
