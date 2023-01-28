package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/abekoh/mapc/mapcstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func modeTestMapc() *mapc.MapC {
	m := mapc.New()
	m.Option.FieldMappers = []mapcstd.FieldMapper{
		mapcstd.HashMap{
			"ID":   "ID",
			"Name": "Name",
			"Age":  "Age",
		},
	}
	m.Option.FuncName = "ModeTestMapper"
	m.Option.OutPath = "testdata/ab/generated.go"
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	return m
}

func Test_PrioritizeGenerated(t *testing.T) {
	m := modeTestMapc()
	m.Option.Mode = mapc.PrioritizeGenerated
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

import "time"

// ModeTestMapper maps AUser to BUser.
// This function is generated by mapc.
// You can edit mapping fields.
func ModeTestMapper(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
`, got)
}

func Test_PrioritizeExisted(t *testing.T) {
	m := modeTestMapc()
	m.Option.Mode = mapc.PrioritizeExisted
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

import "time"

// ModeTestMapper maps AUser to BUser.
// This function is generated by mapc.
// You can edit mapping fields.
func ModeTestMapper(x AUser) BUser {
	return BUser{
		Name:         "Example",
		RegisteredAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		ID:           x.ID,
		Age:          x.Age,
	}
}
`, got)
}

func Test_Deterministic(t *testing.T) {
	m := modeTestMapc()
	m.Option.Mode = mapc.Deterministic
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

// ModeTestMapper maps AUser to BUser.
// This function is generated by mapc.
// DO NOT EDIT this function.
func ModeTestMapper(x AUser) BUser {
	return BUser{
		ID:   x.ID,
		Name: x.Name,
		Age:  x.Age,
		// RegisteredAt:
	}
}
`, got)
}
