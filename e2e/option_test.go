package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OptionWithGroup(t *testing.T) {
	m := mapc.New()
	m.Option.FuncComment = false
	m.Option.NoMapperFieldComment = false
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Option.FieldMappers = []mapcstd.FieldMapper{
		mapcstd.HashMap{"ID": "ID"},
	}
	m.Option.FuncName = "Func1"
	m.Register(ab.AUser{}, ab.BUser{})
	g := m.Group(func(option *mapc.Option) {
		option.FieldMappers = append(
			option.FieldMappers,
			mapcstd.HashMap{"Name": "Name"},
		)
		option.FuncName = "Func2"
	})
	g.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	assert.Len(t, gs, 2)
	got1, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func Func1(x AUser) BUser {
	return BUser{
		ID: x.ID,
	}
}
`, got1)
	got2, err := gs[1].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func Func2(x AUser) BUser {
	return BUser{
		ID:   x.ID,
		Name: x.Name,
	}
}
`, got2)
}

func Test_FuncName(t *testing.T) {
	m := mapc.New()
	m.Option.FuncComment = false
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Option.FuncName = "CustomFuncName"
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func CustomFuncName(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, got)
}

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
