package e2e

import (
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ab_same_package(t *testing.T) {
	m := mapc.New()
	m.Register(ab.AUser{}, ab.BUser{}, func(option *mapc.Option) {
		option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
		option.OutPath = "src/foo/bar.go"
	})
	gs, errs := m.Generate()
	requireNoError(t, errs)
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
